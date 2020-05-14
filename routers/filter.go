package routers

import (
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/smarteng/mindoc/conf"
	"github.com/smarteng/mindoc/models"
)

func init() {
	var FilterUser = func(ctx *context.Context) {
		_, ok := ctx.Input.Session(conf.LoginSessionName).(models.Member)
		if !ok {
			if ctx.Input.IsAjax() {
				jsonData := make(map[string]interface{}, 3)
				jsonData["errcode"] = 403
				jsonData["message"] = "请登录后再操作"

				returnJSON, _ := json.Marshal(jsonData)
				if _, err := ctx.ResponseWriter.Write(returnJSON); err != nil {
					panic("500")
				}
			} else {
				refer := url.PathEscape(conf.BaseUrl + ctx.Request.URL.RequestURI())
				ctx.Redirect(302, conf.URLFor("AccountController.Login")+"?url="+refer)
			}
		}
	}
	beego.InsertFilter("/manager", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/manager/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/setting", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/setting/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/book", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/book/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/api/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/manage/*", beego.BeforeRouter, FilterUser)

	var FinishRouter = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Add("MinDoc-Version", conf.VERSION)
		ctx.ResponseWriter.Header().Add("MinDoc-BuildTime", conf.BUILD_TIME)
		ctx.ResponseWriter.Header().Add("X-XSS-Protection", "1; mode=block")
	}

	var StartRouter = func(ctx *context.Context) {
		sessionID := ctx.Input.Cookie(beego.AppConfig.String("sessionname"))
		if sessionID != "" {
			//sessionId必须是数字字母组成，且最小32个字符，最大1024字符
			if ok, err := regexp.MatchString(`^[a-zA-Z0-9]{32,512}$`, sessionID); !ok || err != nil {
				panic("401")
			}
		}
	}
	beego.InsertFilter("/*", beego.BeforeStatic, StartRouter, false)
	beego.InsertFilter("/*", beego.BeforeRouter, FinishRouter, false)
}
