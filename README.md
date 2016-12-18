# gbb
昨天跑得好好的程序突然出了问题，查看它的版本号，机器冷冰冰地告诉你👇

``` shell
$ xxx --version
xxx version 1.0.12
```
如果没有详细的发布记录信息，我想此时的你一定是崩溃的。因为实在不知道这个`1.0.12`到底是什么时候编译的，更加不知道是从哪份源代码编译而来，想要找出其中的bug，难度大了不少。

那么，同样的场景下，机器告诉你的信息是这样，那debug是否容易多了呢？！

``` shell
$ xxx --version
xxx version 1.0.12
date: 2016-12-18T15:37:09+08:00
commit: db8b606cfc2b24a24e2e09acac24a52c47b68401
```

如果以上的场景你也似曾相识，那么也许`gbb`就能帮到你，耐心往下👀吧。
## 安装
1. 拉取源代码

	``` shell
	$ go get -u github.com/voidint/gbb
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
可以看到当前目录下已经多了一个可执行的二进制文件。没错，这个`./gbb`就是使用已经安装的`gbb`编译源代码后的产物。

如果是一个全新的项目，该怎么使用`gbb`来代替`go build/install`或者`gb`来完成日常的代码编译工作呢？很简单，跟着下面的步骤尝试一下，立马就能学会了。

### 准备
既然需要演示使用方法，必然就需要有个go项目。我这里就以`gbb`项目为例来展开。

为了从零开始我们的演示，请先把源代码目录下的`gbb.json`文件删除。`gbb.json`的作用以及文件内容的含义暂且不表，下文自然会提到。

``` 
$ rm -f gbb.json
```

首先，明确下`gbb`工具要干什么事？我知道这个措辞很烂，在没有更好的措辞之前，先将就着看吧。
> 对go install/build、gb等golang编译工具进行包装，使编译得到的二进制文件的版本信息中包含编译时间戳、git commit等信息。

其次，看`gbb`的版本信息👇

``` shell
$ gbb version
gbb version v0.0.1
date: 2016-12-17T15:37:09+08:00
commit: db8b606cfc2b24a24e2e09acac24a52c47b68401
```

这个版本信息，除了常规的`v0.0.1`，还有这个`gbb`二进制文件编译生成的时间，以及项目所使用的源代码管理工具`git`的最近一次`commit`号。这样的版本信息是否比简单的一个`v0.0.1`要更加友好呢？丰富的版本信息也为`debug`降低了难度，因为这个二进制能和仓库中的源代码唯一对应了。

### step0
为了在版本信息中显示`编译时间`和`commit号`这两个关键信息，需要先定义两个变量（变量不需要赋初值）。

```
package build
var (
	Date   string
	Commit string
)

```
然后，在代码中打印版本号的位置上将这些信息格式化输出，类似👇

``` go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/voidint/gbb/build"
)

var (
	// Version 版本号
	Version = "v0.0.1"
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
在项目目录*合适的地方*执行`gbb init`生成`gbb.json`文件。

合适的地方指哪些地方？一般规律是这样：

- 若使用的是`go build/install`工具编译代码(`gbb init`执行过程中填写的`tool`项对应的值)，那么这个**合适的地方**就是`main`方法所在目录。
- 若使用`gb`工具编译代码，那么这个**合适的地方**就是项目根目录。

按照`gbb init`的提示，逐步填写完信息并最终生成`gbb.json`文件。

``` shell
$ gbb init
This utility will walk you through creating a gbb.json file.
It only covers the most common items, and tries to guess sensible defaults.

Press ^C at any time to quit.
version: (0.0.1)
tool: (go_install) go_build
package: (main) github.com/voidint/gbb/build
variable: Date
value: {{.date}}
Do you want to continue?[y/n] y
variable: Commit
value: {{.gitCommit}}
Do you want to continue?[y/n] n
About to write to /Users/voidint/cloud/workspace/go/projects/src/github.com/voidint/gbb/gbb.json:

{
    "version": "0.0.1",
    "tool": "go install",
    "package": "github.com/voidint/gbb/build",
    "variables": [
        {
            "variable": "Date",
            "value": "{{.date}}"
        },
        {
            "variable": "Commit",
            "value": "{{.gitCommit}}"
        }
    ]
}

Is this ok?[y/n] y
```

关于`gbb.json`，请参见下文的[详细说明](https://github.com/voidint/gbb#gbbjson)。

### step2
在`gbb.json`文件所在目录编译（若目录下没有`gbb.json`文件，`gbb init`会被自动调用）。

```
$ gbb --debug
==> go build -ldflags  '-X "github.com/voidint/gbb/build.Date=2016-12-17T22:18:32+08:00" -X "github.com/voidint/gbb/build.Commit=db8b606cfc2b24a24e2e09acac24a52c47b68401"'
```
编译完后在目录下多出一个编译后的二进制文件，接着打印版本信息，看看是否实现我们设定的目标了。

```
$ ./gbb version
gbb version v0.0.1
date: 2016-12-17T22:18:32+08:00
commit: db8b606cfc2b24a24e2e09acac24a52c47b68401
```
😊

## gbb.json
`gbb.json`可以认为是`gbb`工具的配置文件，通过`gbb init`自动创建（感谢`npm init`）。通常它的格式是这样：

``` json
{
    "version": "0.0.1",
    "tool": "gb_build",
    "package": "build",
    "variables": [
        {
            "variable":"Date",
            "value":"{{.date}}"
        },
        {
            "variable":"Commit",
            "value":"{{.gitCommit}}"
        }
    ]
}
```

- `version`: 版本号。预留字段。
- `tool`: gbb实际调用的编译工具。已知的可选值包括：`go_build`、`go_install`、`gb_build`。注意：这个值不能包含空格[issue](https://github.com/voidint/gbb/issues/1)，因此暂时通过下划线`_`连接。
- `pakcage`: 包名，也就是定义`Date`、`Commit`这类变量的包全路径，如`github.com/voidint/gbb/build`。
- `variables`: 变量列表。列表中的每个元素都包含`variable`和`value`两个属性。
	- `variable`属性表示变量名，比如`Date`。
	- `value`属性表示变量值表达式，比如`{{.date}}`。内置变量表达式[列表](https://github.com/voidint/gbb/blob/master/variable/registry.go)。