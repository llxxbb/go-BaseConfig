package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestConfig_default(t *testing.T) {
	c := BaseConfig{}
	FillConfig(&c, &c)
	assert.Equal(t, "back-normal", c.ProjectName)
	assert.Equal(t, "v0.0.1", c.ProjectVersion)
	assert.Equal(t, 123, c.ProjectId)
	assert.Equal(t, "product", c.Env)
	wd, _ := os.Getwd()
	assert.Equal(t, wd, c.WorkPath)
	assert.Equal(t, true, len(c.Host) > 0)
	println(c.LogPath)
	assert.Equal(t, true, len(c.LogPath) > 0)
}

func TestFindField(t *testing.T) {
	c := MyConfig{}
	c.Mysql.DBName = "lxb"
	root := reflect.ValueOf(&c).Elem()
	to := findField(&root, "Mysql.DBName")
	assert.Equal(t, "lxb", to.String())
}

// 定义自己的配置
type MyConfig struct {
	BaseConfig
	Mysql       mysql.Config
	MaxOpen     int
	MaxIdle     int
	TypeInt     int
	TypeFloat32 float32
	TypeFloat64 float64
	TypeInt8    int8
	TypeInt16   int16
	TypeInt32   int32
	TypeInt64   int64
	TypeUInt    uint
	TypeUInt8   uint8
	TypeUInt16  uint16
	TypeUInt32  uint32
	TypeUInt64  uint64
}

// 实现 ConfigI 接口
func (c *MyConfig) AppendFieldMap(fMap map[string]string) {
	fMap["mysql.user"] = "Mysql.User"
	fMap["mysql.password"] = "Mysql.Passwd"
	fMap["mysql.address"] = "Mysql.Addr"
	fMap["mysql.db"] = "Mysql.DBName"
	fMap["mysql.conns.timeout"] = "Mysql.Timeout"
	fMap["mysql.conns.readTimeout"] = "Mysql.ReadTimeout"
	fMap["mysql.conns.maxOpen"] = "MaxOpen"
	fMap["mysql.conns.maxIdle"] = "MaxIdle"
	fMap["type.int"] = "TypeInt"
	fMap["type.float32"] = "TypeFloat32"
	fMap["type.float64"] = "TypeFloat64"
	fMap["type.int8"] = "TypeInt8"
	fMap["type.int16"] = "TypeInt16"
	fMap["type.int32"] = "TypeInt32"
	fMap["type.int64"] = "TypeInt64"
	fMap["type.uint"] = "TypeUInt"
	fMap["type.uint8"] = "TypeUInt8"
	fMap["type.uint16"] = "TypeUInt16"
	fMap["type.uint32"] = "TypeUInt32"
	fMap["type.uint64"] = "TypeUInt64"
}

func TestCustomConfig(t *testing.T) {
	c := MyConfig{}
	FillConfig(&c, &c.BaseConfig)
	assert.Equal(t, "user", c.Mysql.User)
	assert.Equal(t, "password", c.Mysql.Passwd)
	assert.Equal(t, "localhost:3306", c.Mysql.Addr)
	assert.Equal(t, "testdb", c.Mysql.DBName)
	assert.Equal(t, 40, c.MaxOpen)
	assert.Equal(t, 2, c.MaxIdle)
	assert.Equal(t, -5, c.TypeInt)
	assert.Equal(t, float32(6.6), c.TypeFloat32)
	assert.Equal(t, 8.8, c.TypeFloat64)
	assert.Equal(t, int8(-8), c.TypeInt8)
	assert.Equal(t, int16(-16), c.TypeInt16)
	assert.Equal(t, int32(-32), c.TypeInt32)
	assert.Equal(t, int64(-64), c.TypeInt64)
	assert.Equal(t, uint(5), c.TypeUInt)
	assert.Equal(t, uint8(8), c.TypeUInt8)
	assert.Equal(t, uint16(16), c.TypeUInt16)
	assert.Equal(t, uint32(32), c.TypeUInt32)
	assert.Equal(t, uint64(64), c.TypeUInt64)
}
