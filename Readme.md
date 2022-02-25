# Go GetOpt

[English](./Readme-EN.md)

## 你终于来了，我们等了好久！

## 希望我们的代码不会让你失望。

---

Go GetOpt，让你在 go 里解析命令行参数~~无聊地~~跟写 shell 脚本一样。

---

为了不引起混淆，以下说明将使用

`go getopt` 表示本代码仓库

`shell getopt`、`getopt 命令` 表示 util-linux 中的 getopt 二进制程序

`getopt`（或 `C getopt`）表示 `libc` 中的 `getopt` 方法

但在某个上下文（如标题说明了该段是 shell getopt）中可能有时会直接使用 `getopt` 指代。请各位注意区分。

---

## 项目状态

别急，作者还在呢。在 [issue](https://gitee.com/go-getopt/go-getopt/issues/I4UAFT) 下评论 `"/ping @FlyingOnion"`，我就会回复 `"/pong @你"`（手动的）。也可以在 Github https://github.com/go-getopt/go-getopt 提工单。不过有时忙起来可能会~~几百年以后~~几天后再回。

GitHub 上没开 PR，提 PR 请右转 [Gitee 上游仓库](https://gitee.com/go-getopt/go-getopt)（因为众所周知的原因，我已经半放弃 GitHub 了）。

---

## 怎么用

```shell
go get gitee.com/go-getopt/go-getopt
```

```go
package main

import (
    "fmt"
    "os"

    // 这里为了方便，直接使用 . 进行 import。
    // 这样可以直接调用 GetOpt、Get 和 Shift 方法。
    . "gitee.com/go-getopt/go-getopt"
)

func main() {
    // 传入 os.Args、options 和 longOptions 字符串参数即可
    err := GetOpt(os.Args, "ab:c::", "a-long,b-long:,c-long::")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // 解析后的参数列表存储在全局变量 Args 中
    fmt.Println("Arguments:", Args)
    fmt.Println("Program name:", Args[0])

    // 接下来的步骤就和 shell 差不多了
    for loop := true; loop; {
        switch Get(1) {
        case "-a", "--a-long":
            fmt.Println("Option a")
            Shift(1)
        case "-b", "--b-long":
            fmt.Println("Option b, argument '" + Get(2) + "'")
            Shift(2)
        case "-c", "--c-long":
            if Get(2) == "" {
                fmt.Println("Option c, no argument")
            } else {
                fmt.Println("Option c with arg '" + Get(2) + "'")
            }
            Shift(2)
        case "--":
            Shift(1)
            loop = false
        default:
            fmt.Fprintln(os.Stderr, "Error: wrong argument '"+arg1+"'")
            os.Exit(1)
        }
    }
    fmt.Println("Remains:", Args[1:])
}
```

对比一下 shell getopt 解析命令行的脚本：

```shell
# 检查 getopt 命令是否正常运行
getopt --test > /dev/null

[ $? -ne 4 ] &&
    echo "Error: command 'getopt --test' failed in this environment." &&
    exit 1

# 设定 options 和 longOptions，调用 getopt 命令
options=ab:c::
longOptions=a-long,b-long:,c-long::

parsed=$(getopt --options=$options --longoptions=$longOptions --name "$0" -- "$@")

[ $? -ne 0 ] &&
    echo "Error: failed to parse cmd arguments" &&
    exit 1

eval "set -- $parsed"

# 循环判断是哪个 flag，处理完后 shift 掉
while true; do
    case "$1" in
        -a|--a-long)
            echo 'Option a'
            shift
        ;;
        -b|--b-long)
            echo "Option b, argument '$2'"
            shift 2
        ;;
        -c|--c-long)
            [ -n "$2" ] && \
            echo "Option c, argument '$2'" || \
            echo 'Option c, no argument'
            shift 2
        --)
            shift
            break
        *)
            echo "Error: wrong argument '$1'"
            exit 1
            ;;
    esac
done
echo "Remains: $@"
```

## Go GetOpt 适合哪些人用

通常来说，你能看到这里，说明你对这~~破~~代码有点兴趣了。

如果你符合以下的一条或多条，可能这个库会适合你：

- 想用一个库让 go 程序解析命令行参数方便一点。
- 只想解析出字符串形式参数，然后自己做处理或转换。
- 不想让类型断言、类型转换代码到处乱飞，也不需要调用的库提供 `GetInt`、`MustGetInt` 之类的方法。
- ~~忘不了前任~~ 习惯了写 shell，想找个差不多的库~~接盘~~。
- 不喜欢 `flag`、`pflag` 这种类型的解析方式（`pflag` 也很久没维护了）。
- 不想用 `cobra` 这种很繁琐的库。

如果我说中了，那这~~破~~东西你尽管拿去用。不好用请尽管提 issue。

## 子命令（subcommand）

子命令其实并不是必须的，我们写 shell 的时候可以用多个脚本（或方法）共同完成。

以 `git` 和 `git push -f origin master` 为例。以下分别给出 shell 和 go 的处理方法。

shell

```shell
# git

# 检查参数 1 是否为 "push"。如果是的话 shift 掉参数 1，然后将剩下的参数传到 git-push 中。
[ "$1" == "push" ] && shift && git-push "$@" && exit

# 若不是，则继续 git 命令的 getopt 解析。
```

go

```go
if len(os.Args) >= 2 && os.Args[1] == "push" {
    gitPushArgs := make([]string, 0, len(os.Args)-1)
    gitPushArgs = append(gitPushArgs, "git-push")
    gitPushArgs = append(gitPushArgs, os.Args[2:]...)
    GetOpt(gitPushArgs, options, longOptions)
    // continue for loop of "git push"
    return
}
```

## 其他问题

### 并发（协程）安全吗？

没办法做到，也没必要。C 的 `getopt` 和 `getopt_long` 方法本身就不是并发安全的（用了全局变量 `optind`、`optarg` 来存储中间状态）。

而且，命令行应该只需要解析一次就可以了吧。有必要多次解析吗🤔？

### 支持哪些平台？

这个库是用 `cgo` 包装 `libc` 中的 `getopt_long` 方法实现的。原理和 shell getopt 命令行程序差不多。目前主流的 `libc` 都是支持的：

- `mingw` / `cygwin` ~~/ `msvc`~~ Windows 系
- `glibc` Debian 系、CentOS 系
- `musl` Alpine
- `uclibc-ng` BusyBox

对不起搞错了。`MSVC` 没有 `getopt` 和 `getopt_long`。用 Windows 写 go 程序的人可以采用以下（可能有用的）替代方案。

- https://github.com/DavidGamba/go-getoptions
- https://github.com/pborman/getopt
- https://github.com/droundy/goopt

### 怎么不用纯 go 写一个呢？

C 很久没用了，所以就没再仔细研究 `getopt` 和 `getopt_long` 的源码。

各 `libc` 的 `getopt` 和 `getopt_long` 源码的地址请自取，有兴趣的也可以自己研究一下。

~~少年，我看你骨骼精奇，是个搞 996 的好苗子！~~

**musl**

https://git.musl-libc.org/cgit/musl/tree/src/misc/getopt.c

https://git.musl-libc.org/cgit/musl/tree/src/misc/getopt_long.c

**glibc**

https://sourceware.org/git/?p=glibc.git;a=blob;f=posix/getopt.c;h=e9661c79faa8920253bc37747b193d1bdcb288ef;hb=HEAD

https://sourceware.org/git/?p=glibc.git;a=blob;f=posix/getopt1.c;h=990eee1b64fe1ee03e8b22771d2e88d5bba3ac68;hb=HEAD

**uclibc-ng**

https://gogs.waldemar-brodkorb.de/oss/uclibc-ng/src/master/libc/unistd/getopt.c

https://gogs.waldemar-brodkorb.de/oss/uclibc-ng/src/master/libc/unistd/getopt_long-simple.c