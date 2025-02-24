package logger

import (
    "log"
    "strings"

    "github.com/deatil/lakego-doak/lakego/register"
    "github.com/deatil/lakego-doak/lakego/facade/config"
    "github.com/deatil/lakego-doak/lakego/logger"
    "github.com/deatil/lakego-doak/lakego/logger/interfaces"
    logrusDriver "github.com/deatil/lakego-doak/lakego/logger/driver/logrus"
)

// 初始化
func init() {
    // 注册默认
    Register()
}

/**
 * 验证码
 *
 * @create 2021-10-12
 * @author deatil
 */
func New(once ...bool) *logger.Logger {
    // 默认驱动
    driver := GetDefaultDriver()

    return NewLogger(driver, once...)
}

// 验证码
func NewLogger(driverName string, once ...bool) *logger.Logger {
    // 配置
    conf := config.New("logger")

    // 驱动列表
    drivers := conf.GetStringMap("drivers")

    // 转为小写
    driverName = strings.ToLower(driverName)

    // 获取配置
    driverConfig, ok := drivers[driverName]
    if !ok {
        log.Print("日志驱动[" + driverName + "]配置不存在")
    }

    // 驱动配置
    driverConf := driverConfig.(map[string]any)

    driverType := driverConf["type"].(string)
    driver := register.
        NewManagerWithPrefix("logger").
        GetRegister(driverType, driverConf, once...)
    if driver == nil {
        log.Print("日志驱动[" + driverType + "]没有被注册")
    }

    return logger.New(driver.(interfaces.Driver))
}

// 自定义数据
// import "github.com/deatil/lakego-doak/lakego/facade/logger"
// logger.LogrusWithField(logger.New(), "system", "lakego").Info("logger test")
func LogrusWithField(log *logger.Logger, key string, value any) *logrusDriver.Entry {
    return log.WithField(key, value).(*logrusDriver.Entry)
}

// 批量自定义数据
func LogrusWithFields(log *logger.Logger, fields map[string]any) *logrusDriver.Entry {
    return log.WithFields(fields).(*logrusDriver.Entry)
}

// 默认驱动
func GetDefaultDriver() string {
    return config.New("logger").GetString("default")
}

// 注册
func Register() {
    // 注册驱动
    register.
        NewManagerWithPrefix("logger").
        RegisterMany(map[string]func(map[string]any) any {
            // logrus 日志
            "logrus": func(conf map[string]any) any {
                driver := logrusDriver.New()

                driver.WithConfig(conf)

                return driver
            },
        })
}
