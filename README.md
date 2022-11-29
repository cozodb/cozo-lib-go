# Cozo for Golang

## Building

You need to download the compiled C library files for your system
(files starting with `libcozo_c`), uncompress it:

```bash
gunzip libcozo_c*
```

and rename it to
`libcozo_c` with the original extension (`.lib` on Windows, `.a` on everything else).

Then you need to set the environment variable

```bash
export CGO_LDFLAGS="-L/<absolute-path-to-directory-containing-the-library>"
```

for example, if you placed the library in `/home/xxx/libs`, you should use

```bash
export CGO_LDFLAGS="-L/home/xxx/libs"
```

With the environment variable set, you can run `go build`, etc. for your project.