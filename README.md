# swag

🌍 *[English](README.md) ∙ [简体中文](README_zh-CN.md)*

[![Release](https://img.shields.io/github/v/release/liasica/swag?style=flat-square)](https://github.com/liasica/swag/releases)

A fork of the [swaggo/swag](https://github.com/swaggo/swag) `v2` branch. The module path is renamed to `github.com/liasica/swag/v2`, so the CLI installs and runs straight from the Go toolchain.

Annotation syntax and CLI reference live in the [upstream documentation](https://github.com/swaggo/swag/tree/v2).

## Install

```sh
go install github.com/liasica/swag/v2/cmd/swag@latest
```

Pin a release:

```sh
go install github.com/liasica/swag/v2/cmd/swag@v2.0.0-rc5-e73d748-adaptation
```

Run without installing:

```sh
go run -mod=mod github.com/liasica/swag/v2/cmd/swag@latest init
```

Container image:

```sh
docker run --rm -v $(pwd):/code ghcr.io/liasica/swag:latest init
```

Pre-built binaries for Linux and macOS (amd64, arm64) are attached to each [release](https://github.com/liasica/swag/releases).

## Changes against upstream

| Change | Files |
| --- | --- |
| Emits an `x-enum-descriptions` array next to `x-enum-varnames` and `x-enum-comments`, on both the Swagger 2.0 and the OpenAPI 3 code path | `enums.go`, `parser.go`, `parserv3.go` |
| OpenAPI 3 query parameters resolve `$ref` properties through `components.schemas` instead of skipping them | `operationv3.go` |
| `swag fmt` separates `//` from the annotation with a space instead of a tab | `formatter.go` |
| `swag fmt` pads generic brackets: `dto.Response[pagination.Result[T]]` becomes `dto.Response[ pagination.Result[ T ] ]` | `formatter.go` |
| `Formatter.Format` no longer rewrites the caller's input buffer in place, so `swag fmt` writes the file even when the formatted result is not longer than the original | `formatter.go` |

## Branches

| Branch | Content |
| --- | --- |
| `master` | mirror of upstream `master` (v1) |
| `v2` | mirror of upstream `v2` |
| `adaptation` | v1 fork, tagged `v1.x.y-<upstream commit>-adaptation` |
| `v2-adaptation` | v2 fork, tagged `v2.x.y-<upstream commit>-adaptation` |

## Release

Tags are named `<upstream version>-<upstream commit>-adaptation`, for example `v2.0.0-rc5-e73d748-adaptation`. Pushing a `v*` tag runs two workflows: GoReleaser publishes the binaries, and the Docker workflow pushes the image to GHCR.

```sh
./release.sh
```

The script reads the upstream version and commit from the `v2` branch, tags the current `v2-adaptation` head and pushes the tag.

## Sync with upstream

```sh
git fetch upstream
git checkout v2 && git merge --ff-only upstream/v2
git checkout v2-adaptation && git merge v2
./rename.sh
go mod tidy && go build ./... && go test ./...
```

`rename.sh` rewrites the module path and the generated definition names that upstream commits bring back in.

## License

[MIT](license), inherited from upstream.
