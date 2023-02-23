/* Usage:
- Define you own config struct
- Make `Config` as a field of you config
- implement interface `ConfigI`
- call `NewConfig()`
*/

package config

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type BaseConfig struct {
	ProjectName    string // 项目名称
	ProjectVersion string // 项目本身的版本信息
	Env            string // 部署环境
	Port           string // 对外服务的端口号
	GinRelease     bool   // gin 是否以 release 模式工作
	LogRoot        string // 保存日志的根路径

	// 注意：下面的配置项目在运行时自动设置，无需配置。
	Host     string // 实例部署位置
	WorkPath string // 项目启动所在的目录
	LogPath  string // 日志输出的位置，由 log.root 主机IP 项目名 等组成
}

type ConfigI interface {
	AppendFieldMap(map[string]string)
	Print()
}

const (
	KEY_ENV     = "env"
	VAL_PRODUCT = "product"
)

var FDefault []byte

// 用于反射，将配置文件中的配置项映射到 `Config` 对象的属性上
var fieldMap = map[string]string{
	"prj.name":    "ProjectName",
	"prj.version": "ProjectVersion",
	"port":        "Port",
	"gin.release": "GinRelease",
	"log.root":    "LogRoot",
}

const (
	_fileDefault  = "default"
	_fileName     = "config"
	_fileType     = "yaml"
	_nameSplitter = "_"
)

func FillConfig(cus ConfigI, base *BaseConfig) {
	cus.AppendFieldMap(fieldMap)
	// 可以从环境变量中取值
	viper.SetDefault(KEY_ENV, VAL_PRODUCT)
	viper.AutomaticEnv()
	base.Env = viper.GetString(KEY_ENV)

	// 读取缺省配置文件
	mergeFile(cus, _fileDefault)

	// 读取 profile 对应的配置文件
	mergeFile(cus, base.Env)

	// 读取环境变量
	mergeFile(cus, "")
	// 设置工作目录、日志目录等
	base.setWdAndLogPath()
}

func (c *BaseConfig) setWdAndLogPath() {
	var e error
	c.WorkPath, e = os.Getwd()
	if e != nil {
		panic(e)
	}
	c.Host, e = getOutBoundIP()
	if e != nil {
		panic(e)
	}
	c.LogPath = c.LogRoot + "/" + c.Host + "-" + c.ProjectName
}

// 读取并合并配置项
// `path` 为 "" 则从环境变量中读取
func mergeFile(c ConfigI, profile string) {
	if profile != "" {
		path := _fileName + _nameSplitter + profile
		viper.SetConfigName(path)      // name of config file (without extension)
		viper.SetConfigType(_fileType) // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath(".")       // look for config in the working directory
		viper.AddConfigPath("../cmd")  // used for unit test
		err := viper.ReadInConfig()    // Find and read the config file
		if err != nil {                // Handle errors reading the config file
			if profile == _fileDefault {
				if viper.ReadConfig(bytes.NewReader(FDefault)) != nil {
					panic(fmt.Errorf("load config_default.yml from embed: %w", err))
				}
				fmt.Printf("%s loaded from embed\n", path)
			} else {
				panic(fmt.Errorf("fatal error: read config file: %s, %w", path, err))
			}
		} else {
			fmt.Printf("%s loaded\n", path)
		}
	}

	// 利用反射进行赋值
	cfgType := reflect.ValueOf(c).Elem()
	for mKey, mVal := range fieldMap {
		fileVal := viper.GetString(mKey)
		if fileVal == "" {
			continue
		}
		var rV reflect.Value
		var err error
		field := findField(&cfgType, mVal)
		typeName := field.Type().Name()
		switch typeName {
		case "string":
			rV = reflect.ValueOf(fileVal)
		case "bool":
			rtn, e := strconv.ParseBool(fileVal)
			if e == nil {
				rV = reflect.ValueOf(rtn)
			} else {
				err = e
			}
		case "int":
			rtn, e := strconv.ParseInt(fileVal, 0, 64)
			if e == nil {
				rV = reflect.ValueOf(int(rtn))
			} else {
				err = e
			}
		case "float32":
			rtn, e := strconv.ParseFloat(fileVal, 32)
			if e == nil {
				rV = reflect.ValueOf(float32(rtn))
			} else {
				err = e
			}
		case "float64":
			rtn, e := strconv.ParseFloat(fileVal, 64)
			if e == nil {
				rV = reflect.ValueOf(float64(rtn))
			} else {
				err = e
			}
		case "Duration":
			rtn, e := time.ParseDuration(fileVal)
			if e == nil {
				rV = reflect.ValueOf(rtn)
			} else {
				err = e
			}
		default:
			panic(fmt.Sprintf("config item: %s, unhandled type.", mKey))
		}
		if err != nil {
			panic(fmt.Sprintf("config item: %s, value type error. %v", mKey, err))
		}
		field.Set(rV)
	}
}

func findField(parent *reflect.Value, field string) *reflect.Value {
	fs := strings.Split(field, ".")
	rtn := parent.FieldByName(fs[0])
	if len(fs) == 1 {
		return &rtn
	} else {
		return findField(&rtn, field[len(fs[0])+1:])
	}
}

func (c *BaseConfig) AppendFieldMap(fm map[string]string) {

}

func (c *BaseConfig) Print() {
	zap.L().Info("------------ project info ------------")
	zap.L().Info("-- ", zap.String("ProjectName", c.ProjectName))
	zap.L().Info("-- ", zap.String("ProjectVersion", c.ProjectVersion))
	zap.L().Info("-- ", zap.Bool("GinRelease", c.GinRelease))
	zap.L().Info("------------ endpoint info ------------")
	zap.L().Info("-- ", zap.String("Env", c.Env))
	zap.L().Info("-- ", zap.String("Host", c.Host))
	zap.L().Info("-- ", zap.String("Port", c.Port))
	zap.L().Info("------------ path info ------------")
	zap.L().Info("-- ", zap.String("WorkPath", c.WorkPath))
	zap.L().Info("-- ", zap.String("LogPath", c.LogPath))
}

func getOutBoundIP() (ip string, err error) {
	out := "114.114.114.114:53"
	fmt.Printf("use %s test local out bound,", out)
	conn, err := net.Dial("udp", out)
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(" the result is: " + localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}
