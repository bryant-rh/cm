package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mattn/go-isatty"
	"github.com/yitter/idgenerator-go/idgen"
)

func SnowId() (int64, error) {

	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	// Generate a snowflake ID.
	id := node.Generate()
	return id.Int64(), nil
}

type IGenID interface {
	ID() int64
}

type IDMgr struct {
	idGen *idgen.DefaultIdGenerator
}

func NewIdMgr(node uint16) (*IDMgr, error) {
	opts := idgen.NewIdGeneratorOptions(node)
	idGenerator := idgen.NewDefaultIdGenerator(opts)
	return &IDMgr{
		idGen: idGenerator,
	}, nil
}

func (c *IDMgr) ID() int64 {
	return c.idGen.NewLong()
}

func StringtoInt64(s string) int64 {
	num, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return int64(num)
}

func Int64toString(i int64) string {
	return strconv.FormatInt(i, 10)
}
func InttoString(i int) string {
	return strconv.Itoa(i)
}

func StringtoSlice(s string) []string {
	if s == "" {
		return []string{}
	} else {
		return strings.Split(s, ",")

	}
}

func SlicetoString(s []string) string {
	return strings.Join(s, ",")
}

func SliceDelOne(a string, b []string) []string {
	target := b[:0]
	for _, item := range b {
		if item != a {
			target = append(target, item)
		}
	}
	return target

}

type Res struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func ReturnMsg(ctx *gin.Context, code int, data interface{}, msg string) {
	ctx.JSON(code, Res{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

//
func StructToSlice(f interface{}) []string {
	v := reflect.ValueOf(f)
	ss := make([]string, v.NumField())
	for i := range ss {
		ss[i] = fmt.Sprintf("%v", v.Field(i))
	}
	return ss
}

func IsTerminal(f *os.File) bool {
	return isatty.IsTerminal(f.Fd()) || isatty.IsCygwinTerminal(f.Fd())
}
func RedColor(s string) string {
	return fmt.Sprintf("%s", aurora.Red(s))
}

func YellowColor(s string) string {
	return fmt.Sprintf("%s", aurora.Yellow(s))
}

func GreenColor(s string) string {
	return fmt.Sprintf("%s", aurora.Green(s))
}

func WriteToFile(filename string, v interface{}, marshal func(v interface{}) ([]byte, error)) error {
	bytes, _ := marshal(v)
	dir := filepath.Dir(filename)
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, bytes, os.ModePerm)
}
