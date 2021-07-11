package permission

import (
    "strings"
    "github.com/gin-gonic/gin"

    "lakego-admin/lakego/config"
    "lakego-admin/lakego/lake"
    "lakego-admin/lakego/http/code"
    "lakego-admin/lakego/http/response"
    "lakego-admin/lakego/facade/casbin"

    "lakego-admin/admin/auth/admin"
)

// 权限检测
func Permission() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        if !shouldPassThrough(ctx) {
            // 权限检测
            permissionCheck(ctx)
        }

        ctx.Next()
    }
}

// 权限检测
func permissionCheck(ctx *gin.Context) bool {
    if checkSuperAdmin(ctx) {
        return true
    }

    requestPath := ctx.Request.URL.String()
    method := strings.ToUpper(ctx.Request.Method)

    adminId, _ := ctx.Get("admin_id")

    c := casbin.New()
    ok, err := c.Enforce(adminId.(string), requestPath, method)

    if err != nil {
        response.Error(ctx, code.AuthError, "你没有访问权限")
        return false
    } else if !ok {
        response.Error(ctx, code.AuthError, "你没有访问权限")
        return false
    }

    return true
}

// 超级管理员检测
func checkSuperAdmin(ctx *gin.Context) bool {
    adminInfo, _ := ctx.Get("admin")

    if adminInfo == nil {
        response.Error(ctx, code.AuthError, "你没有访问权限")
        return false
    }

    isSuperAdministrator := adminInfo.(*admin.Admin).IsSuperAdministrator()
    if isSuperAdministrator {
        return true
    }

    return false
}

// 过滤
func shouldPassThrough(ctx *gin.Context) bool {
    // 默认过滤
    excepts := []string{
        "GET:" + lake.AdminUrl("passport/captcha"),
        "POST:" + lake.AdminUrl("passport/login"),
        "PUT:" + lake.AdminUrl("passport/refresh-token"),
    }
    for _, except := range excepts {
        if lake.MatchPath(ctx, except, "") {
            return true
        }
    }

    // 自定义权限过滤
    authenticateExcepts := config.New("auth").GetStringSlice("Auth.PermissionExcepts")
    for _, ae := range authenticateExcepts {
        newStr := strings.Split(ae, ":")
        if len(newStr) == 2 {
            newUtl := newStr[0] + ":" + lake.AdminUrl(newStr[1])
            if lake.MatchPath(ctx, newUtl, "") {
                return true
            }
        }
    }

    return false
}
