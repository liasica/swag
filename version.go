package swag

import "runtime/debug"

const defaultVersion = "v2.0.0"

// Version of swag. 发布的二进制由 ldflags 注入，go install 安装时回落到模块版本
var Version = defaultVersion

func init() {
	if Version != defaultVersion {
		return
	}

	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return
	}

	Version = info.Main.Version
}
