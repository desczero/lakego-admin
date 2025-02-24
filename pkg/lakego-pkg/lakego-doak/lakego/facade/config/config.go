package config

import (
    "github.com/deatil/lakego-doak/lakego/register"
    "github.com/deatil/lakego-doak/lakego/path"

    "github.com/deatil/lakego-doak/lakego/config"
    "github.com/deatil/lakego-doak/lakego/config/interfaces"
    viper_adapter "github.com/deatil/lakego-doak/lakego/config/adapter/viper"
)

var (
    // 默认驱动
    defaultAdapter = "viper"

    // 配置目录
    defaultConfigPath = "{root}/config"
)

// 初始化
func init() {
    // 注册默认
    Register()
}

// 配置别名
type Config = config.Config

/**
 * 配置
 *
 * @create 2021-9-25
 * @author deatil
 */
func New(name ...string) *config.Config {
    adapter := defaultAdapter

    if len(name) > 0 {
        return NewConfig(adapter).WithFile(name[0])
    }

    return NewConfig(adapter)
}

// 实例化
func NewWithAdapter(name string, adapter string) *config.Config {
    return NewConfig(adapter).WithFile(name)
}

// 配置
func NewConfig(name string, once ...bool) *config.Config {
    adapter := register.
        NewManagerWithPrefix("config").
        GetRegister(name, nil, once...)
    if adapter == nil {
        panic("配置驱动[" + name + "]没有被注册")
    }

    conf := &config.Config{}
    conf.WithAdapter(adapter.(interfaces.Adapter))

    return conf
}

// 设置默认驱动
func SetAdapter(name string) {
    defaultAdapter = name
}

// 设置配置路径
func SetConfigPath(cfgPath string) {
    defaultConfigPath = cfgPath
}

// 注册磁盘
func Register() {
    // 注册可用驱动
    register.
        NewManagerWithPrefix("config").
        Register("viper", func(conf map[string]any) any {
            adapter := viper_adapter.New()

            // 配置文件夹
            configPath := path.FormatPath(defaultConfigPath)

            // 设置 env 前缀
            adapter.SetEnvPrefix("LAKEGO")
            adapter.AutomaticEnv()
            adapter.WithPath(configPath)

            return adapter
        })
}

