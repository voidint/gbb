# gbb
[![Build Status](https://travis-ci.org/voidint/gbb.png)](https://travis-ci.org/voidint/gbb)
[![Coverage Status](https://coveralls.io/repos/github/voidint/gbb/badge.svg)](https://coveralls.io/github/voidint/gbb)

## 目录
- [应用场景](#应用场景)
	- [场景一](#场景一)
	- [场景二](#场景二)
- [特性](#特性)
- [安装](#安装)
- [基本使用](#基本使用)
	- [准备](#准备)
	- [step0](#step0)
	- [step1](#step1)
	- [step2](#step2)
- [gbb.json](#gbbjson)
- [changelog](#changelog)
	
## 应用场景
### 场景一
如果项目中包含了多个main入口文件，比如👇

```shell
$ tree ./github.com/voidint/test
./github.com/voidint/test
├── cmd
│   ├── apiserver
│   │   └── main.go
│   ├── dbtool
│   │   └── main.go
│   └── init
│       └── main.go
└── gbb.json

4 directories, 4 files
```
对于这样子目录结构，该怎么去编译这些个程序？假设使用原生的`go build/install`工具，也许会这么做：

- 输入完整的路径编译

	``` shell
	$ go install github.com/voidint/test/cmd/apiserver
	$ go install github.com/voidint/test/cmd/dbtool
	$ go install github.com/voidint/test/cmd/init
	```
	
- 逐个切换工作目录后执行`go build/install`

	``` shell
	$ cd github.com/voidint/test/cmd/apiserver && go install && cd -
	$ cd github.com/voidint/test/cmd/dbtool && go install && cd -
	$ cd github.com/voidint/test/cmd/init && go install && cd -
	```
操作完之后是否会觉得很繁琐？如果一天需要编译这个项目几十次，那会相当低效。可惜，目前`go build/install`好像并不支持在项目根目录下编译子孙目录中所有的main入口文件。

### 场景二
昨天跑得好好的程序突然出了问题，查看它的版本号，机器冷冰冰地告诉你👇

``` shell
$ xxx --version
xxx version 1.0.12
```
如果没有详细的发布记录，那么此时的你一定是崩溃的。因为实在不知道这个`1.0.12`到底是什么时候编译的，更加不知道是从哪份源代码编译而来，想要找出其中的bug，难度大了不少。

那么，同样的场景下，机器告诉你的信息是这样，那debug是否容易多了呢？！

``` shell
$ xxx --version
xxx version 1.0.12
date: 2016-12-18T15:37:09+08:00
commit: db8b606cfc2b24a24e2e09acac24a52c47b68401
```


如果以上的场景你也似曾相识，那么也许`gbb`就能帮到你，耐心往下👀吧。

## 特性
根据以上的场景描述，可以简单地将主要特性归纳为如下几条：

- 一键编译当前目录下所有go package。
- 支持编译时自动“嵌入”信息到二进制可执行文件，典型的如嵌入`编译时间`和源代码`Commit`信息到二进制可执行文件的版本信息当中。
- 首次运行会在项目根目录自动生成`gbb.json`配置文件，今后编译操作所需的信息都从此文件读取，不再打扰用户，做安静美男子。

## 安装
1. 拉取源代码

	``` shell
	$ go get -u -v github.com/voidint/gbb
	```
1. 编译（默认情况下`go get`就会编译）

	```
	$ cd $GOPATH/src/github.com/voidint/gbb && go install
	```
1. 将可执行文件`gbb`放置到环境变量`PATH`
1. 执行`which gbb`确认是否安装成功
1. 若`gbb`重名，那么建议设置别名，比如`alias gbb=gbb2`。


## 基本使用
`gbb`是自举的，换句话说，使用以上步骤安装的`gbb`可执行二进制文件是可以编译gbb源代码的。类似👇

```shell
$ cd $GOPATH/src/github.com/voidint/gbb && gbb --debug
==> go build -ldflags  '-X "github.com/voidint/gbb/build.Date=2016-12-17T17:00:04+08:00" -X "github.com/voidint/gbb/build.Commit=db8b606cfc2b24a24e2e09acac24a52c47b68401"'

$ ls -l ./gbb
-rwxr-xr-x  1 voidint  staff  4277032 12 17 17:00 ./gbb
```
可以看到当前目录下已经多了一个可执行二进制文件。没错，这个`./gbb`就是使用已经安装的`gbb`编译源代码后的产物。

怎么使用`gbb`来代替`go build/install`或者`gb`来完成日常的代码编译工作呢？简单，跟着下面的步骤尝试一下，立马就学会了。

### 准备
既然需要演示使用方法，必然就需要有个go项目。下面以`gbb`项目为例来展开。

为了从零开始我们的演示，请先把源代码目录下的`gbb.json`文件删除。`gbb.json`的作用以及文件内容的含义暂且不表，下文自然会提到。

``` 
$ rm -f gbb.json
```

首先，明确下使用`gbb`工具能干什么事？

如场景一所描述的那样，如果日常都是使用`go build/install`去应对编译工作，并且也不需要在二进制可执行文件中“嵌入”什么信息，那么，请跳过下面的step0，直接阅读[step1](https://github.com/voidint/gbb#step1)。

如果对“嵌入”编译时间、Commit这类信息到二进制可执行文件中有一定兴趣，那么建议从头至尾通读一遍吧。

### step0
为了在版本信息中显示`编译时间`和`commit号`这两个关键信息（并不限于这两个信息），需要先定义两个变量（变量不需要赋初值）。

```
package build
var (
	Date   string
	Commit string
)

```
然后，设法在功能代码中用上这两个变量。类似👇。

``` go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/build"
)

var (
	// Version 版本号
	Version = "0.1.0"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gbb version %s\n", Version)
		if build.Date != "" {
			fmt.Printf("date: %s\n", build.Date)
		}
		if build.Commit != "" {
			fmt.Printf("commit: %s\n", build.Commit)
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
```

### step1
在项目根目录执行`gbb init`，按照`gbb init`的提示，逐步填写完信息并最终生成`gbb.json`文件。关于`gbb.json`，请参见下文的[详细说明](https://github.com/voidint/gbb#gbbjson)。

如果是场景一的使用场景，那么只需要填写`tool`（实际调用的编译工具）后按要求终止流程即可。

``` shell
$ gbb init
This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.

Press ^C at any time to quit.
tool: (go_install) go_build
Do you want to continue?[y/n] n
About to write to /Users/voidint/cloud/workspace/go/lib/src/github.com/voidint/gbb/gbb.json:

{
    "version": "0.4.0",
    "tool": "go build"
}

Is this ok?[y/n] y
```

如果满足场景二所描述的情况，那么还需要继续信息搜集流程。

``` shell
$ gbb init
This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.

Press ^C at any time to quit.
tool: (go_install) go_build
Do you want to continue?[y/n] y
importpath: (main) github.com/voidint/gbb/build
variable: Date
value: {{.Date}}
Do you want to continue?[y/n] y
variable: Commit
value: {{.GitCommit}}
Do you want to continue?[y/n] n
About to write to /Users/voidint/cloud/workspace/go/lib/src/github.com/voidint/gbb/gbb.json:

{
    "version": "0.4.0",
    "tool": "go build",
    "importpath": "github.com/voidint/gbb/build",
    "variables": [
        {
            "variable": "Date",
            "value": "{{.Date}}"
        },
        {
            "variable": "Commit",
            "value": "{{.GitCommit}}"
        }
    ]
}

Is this ok?[y/n] y
```



### step2
在项目根目录执行`gbb --debug`，就会按照同级目录下`gbb.json`中配置执行编译操作。若`gbb.json`文件不存在，`gbb init`会被自动调用。

```
$ gbb --debug
==> go build -ldflags  '-X "github.com/voidint/gbb/build.Date=2016-12-17T22:18:32+08:00" -X "github.com/voidint/gbb/build.Commit=db8b606cfc2b24a24e2e09acac24a52c47b68401"'
```
编译完后在目录下多出一个编译后的二进制文件。试着输出版本信息，看看是否实现我们设定的目标了。

```
$ ./gbb version
gbb version 0.4.0
date: 2016-12-17T22:18:32+08:00
commit: db8b606cfc2b24a24e2e09acac24a52c47b68401
```
😊

## gbb.json
`gbb.json`可以认为是`gbb`工具的配置文件，通过`gbb init`自动创建（感谢`npm init`）。通常它的格式是这样：

``` json
{
    "version": "0.4.0",
    "tool": "go build",
    "importpath": "github.com/voidint/gbb/build",
    "variables": [
        {
            "variable": "Date",
            "value": "{{.Date}}"
        },
        {
            "variable": "Commit",
            "value": "{{.GitCommit}}"
        }
    ]
}
```

- `version`: 版本号。gbb根据自身版本号自动写入gbb.json。
- `tool`: gbb实际调用的编译工具。可选值包括：`go_build`、`go_install`、`gb_build`。注意：这个值不能包含空格[issue](https://github.com/voidint/gbb/issues/1)，因此暂时通过下划线`_`连接。
- `importpath`: 包导入路径，也就是`Date`、`Commit`这类变量所在包的导入路径，如`github.com/voidint/gbb/build`。
- `variables`: 变量列表。列表中的每个元素都包含`variable`和`value`两个属性。
	- `variable`变量名，比如`Date`。
	- `value`变量表达式
		- 内置变量表达式
			- `{{.Date}}`: 输出[RFC3339](http://www.ietf.org/rfc/rfc3339.txt)格式的系统时间。
			- `{{.GitCommit}}`: 输出当前分支最近一次`git hash`字符串。
		- 命令形式的变量表达式
			- 以`$(`开头，`)`结尾，中间的字符串内容会被当做命令被执行。如表达式`$(date)`，`date`命令的输出将会作为变量表达式最终的求值结果。
	
	
## changelog
### 0.4.0 - 2017/04/08
- 支持编译当前目录下所有go package，不再仅限于编译main package。[#10](https://github.com/voidint/gbb/issues/10)
- `gbb.json`中的配置项`package`重命名为`importpath`。[#9](https://github.com/voidint/gbb/issues/9)
- 新增命令行选项`--config`用于自定义配置文件路径。[#16](https://github.com/voidint/gbb/issues/16)
- 切换目录并编译后重新切换回源目录。[#17](https://github.com/voidint/gbb/issues/17)
- 当gbb.json的版本号高于gbb程序版本号时给出程序升级提醒。[#19](https://github.com/voidint/gbb/issues/19)

### 0.3.0 - 2017/01/09
- 若开启debug模式`gbb --debug`，那么变量表达式求值过程详情也一并输出。[#12](https://github.com/voidint/gbb/issues/12) [#6](https://github.com/voidint/gbb/issues/6)
- 变量表达式首字母大写。[#11](https://github.com/voidint/gbb/issues/11)
- 支持命令形式的变量表达式。[#7](https://github.com/voidint/gbb/issues/7)

### 0.2.0 - 2016/12/30
- `gbb.json`中的配置项——`package`和`variables`由必选项改为可选项。其中，在`variables`选项为空的情况下，实际在调用编译工具编译时不再加上形如`-ldflags '-X "xxx.yyy=zzz"'`的参数。[#8](https://github.com/voidint/gbb/issues/8)
- 若程序版本号与`gbb.json`中的`version`值不一致，就会强制重新生成`gbb.json`文件。

### 0.1.1 - 2016/12/24
- 支持通过`gbb init`初始化配置信息并生成`gbb.json`配置文件。
- 支持在项目根目录下，一键编译所有入口源代码文件，并生成一个或者多个可执行二进制文件。[#4](https://github.com/voidint/gbb/issues/4)
- 支持调用`gb`或者`go build/install`，并为编译生成的可执行文件提供丰富的版本信息中，包括但不限于：`编译时间`、`源代码版本控制commit`等。