package controller

import (
    "github.com/gin-gonic/gin"

    "github.com/deatil/lakego-doak/lakego/pipeline"
    "github.com/deatil/lakego-doak/lakego/exception"
    "github.com/deatil/lakego-doak/lakego/filesystem"
    "github.com/deatil/lakego-doak/lakego/support/str"
    "github.com/deatil/lakego-doak/lakego/support/datebin"
    "github.com/deatil/lakego-doak/lakego/support/snowflake"
    "github.com/deatil/lakego-doak/lakego/support/cryptobin"
    "github.com/deatil/lakego-doak/lakego/support/hash"

    "github.com/deatil/lakego-doak-admin/admin/support/controller"
)

/**
 * 数据
 *
 * @create 2022-1-9
 * @author deatil
 */
type Data struct {
    controller.Base
}

/**
 * 信息
 */
func (this *Data) Index(ctx *gin.Context) {
    this.Fetch(ctx, "example::index", gin.H{
        "msg": "测试数据",
    })
}

/**
 * 信息2
 */
func (this *Data) Show(ctx *gin.Context) {
    this.Fetch(ctx, "example::show.index", map[string]interface{}{
        "msg": "测试数据",
    })
}

/**
 * Error 测试
 */
func (this *Data) Error(ctx *gin.Context) {
    // 报错测试
    data := []string{}
    exception.
        Try(func(){
            panic("exception error test")

            // exception.Throw("exception error test 222")
        }).
        Catch(func(e exception.Exception){
            // data = e.GetMessage()

            trace := e.GetTrace()
            for _, ev := range trace {
                data = append(data, ev.LongString())
            }
        })

    // 管道测试
    data2 := pipeline.NewPipeline().
        WithCarryCallback(func(carry interface{}) interface{} {
            return carry
        }).
        Send("开始的数据").
        Through(
            func(data interface{}, next pipeline.NextFunc) interface{} {
                old := data.(string)
                old = old + ", 第1次数据1"

                data2 := next(old)
                data2 = data2.(string) + ", 第1次数据2"

                return data2
            },
            func(data interface{}, next pipeline.NextFunc) interface{} {
                old := data.(string)
                old = old + ", 第2次数据1"

                data2 := next(old)
                data2 = data2.(string) + ", 第2次数据2"

                return data2
            },
            func(data interface{}, next pipeline.NextFunc) {
            },
            PipelineEx{},
        ).
        ThenReturn()

    // hub 测试
    hub := pipeline.NewHub()
    hub.Pipeline("hub", func(pipe *pipeline.Pipeline, object interface{}) interface{} {
        data := pipe.
            Send(object).
            Through(
                func(data interface{}, next pipeline.NextFunc) interface{} {
                    old := data.(string)
                    old = old + ", 第1次数据1"

                    data2 := next(old)
                    data2 = data2.(string) + ", 第1次数据2"

                    return data2
                },
            ).
            ThenReturn()

        return data
    })
    data3 := hub.Pipe("hub 测试", "hub")

    // 雪花算法
    snowflakeId, _ := snowflake.Make(5)

    // 文件管理
    fs := filesystem.New()
    img := "./runtime/test/13123/123321.txt"
    // img2 := "./runtime/test/13123"
    // txt := "./runtime/test/1/data.txt"
    fsInfo := fs.IsWritable(img)

    // 字符处理
    strData := "t_ydfd_ydf"
    newStrData := str.LowerCamel(strData)

    // 时间
    date := datebin.
        Now().
        ToDatetimeString()
    date2 := datebin.
        Parse("2032-03-15 12:06:17").
        ToDatetimeString()

    // 加密测试
    cypt := cryptobin.
        FromString("test-pass").
        SetKey("dfertf12dfertf12dfertf12").
        Encrypt().
        ToHexString()
    cyptde := cryptobin.
        FromHexString("deec911b634a2435bb7bf8a6fb1f9eae").
        SetKey("dfertf12dfertf12dfertf12").
        Decrypt().
        ToString()

    // 签名
    rsa := cryptobin.NewRsa()
    rsaPriKey, _ := rsa.MakePassPKCS8PrvKey(1024, "123", "DES")
    rsaPubKey, _ := rsa.MakePKCS8PubKeyFromPassPKCS8PrvKey(rsaPriKey, "123")

    // 签名
    hashData := hash.MD5SHA1("123")

    this.SuccessWithData(ctx, "Error 测试", gin.H{
        "error": data,
        "data2": data2,
        "data3": data3,
        "snowflakeId": snowflakeId,

        "img": img,
        "fsInfo": fsInfo,

        "str": newStrData,

        "date": date,
        "date2": date2,

        "cypt": cypt,
        "cyptde": cyptde,

        "rsaPriKey": string(rsaPriKey),
        "rsaPubKey": string(rsaPubKey),

        "hashData": hashData,
    })
}

/* ======================== */

// 管道测试
type PipelineEx struct {}

func (this PipelineEx) Handle(data interface{}, next pipeline.NextFunc) interface{} {
    old := data.(string)

    old = old + ", struct 数据开始"

    data2 := next(old)

    data2 = data2.(string) + ", struct 数据结束"

    return data2
}
