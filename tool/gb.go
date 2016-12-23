package tool

import "github.com/voidint/gbb/config"

// GBBuilder gb编译工具
type GBBuilder struct {
	conf  *config.Config
	debug bool
}

// NewGBBuilder 返回gb编译工具实例
func NewGBBuilder(conf *config.Config, debug bool) *GBBuilder {
	return &GBBuilder{
		conf:  conf,
		debug: debug,
	}
}

// Build 切换到指定工作目录后调用编译工具编译。
func (b *GBBuilder) Build(rootDir string) error {
	return buildDir(b.conf, b.debug, rootDir)
}
