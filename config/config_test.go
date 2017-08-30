package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoad(t *testing.T) {
	Convey("从磁盘文件加载配置信息", t, func() {
		filename := filepath.Join(os.TempDir(), "gbb.json")
		Convey("预期配置文件存在，实际文件存在", func() {
			Convey("配置文件内容为合法JSON", func() {
				c := Config{
					Version:    "0.1.1",
					Tool:       "go install",
					Importpath: "github.com/voidint/gbb/build",
					Variables: []Variable{
						{Variable: "Date", Value: "{{.Date}}"},
						{Variable: "Commit", Value: "{{.GitCommit}}"},
					},
				}
				err := Save(&c, filename)
				So(err, ShouldBeNil)
				defer os.Remove(filename)

				loadC, err := Load(filename)
				So(err, ShouldBeNil)
				So(loadC, ShouldNotBeNil)
				So(loadC.Version, ShouldEqual, c.Version)
				So(loadC.Tool, ShouldEqual, c.Tool)
				So(loadC.Importpath, ShouldEqual, c.Importpath)
				So(len(loadC.Variables), ShouldEqual, len(c.Variables))
				So(loadC.Variables[0].Variable, ShouldEqual, c.Variables[0].Variable)
				So(loadC.Variables[0].Value, ShouldEqual, c.Variables[0].Value)
			})

			Convey("配置文件内容非合法JSON", func() {
				So(ioutil.WriteFile(filename, []byte("hello world"), 0666), ShouldBeNil)
				defer os.Remove(filename)
				loadC, err := Load(filename)
				So(err, ShouldNotBeNil)
				So(loadC, ShouldBeNil)
			})

			Convey("读取配置文件时报错", func() {
				var ErrRead = errors.New("failed to read file ")
				monkey.Patch(ioutil.ReadFile, func(filename string) ([]byte, error) {
					return nil, ErrRead
				})
				defer monkey.Unpatch(ioutil.ReadFile)
				loadC, err := Load(filename)
				So(err, ShouldNotBeNil)
				So(err, ShouldEqual, ErrRead)
				So(loadC, ShouldBeNil)
			})
		})

		Convey("预期配置文件存在，实际文件不存在", func() {
			os.Remove(filename)
			loadC, err := Load(filename)
			So(err, ShouldNotBeNil)
			So(loadC, ShouldBeNil)
		})
	})
}

func TestSave(t *testing.T) {
	Convey("将配置对象持久化到磁盘文件", t, func() {
		filename := filepath.Join(os.TempDir(), "gbb.json")
		Convey("配置对象和磁盘文件都合法", func() {
			c := Config{
				Version:    "0.1.1",
				Tool:       "go install",
				Importpath: "github.com/voidint/gbb/build",
				Variables: []Variable{
					{Variable: "Date", Value: "{{.Date}}"},
					{Variable: "Commit", Value: "{{.GitCommit}}"},
				},
			}
			So(Save(&c, filename), ShouldBeNil)
			defer os.Remove(filename)

			loadC, err := Load(filename)
			So(err, ShouldBeNil)
			So(loadC, ShouldNotBeNil)
			So(loadC.Version, ShouldEqual, c.Version)
			So(loadC.Tool, ShouldEqual, c.Tool)
			So(loadC.Importpath, ShouldEqual, c.Importpath)
			So(len(loadC.Variables), ShouldEqual, len(c.Variables))
			So(loadC.Variables[0].Variable, ShouldEqual, c.Variables[0].Variable)
			So(loadC.Variables[0].Value, ShouldEqual, c.Variables[0].Value)
		})
		Convey("配置对象为nil", func() {
			So(Save(nil, filename), ShouldBeNil)
		})

		Convey("JSON序列化报错", func() {
			var ErrMarshal = errors.New("marshal error")
			monkey.Patch(json.MarshalIndent, func(v interface{}, prefix, indent string) ([]byte, error) {
				return nil, ErrMarshal
			})
			defer monkey.Unpatch(json.MarshalIndent)

			So(Save(new(Config), "gbb.json"), ShouldEqual, ErrMarshal)
		})
	})
}
