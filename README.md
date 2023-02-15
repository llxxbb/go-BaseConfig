# BaseConfig

该工具是 [viper](https://github.com/spf13/viper) 的扩展，提供了配置文件合并的能力。

## 预设配置项

缺省提供下面的配置项，

```yml
prj:
  name: back-normal
  version: v0.0.1

port: 80

log.root: /web/logs

gin.release: true
```

配置会写入到 `BaseConfig` 结构中，并会依据配置值自动设置下面的字段

- Host：运行该项目实例的主机IP 

- WorkPath

- LogPath：项目的日志输出路径

## 使用说明

- **搜索路径**：配置文件的搜索路径为当前项目执行目录。

- **缺省配置文件**：放于 config_default.yml 文件中，嵌入的请设置 config.FDefault 变量

- **env 环境变量**： env 的缺省值为 product，用于确定要加载的配置文件如，如果 env = dev 则加载 config_dev.yml

- **优先级**：环境变量>特定环境配置>缺省配置

- **打印配置**：调用`Print()`方法

## 代码示例

```go
func TestConfig_default(t *testing.T) {
    c := BaseConfig{}
    // 填充配置
    FillConfig(&c, &c)
    // 验证配置
    assert.Equal(t, "back-normal", c.ProjectName)
    assert.Equal(t, "v0.0.1", c.ProjectVersion)
    assert.Equal(t, "product", c.Env)
    wd, _ := os.Getwd()
    assert.Equal(t, wd, c.WorkPath)
    assert.Equal(t, true, len(c.Host) > 0)
    println(c.LogPath)
    assert.Equal(t, true, len(c.LogPath) > 0)
}
```

## 扩展自己的配置

```go
// 定义自己的配置，增加 Mysql 的配置
type MyConfig struct {
    BaseConfig
    Mysql   mysql.Config
    MaxOpen int
    MaxIdle int
}

// 实现 ConfigI 接口，用于配置文件和结构体字段的映射
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
```
