## Changelog

### 0.1.1 - 2016/12/24
- 支持通过`gbb init`初始化配置信息并生成`gbb.json`配置文件。
- 支持在项目根目录下，一键编译所有入口源代码文件，并生成一个或者多个可执行二进制文件。
- 支持调用`gb`或者`go build/install`，并为编译生成的可执行文件提供丰富的版本信息中，包括但不限于：`编译时间`、`源代码版本控制commit`等。

### 0.2.0 - 2016/12/30
- `gbb.json`中的配置项——`package`和`variables`由必选项改为可选项。其中，在`variables`选项为空的情况下，实际在调用编译工具编译时不再加上形如`-ldflags '-X "xxx.yyy=zzz"'`的参数。[#8](https://github.com/voidint/gbb/issues/8)
- 若程序版本号与`gbb.json`中的`version`值不一致，就会强制重新生成`gbb.json`文件。