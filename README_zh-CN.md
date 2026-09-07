# swag

🌍 *[English](README.md) ∙ [简体中文](README_zh-CN.md)*

[![Release](https://img.shields.io/github/v/release/liasica/swag?style=flat-square)](https://github.com/liasica/swag/releases)

[swaggo/swag](https://github.com/swaggo/swag) `v2` 分支的 fork，模块路径改为 `github.com/liasica/swag/v2`，可直接用 Go 工具链安装和运行。

注解语法与命令行参数见[上游文档](https://github.com/swaggo/swag/tree/v2)。

## 安装

```sh
curl -sSfL https://raw.githubusercontent.com/liasica/swag/v2-adaptation/install.sh | bash
```

脚本从最新发布中下载与当前平台匹配的二进制，用 `checksums.txt` 校验后安装到 `/usr/local/bin`。`SWAG_VERSION` 指定版本，`SWAG_BIN_DIR` 指定安装目录：

```sh
curl -sSfL https://raw.githubusercontent.com/liasica/swag/v2-adaptation/install.sh | SWAG_VERSION=v2.0.0-rc5-e73d748-adaptation SWAG_BIN_DIR=~/.local/bin bash
```

容器镜像：

```sh
docker run --rm -v $(pwd):/code ghcr.io/liasica/swag:latest init
```

走 Go 工具链：

```sh
go install github.com/liasica/swag/v2/cmd/swag@latest
go run -mod=mod github.com/liasica/swag/v2/cmd/swag@latest init
```

Linux 与 macOS（amd64、arm64）的二进制挂在每个[发布](https://github.com/liasica/swag/releases)上。

## 相对上游的改动

| 改动 | 文件 |
| --- | --- |
| 在 `x-enum-varnames`、`x-enum-comments` 之外输出 `x-enum-descriptions` 数组，Swagger 2.0 与 OpenAPI 3 两条路径都生效 | `enums.go`、`parser.go`、`parserv3.go` |
| OpenAPI 3 的 query 参数遇到 `$ref` 属性时从 `components.schemas` 取回定义，不再跳过 | `operationv3.go` |
| `swag fmt` 用空格而不是 tab 分隔 `//` 与注解 | `formatter.go` |
| `swag fmt` 给泛型方括号补空格：`dto.Response[pagination.Result[T]]` 变为 `dto.Response[ pagination.Result[ T ] ]` | `formatter.go` |
| `Formatter.Format` 不再就地改写调用方传入的缓冲区，格式化结果不比原文长时 `swag fmt` 也会正常写回文件 | `formatter.go` |

## 分支

| 分支 | 内容 |
| --- | --- |
| `master` | 上游 `master`（v1）镜像 |
| `v2` | 上游 `v2` 镜像 |
| `adaptation` | v1 fork，标签为 `v1.x.y-<上游 commit>-adaptation` |
| `v2-adaptation` | v2 fork，标签为 `v2.x.y-<上游 commit>-adaptation` |

## 发版

标签格式为 `<上游版本>-<上游 commit>-adaptation`，例如 `v2.0.0-rc5-e73d748-adaptation`。推送 `v*` 标签会触发两个工作流：GoReleaser 发布二进制，Docker 工作流把镜像推到 GHCR。

```sh
./release.sh
```

脚本从 `v2` 分支读取上游版本与 commit，给当前 `v2-adaptation` 的 HEAD 打标签并推送。

## 同步上游

```sh
git fetch upstream
git checkout v2 && git merge --ff-only upstream/v2
git checkout v2-adaptation && git merge v2
./rename.sh
go mod tidy && go build ./... && go test ./...
```

`rename.sh` 负责改写上游提交带回来的模块路径与生成的定义名。

## 许可

[MIT](license)，沿用上游。
