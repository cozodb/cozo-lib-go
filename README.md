# Cozo for Golang

[![Go](https://img.shields.io/github/v/release/cozodb/cozo-lib-go)](https://github.com/cozodb/cozo-lib-go)

This document describes how to set up the Cozo module for use in Golang projects.
To learn how to use CozoDB (CozoScript), read the [docs](https://docs.cozodb.org/en/latest/index.html).

## Setup

You need to download the compiled C library files for your system
(files starting with `libcozo_c`) from the [release page](https://github.com/cozodb/cozo/releases), 
uncompress it somewhere, and set a few environment
variables before you can use this library.
The following script ([here](pull_libs.sh)) does this for you on UNIX systems:

```bash
COZO_VERSION=0.7.5

COZO_PLATFORM=x86_64-unknown-linux-gnu # for Linux
#COZO_PLATFORM=aarch64-apple-darwin # uncomment for ARM Mac
#COZO_PLATFORM=x86_64-apple-darwin # uncomment  for Intel Mac
#COZO_PLATFORM=x86_64-pc-windows-gnu # uncomment for Windows PC

URL=https://github.com/cozodb/cozo/releases/download/v${COZO_VERSION}/libcozo_c-${COZO_VERSION}-${COZO_PLATFORM}.a.gz

mkdir libs
echo "Download from ${URL}"
curl -L $URL -o libs/libcozo_c.a.gz
gunzip -f libs/libcozo_c.a.gz
export CGO_LDFLAGS="-L/${PWD}/libs"
```

If you used the [script](pull_libs.sh), you need to call it thus
to make the environment variables available: `. ./pull_libs.sh`

Now you can run `go build`, etc. for your project.

### Note for Windows users

On Windows, in addition to following the above instructions, 
you must have [MinGW](https://www.mingw-w64.org/) installed 
(Go doesn't seem to work with MSVC compiliers), and you must use 
the GNU version of the static library (`libcozo_c` ending in `x86_64-pc-windows-gnu.a`),
not the MSVC version.

Or just use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install).
It is much easier and Cozo runs much faster under WSL.

## API

See the [test file](cozo_test.go) for example usage.

```go
/**
 * Constructor, the returned database must be closed after use.
 *
 * @param engine:  'mem' is for the in-memory non-persistent engine.
 *                 'sqlite', 'rocksdb' and maybe others are available,
 *                 depending on compile time flags.
 * @param path:    path to store the data on disk,
 *                 may not be applicable for some engines such as 'mem'
 * @param options: defaults to nil, ignored by all the engines in the published NodeJS artefact
 */
func New(engine string, path string, options Map) (CozoDB, error)

/**
 * You must call this method for any database you no longer want to use:
 * otherwise the native resources associated with it may linger for as
 * long as your program runs. Simply `delete` the variable is not enough.
 */
func (db *CozoDB) Close()

/**
 * Runs a query
 *
 * @param query: the query
 * @param params: the parameters as key-value pairs, defaults to {} if nil
 */
func (db *CozoDB) Run(query string, params Map) (NamedRows, error)

/**
 * Export several relations
 *
 * @param relations:  names of relations to export, in an array.
 */
func (db *CozoDB) ExportRelations(relations []string) (Map, error)

/**
 * Import several relations
 *
 * Note that triggers are _not_ run for the relations, if any exists.
 * If you need to activate triggers, use queries with parameters.
 *
 * @param data: in the same form as returned by `exportRelations`. The relations
 *              must already exist in the database.
 */
func (db *CozoDB) ImportRelations(payload Map) error

/**
 * Backup database
 *
 * @param path: path to file to store the backup.
 */
func (db *CozoDB) Backup(path string) error

/**
 * Restore from a backup. Will fail if the current database already contains data.
 *
 * @param path: path to the backup file.
 */
func (db *CozoDB) Restore(path string) error

/**
 * Import several relations from a backup. The relations must already exist in the database.
 *
 * Note that triggers are _not_ run for the relations, if any exists.
 * If you need to activate triggers, use queries with parameters.
 *
 * @param path: path to the backup file.
 * @param relations: the relations to import.
 */
func (db *CozoDB) ImportRelationsFromBackup(path string, relations []string) error
```

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
Windows flags, and additional macOS flags. If for some reason these flags
are not enough, you can google what the compiler tells you is missing
to see what flags you should add (and open an issue about the problem).