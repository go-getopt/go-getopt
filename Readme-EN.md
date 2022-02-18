# Go GetOpt

[中文](./Readme.md)

## Welcome! You finally found this ~~shit~~ repo.

## Wish this ~~shit~~ code satisfy your need.

---

Go GetOpt is a go library that helps you parse cmd args like writing a shell script.

---

To make it clear, we will use

`go getopt` to represent this repo,

`shell getopt`, `getopt cmd` to represent getopt binary execution in util-linux, and

`getopt` or `C getopt` to represent `getopt` function in `libc`.

But we may still use abbreviation (use the word `getopt` directly) in a context without ambiguity. Please note the distinction.

---

## Project Status

Yes. The project is still actively maintained. You can ping me (comment `/ping @FlyingOnion`) at [Gitee issue](https://gitee.com/go-getopt/go-getopt/issues/I4UAFT) or [GitHub issue](https://github.com/go-getopt/go-getopt/issues/1). I will manually reply `/pong @you` ~~a century~~ ~~a decade~~ ~~several years~~ ~~several months~~ several days (if I'm busy), or several hours later.

Please don't do stupid things like manual DOS😕.

Pull requests are not open in GitHub. Please open requests at [Gitee upstream repo](https://gitee.com/go-getopt/go-getopt).

## Usage

```shell
go get gitee.com/go-getopt/go-getopt
```

```go
package main

import (
    "fmt"
    "os"

    // We use . to use GetOpt, Get and Shift functions conveniently.
    // You can also use common import with package name, if needed
    . "gitee.com/go-getopt/go-getopt"
)

func main() {
    // Pass os.Args, options, and longOptions string
    err := GetOpt(os.Args, "ab:c::", "a-long,b-long:,c-long::")
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    // Args is an exported global variable that stores all parsed arguments
    fmt.Println("Arguments:", Args)
    fmt.Println("Program name:", Args[0])

    // a for loop to deal with each arg like shell script does
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

Compared with shell script getopt:

```shell
getopt --test > /dev/null

[ $? -ne 4 ] &&
    echo "Error: command 'getopt --test' failed in this environment." &&
    exit 1

options=ab:c::
longOptions=a-long,b-long:,c-long::

parsed=$(getopt --options=$options --longoptions=$longOptions --name "$0" -- "$@")

[ $? -ne 0 ] &&
    echo "Error: failed to parse cmd arguments" &&
    exit 1

eval "set -- $parsed"

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
            [ -n "$2" ] &&
            echo "Option c, argument '$2'" ||
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

## Should I Use Go GetOpt

I think go getopt is suitable if you meet one or more of the following conditions:

- you are searching for a go library to parse cmd args
- you only want to get parsed strings, and you can do further conversions by yourself
- you don't want too much code of type assertions and type conversions, and you don't need the library to provide functions like `GetInt` or `MustGetInt`
- you get used to write shell script, and want to find a similar substitution in go
- you don't like the style of `flag` or `pflag`
- you don't want to use heavy and complex `cobra`

Feel free to use this library. If you have any problems just post an issue.

## Subcommand

Subcommand support is **unnecessary**, in my opinion. In shell, you could combine multiple scripts to handle your need, and so does go. Take `git push -f origin master` as an example.

Shell

```shell
# git

# Check if $1 is "push". If so, shift $1 and pass $@ to git-push script
[ "$1" == "push" ] && shift && git-push "$@" && exit

# If not, continue getopt git command
```

go

```go
if len(os.Args) >= 2 && os.Args[1] == "push" {
    gitPushArgs := make([]string, 0, len(os.Args)-1)
    gitPushArgs = append(gitPushArgs, "git-push")
    gitPushArgs = append(gitPushArgs, os.Args[2:]...)
    GetOpt(gitPushArgs, options, longOptions)
    // for loop of "git-push"
    return
}
GetOpt(os.Args, options, longOptions)
// for loop of "git"
```

## Other FAQs

### Is it concurrency-safe (or goroutine-safe)?

It's impossible, because C `getopt` and `getopt_long` use global variables `optind` and `optarg` to store middle states.

It's unnecessary either. Why do you need to parse cmd args multiple times?

### Which platforms do the library supports?

We use `cgo` to wrap C `getopt_long` function, mainstream `libc`s are all supported:

- `mingw` / `msvc` Windows
- `glibc` Debian family, CentOS family
- `musl` Alpine
- `uclibc-ng` BusyBox

### Why not use pure go to get this done?

Perhaps you means to implement a `getopt_long` with pure go, but I have been away from C for a long time. You can do it by yourself if you are interested. The source code of `getopt` and `getopt_long` functions in widely-used `libc`s are listed below.

**musl**

https://git.musl-libc.org/cgit/musl/tree/src/misc/getopt.c

https://git.musl-libc.org/cgit/musl/tree/src/misc/getopt_long.c

**glibc**

https://sourceware.org/git/?p=glibc.git;a=blob;f=posix/getopt.c;h=e9661c79faa8920253bc37747b193d1bdcb288ef;hb=HEAD

https://sourceware.org/git/?p=glibc.git;a=blob;f=posix/getopt1.c;h=990eee1b64fe1ee03e8b22771d2e88d5bba3ac68;hb=HEAD

**uclibc-ng**

https://gogs.waldemar-brodkorb.de/oss/uclibc-ng/src/master/libc/unistd/getopt.c

https://gogs.waldemar-brodkorb.de/oss/uclibc-ng/src/master/libc/unistd/getopt_long-simple.c

