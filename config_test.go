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
	Mysql   mysql.Config
	MaxOpen int
	MaxIdle int
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
}
