// Copyright (c) 2022 Ruicong Huang (huangrc)
// Go getopt is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// 		 http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package getopt

/*
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <getopt.h>
*/
import "C"
import (
	"errors"
	"strconv"
	"unsafe"
)

// Args contains all parsed arguments by GetOpt, including "--" and positional args.
// To get values, you can either use Get and Shift like writing a shell script,
// or copy it or use indices if you don't want to break original values.
//
// GetOpt 转换后的参数列表（包括 "--" 和之后的位置参数）。
// 获取参数值的时候，可以像 shell 那样用 Get 和 Shift 来处理，
// 不想破坏原值的话，也可以复制一份或者用游标。
// 总之 GetOpt 解析参数的任务已经完成，之后就随你们便啦。
var Args []string = make([]string, 0, 64)

func Shift(n int) {
	argc := len(Args)
	if argc <= 1 {
		return
	}
	m := n
	if m > argc-1 {
		m = argc - 1
	}
	copy(Args[1:], Args[1+m:])
	Args = Args[:argc-m]
}

func Get(n int) (str string) {
	if n >= 0 && n < len(Args) {
		str = Args[n]
	}
	return
}

func errEmptyLongOption(longOptions string, index int) error {
	return errors.New("empty long option at index [" + strconv.Itoa(index) + "]: '" + longOptions + "'")
}

// parseLongOptions parses the longOptions string to a _Go_options slice.
// There may be errors. If we return []C.struct_options directly,
// we should deal with the problem of freeing parsed C-Type strings. It's annoying.
// Functions returning C-Type strings (i.e. *C.char or _Ctype_char) should be always correct without errors.
//
// 解析过程有可能会报错，如果直接返回 []C.struct_options，
// 则报错返回时，已经解析过的 C 字符串需要一个个 free 掉。
// 因此这里多做一步，调用 _Go_options.getCTypeOptions 生成 []C.struct_options。
// 只要是返回值中包含 C *char 字符串的 go 方法，必须正确返回，返回值不能包含 error。
func parseLongOptions(longOptions string) (_Gotype_long_options, error) {
	n := len(longOptions)
	opts := make(_Gotype_long_options, 0, 20)
	for i, j := 0, 0; j < n; j++ {
		if ch := longOptions[j]; !(ch == ',' || ch == ' ' || ch == '\t' || ch == '\n') {
			continue
		}
		// deal with ",," or "  " (multiple splits)
		if j == i {
			i++
			continue
		}
		argOpt := C.no_argument
		if longOptions[j-1] == ':' {
			if j == i+1 {
				// deal with ",:,"
				return nil, errEmptyLongOption(longOptions, i)
			}
			if longOptions[j-2] == ':' {
				if j == i+2 {
					// deal with ",::,"
					return nil, errEmptyLongOption(longOptions, i)
				}
				argOpt = C.optional_argument
			} else {
				argOpt = C.required_argument
			}
		}
		opts = append(opts, _Gotype_long_option{
			name:   longOptions[i : j-argOpt],
			hasArg: argOpt,
		})
		i = j + 1
	}
	return opts, nil
}

func GetOpt(args []string, options, longOptions string) error {
	argc := len(args)
	argv := make([]*C.char, 0, argc)
	for _, arg := range args {
		argv = append(argv, C.CString(arg))
	}
	longOpts, err := parseLongOptions(longOptions)
	if err != nil {
		return err
	}
	cTypeLongOpts := longOpts.getCTypeOptions()

	var longIndex C.int

	p0 := C.int(argc)
	p1 := &argv[0]
	p2 := C.CString(options)
	p3 := &cTypeLongOpts[0]
	p4 := &longIndex

	Args = append(Args[:0], args[0])
	for loop := true; loop; {
		opt := C.getopt_long(p0, p1, p2, p3, p4)
		switch opt {
		case C.EOF:
			loop = false
			continue
		case 0:
			// match a long option
			o := cTypeLongOpts[longIndex]
			Args = append(Args, "--"+C.GoString(o.name))
			if o.has_arg > 0 {
				Args = append(Args, C.GoString(C.optarg))
			}

		case 1:
			// match no options
			Args = append(Args, C.GoString(C.optarg))
		default:
			b := byte(opt)
			Args = append(Args, string([]byte{'-', b}))
			for i := 0; i < len(options)-1; i++ {
				if options[i] == b && options[i+1] == ':' {
					Args = append(Args, C.GoString(C.optarg))
					break
				}
			}
		}
	}
	Args = append(Args, "--")

	for optIndex := C.optind; optIndex < p0; optIndex++ {
		Args = append(Args, C.GoString(argv[optIndex]))
	}
	// fmt.Println(Args)

	// Finished parsing and start cleaning

	for _, arg := range argv {
		C.free(unsafe.Pointer(arg))
	}
	C.free(unsafe.Pointer(p2))
	for _, opt := range cTypeLongOpts {
		C.free(unsafe.Pointer(opt.name))
	}
	return nil
}

type _Gotype_long_option struct {
	name   string
	hasArg int
}

type _Gotype_long_options []_Gotype_long_option

var flag C.int

func (opts _Gotype_long_options) getCTypeOptions() []C.struct_option {
	cTypeOptions := make([]C.struct_option, 0, 1+len(opts))
	for i, opt := range opts {
		cTypeOptions = append(cTypeOptions, C.struct_option{
			name:    C.CString(opt.name),
			has_arg: C.int(opt.hasArg),
			flag:    &flag,
			val:     C.int(i),
		})
	}
	// 添加一个全〇字段结束
	cTypeOptions = append(cTypeOptions, C.struct_option{})
	return cTypeOptions
}
