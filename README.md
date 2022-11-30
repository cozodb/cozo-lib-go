# Cozo for Golang

This document describes how to set up the Cozo module for use in Golang projects.
To learn how to use CozoDB (CozoScript), follow
the [tutorial](https://nbviewer.org/github/cozodb/cozo-docs/blob/main/tutorial/tutorial.ipynb)
first and then read the [manual](https://cozodb.github.io/current/manual/). You can run all the queries
described in the tutorial with an in-browser DB [here](https://cozodb.github.io/wasm-demo/).


## Usage

You need to download the compiled C library files for your system
(files starting with `libcozo_c`), uncompress it:

```bash
gunzip libcozo_c*
```

and rename it to `libcozo_c.a`.

Then you need to set the environment variable

```bash
export CGO_LDFLAGS="-L/<absolute-path-to-directory-containing-the-library>"
```

for example, if you placed the library in `/home/xxx/libs`, you should use

```bash
export CGO_LDFLAGS="-L/home/xxx/libs"
```

With the environment variable set, you can run `go build`, etc. for your project.

### Note for Windows users

On Windows, in addition to following the above instructions, 
you must have [MinGW](https://www.mingw-w64.org/) installed 
(Go doesn't seem to work with MSVC compiliers), and you must use 
the GNU version of the static library (`libcozo_c` ending in `x86_64-pc-windows-gnu.a`).

Or just use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).
It is much easier and Cozo runs much faster under WSL.

## Frequently encountered problems

If you encounter a problem when trying to use this library for your project,
you should first clone this repo and try to run `go test` to see
if the library works at all on your machine.
The following are some frequently-encountered error messages:

### Cannot find -lcozo_c: No such file or directory

You need to set the `CGO_LDFLAGS` variable and tell cgo where to find
the static libraries, as described above.
On Windows (PowerShell), the syntax is `$env:CGO_LDFLAGS = "-L<PATH>"`

### Undefined symbols / undefined reference

This means that some linker flags are missing. If you look at the file 
[cozo.go](cozo.go), you can find the following lines

```
#cgo LDFLAGS: -lcozo_c -lstdc++ -lm
#cgo windows LDFLAGS: -lbcrypt -lwsock32 -lws2_32 -lshlwapi -lrpcrt4
#cgo darwin LDFLAGS: -framework Security
```

The three lines defines the linker flags on all platforms, additional
Windows flags, and additional MacOS flags. If for some reason these flags
are not enough, you can google what the compiler tells you is missing
to see what flags you should add (and open an issue about the problem).