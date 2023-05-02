# Cozo 数据库的 Go 语言库

[![Go](https://img.shields.io/github/v/release/cozodb/cozo-lib-go)](https://github.com/cozodb/cozo-lib-go)

本文叙述的是如何安装设置库本身。有关如何使用 CozoDB（CozoScript）的信息，见 [文档](https://docs.cozodb.org/zh_CN/latest/index.html) 。

## 安装

首先你需要根据你的操作系统和处理器，从 [GitHub](下载页面) 下载对应的预编译的 C 语言库（以 `libcozo_c` 开头的文件）。下载后需要将其解压，然后设置一些环境变量。这个 [脚本](pull_libs.sh) 在类 UNIX 系统里可以帮你把这些做了：

```bash
COZO_VERSION=0.7.0

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

注意，因为静态库比较大，Gitee 不支持存储大文件，所以我们无法提供国内的下载镜像，请想办法从 GitHub 下载。

如果想让脚本里的环境变量在当前命令行生效，你可以这么调用脚本： `. ./pull_libs.sh` 。

接下来就可以像平时一样执行 `go build` 之类的命令了。

### Windows 用户需要额外注意的事项

在 Windows 下，除了上面的步骤外，还需要预先安装 [MinGW](https://www.mingw-w64.org/) （主要是因为 Go 的编译器不支持 MSVC 编译出来的文件）。同时，下载 C 库时，请务必选择 GNU 版本（`libcozo_c` 开头， `x86_64-pc-windows-gnu.a` 结尾的文件），而不是 MSVC 版本。

或者直接用 [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)，然后在 Linux 里面运行。Cozo 在 WSL 底下跑得也更快一些。

## API

可参考 [测试文件](cozo_test.go)。

```go
/**
 * 构造函数。返回的数据库对象在不用时必须先关闭。
 *
 * @param engine:  存储引擎类型。'mem' 为纯内存的非持久化引擎，另外还支持 'sqlite'、'rocksdb' 等。
 * @param path:    存储路径。有些存储引擎用不着。
 * @param options: 默认为 nil，预编译版本的库里这个参数没用。
 */
func New(engine string, path string, options Map) (CozoDB, error)

/**
 * 数据库用完之后使用此方法关闭。如果不关闭的话原生的资源不会被释放。
 */
func (db *CozoDB) Close()

/**
 * 执行查询文本。
 *
 * @param query: 查询文本
 * @param params: 查询中可用的参数，默认为 {}
 */
func (db *CozoDB) Run(query string, params Map) (NamedRows, error)

/**
 * 导出指定的存储表
 *
 * @param relations:  需要导出表的名称
 */
func (db *CozoDB) ExportRelations(relations []string) (Map, error)

/**
 * 导入数据至存储表。
 *
 * 注意此方法不会激活触发器。
 *
 * @param data: 格式与 `exportRelations` 方法返回格式相同。需要导入的表必须预先存在。
 */
func (db *CozoDB) ImportRelations(payload Map) error

/**
 * 备份数据库。
 *
 * @param path: 备份文件路径。
 */
func (db *CozoDB) Backup(path string) error

/**
 * 恢复备份至当前数据库。当前数据库必须为空。
 *
 * @param path: 备份文件路径。
 */
func (db *CozoDB) Restore(path string) error

/**
 * 将备份文件中指定的存储表中的数据插入当前数据库中的同名表。同名表必须预先存在。
 *
 * 注意此方法不会激活任何触发器。
 *
 * @param path: 备份文件路径。
 * @param relations: 需要导入的表名。
 */
func (db *CozoDB) ImportRelationsFromBackup(path string, relations []string) error
```

## 常见问题

如果使用时遇到库加载错误之类的问题，请先克隆此项目，然后尝试在项目根目录下运行 `go test` 来检查最基本的操作能否在你的机器上正常执行。以下是一些常见错误以及其解决方案。

### 找不到 -lcozo_c

需要设置 `CGO_LDFLAGS` 环境变量来告诉 cgo 去哪里找原生库（上面关于安装的小节中有这一步）。在 Windows（PowerShell）中，设置环境变量的语法是 `$env:CGO_LDFLAGS = "-L<PATH>"`。

### 未定义的符号或引用

这是由于链接器缺少参数造成的。在 [cozo.go](cozo.go) 文件中，有以下几行：

```
#cgo LDFLAGS: -lcozo_c -lstdc++ -lm
#cgo windows LDFLAGS: -lbcrypt -lwsock32 -lws2_32 -lshlwapi -lrpcrt4
#cgo darwin LDFLAGS: -framework Security
```

这些行定义了各个平台不同的链接器需要的参数。可以看出，Windows 和 macOS 比起 Linux 来都需要额外的参数。但是实际的操作系统情况很复杂，很可能你的操作系统及处理器组合需要额外的参数。你可以尝试以缺少的符号名称在网络上进行搜索，一般来说搜索结果会指向需要链接的库名。