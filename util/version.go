package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ErrSemanticVersion 非语义化版本号错误
var ErrSemanticVersion = errors.New("malformed semantic version")

var semVerReg = regexp.MustCompile("^\\d+\\.\\d+.\\d+$")

// VersionGreaterThan 对比v0的语义化版本号是否大于v1的语义化版本号。
// 若入参的版本号不符合语义化版本号规范，则返回ErrSemanticVersion。
// 关于语义化版本号内容，请参考http://semver.org/。
func VersionGreaterThan(v0, v1 string) (yes bool, err error) {
	if !semVerReg.MatchString(v0) || !semVerReg.MatchString(v1) {
		return false, ErrSemanticVersion
	}

	const sep = "."
	digitals0 := strings.Split(v0, sep)
	digitals1 := strings.Split(v1, sep)

	major0, _ := strconv.ParseInt(digitals0[0], 10, 64)
	major1, _ := strconv.ParseInt(digitals1[0], 10, 64)

	if major1 > major0 {
		return false, nil
	} else if major1 < major0 {
		return true, nil
	}

	minor0, _ := strconv.ParseInt(digitals0[1], 10, 64)
	minor1, _ := strconv.ParseInt(digitals1[1], 10, 64)

	if minor1 > minor0 {
		return false, nil
	} else if minor1 < minor0 {
		return true, nil
	}

	patch0, _ := strconv.ParseInt(digitals0[2], 10, 64)
	patch1, _ := strconv.ParseInt(digitals1[2], 10, 64)

	return patch0 > patch1, nil
}
