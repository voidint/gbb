## Changelog

### 0.2.0 - 2016/12/30
- `gbb.json`中的配置项——`package`和`variables`由必选项改为可选项。其中，在`variables`选项为空的情况下，实际在调用编译工具编译时不再加上形如`-ldflags '-X "xxx.yyy=zzz"'`的参数。[#8](https://github.com/voidint/gbb/issues/8)
- 若程序版本号与`gbb.json`中的`version`值不一致，就会强制重新生成`gbb.json`文件。