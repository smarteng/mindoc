// package conf 为配置相关.
package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
)

// 登录用户的Session名
const LoginSessionName = "LoginSessionName"

const CaptchaSessionName = "__captcha__"

//允许用户名中出现点号
const RegexpAccount = `^[a-zA-Z][a-zA-Z0-9\.-]{2,50}$`

// PageSize 默认分页条数.
const PageSize = 10

// 用户权限
const (
	// 超级管理员.
	MemberSuperRole SystemRole = iota
	//普通管理员.
	MemberAdminRole
	//普通用户.
	MemberGeneralRole
)

//系统角色
type SystemRole int

const (
	// 创始人.
	BookFounder BookRole = iota
	//管理者
	BookAdmin
	//编辑者.
	BookEditor
	//观察者
	BookObserver
)

//项目角色
type BookRole int

const (
	LoggerOperate   = "operate"
	LoggerSystem    = "system"
	LoggerException = "exception"
	LoggerDocument  = "document"
)
const (
	//本地账户校验
	AuthMethodLocal = "local"
	//LDAP用户校验
	AuthMethodLDAP = "ldap"
)

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

var (
	ConfigurationFile = "./conf/app.conf"
	WorkingDirectory  = "./"
	LogFile           = "./runtime/logs"
	BaseUrl           = ""
	AutoLoadDelay     = 0
)

func init() {
	VERSION = beego.AppConfig.String("version")
}

// app_key
func GetAppKey() string {
	return beego.AppConfig.DefaultString("app_key", "mindoc")
}

func GetDatabasePrefix() string {
	return beego.AppConfig.DefaultString("db_prefix", "md_")
}

//获取默认头像
func GetDefaultAvatar() string {
	return URLForWithCdnImage(beego.AppConfig.DefaultString("avatar", "/images/headimgurl.jpg"))
}

//获取阅读令牌长度.
func GetTokenSize() int {
	return beego.AppConfig.DefaultInt("token_size", 12)
}

//获取默认文档封面.
func GetDefaultCover() string {
	return URLForWithCdnImage(beego.AppConfig.DefaultString("cover", "/images/book.jpg"))
}

//获取允许的商城文件的类型.
func GetUploadFileExt() []string {
	ext := beego.AppConfig.DefaultString("upload_file_ext", "png|jpg|jpeg|gif|txt|doc|docx|pdf")
	temp := strings.Split(ext, "|")
	exts := make([]string, len(temp))
	i := 0
	for _, item := range temp {
		if item != "" {
			exts[i] = item
			i++
		}
	}
	return exts
}

// 获取上传文件允许的最大值
func GetUploadFileSize() int64 {
	size := beego.AppConfig.DefaultString("upload_file_size", "0")

	if strings.HasSuffix(size, "MB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "GB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024 * 1024 * 1024
		}
	}
	if strings.HasSuffix(size, "KB") {
		if s, e := strconv.ParseInt(size[0:len(size)-2], 10, 64); e == nil {
			return s * 1024
		}
	}
	if s, e := strconv.ParseInt(size, 10, 64); e == nil {
		return s * 1024
	}
	return 0
}

//是否启用导出
func GetEnableExport() bool {
	return beego.AppConfig.DefaultBool("enable_export", true)
}

//同一项目导出线程的并发数
func GetExportProcessNum() int {
	exportProcessNum := beego.AppConfig.DefaultInt("export_process_num", 1)

	if exportProcessNum <= 0 || exportProcessNum > 4 {
		exportProcessNum = 1
	}
	return exportProcessNum
}

//导出项目队列的并发数量
func GetExportLimitNum() int {
	exportLimitNum := beego.AppConfig.DefaultInt("export_limit_num", 1)

	if exportLimitNum < 0 {
		exportLimitNum = 1
	}
	return exportLimitNum
}

//等待导出队列的长度
func GetExportQueueLimitNum() int {
	exportQueueLimitNum := beego.AppConfig.DefaultInt("export_queue_limit_num", 10)

	if exportQueueLimitNum <= 0 {
		exportQueueLimitNum = 100
	}
	return exportQueueLimitNum
}

//默认导出项目的缓存目录
func GetExportOutputPath() string {
	exportOutputPath := filepath.Join(beego.AppConfig.DefaultString("export_output_path", filepath.Join(WorkingDirectory, "cache")), "books")

	return exportOutputPath
}

//判断是否是允许商城的文件类型.
func IsAllowUploadFileExt(ext string) bool {
	if strings.HasPrefix(ext, ".") {
		ext = string(ext[1:])
	}
	exts := GetUploadFileExt()
	for _, item := range exts {
		if item == "*" {
			return true
		}
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

//重写生成URL的方法，加上完整的域名
func URLFor(endpoint string, values ...interface{}) string {
	baseURL := beego.AppConfig.DefaultString("baseurl", "")
	pathURL := beego.URLFor(endpoint, values...)

	if baseURL == "" {
		baseURL = BaseUrl
	}
	if strings.HasPrefix(pathURL, "http://") {
		return pathURL
	}
	if strings.HasPrefix(pathURL, "/") && strings.HasSuffix(baseURL, "/") {
		return baseURL + pathURL[1:]
	}
	if !strings.HasPrefix(pathURL, "/") && !strings.HasSuffix(baseURL, "/") {
		return baseURL + "/" + pathURL
	}
	return baseURL + beego.URLFor(endpoint, values...)
}

func URLForNotHost(endpoint string, values ...interface{}) string {
	baseURL := beego.AppConfig.DefaultString("baseurl", "")
	pathURL := beego.URLFor(endpoint, values...)

	if baseURL == "" {
		baseURL = "/"
	}
	if strings.HasPrefix(pathURL, "http://") {
		return pathURL
	}
	if strings.HasPrefix(pathURL, "/") && strings.HasSuffix(baseURL, "/") {
		return baseURL + pathURL[1:]
	}
	if !strings.HasPrefix(pathURL, "/") && !strings.HasSuffix(baseURL, "/") {
		return baseURL + "/" + pathURL
	}
	return baseURL + beego.URLFor(endpoint, values...)
}

// 注册图片的image方法
func URLForWithCdnImage(p string) string {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "//") {
		return p
	}
	cdn := beego.AppConfig.DefaultString("cdnimg", "")
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseURL := beego.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseURL, "/") {
			return baseURL + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseURL, "/") {
			return baseURL + "/" + p
		}
		return baseURL + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

// 注册ccs方法
func URLForWithCdnCSS(p string, v ...string) string {
	cdn := beego.AppConfig.DefaultString("cdncss", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "//") {
		return p
	}
	filePath := WorkingDir(p)

	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseUrl := beego.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseUrl, "/") {
			return baseUrl + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseUrl, "/") {
			return baseUrl + "/" + p
		}
		return baseUrl + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

// 注册js模板方法
func URLForWithCdnJs(p string, v ...string) string {
	cdn := beego.AppConfig.DefaultString("cdnjs", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "//") {
		return p
	}
	filePath := WorkingDir(p)
	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}
	//如果没有设置cdn，则使用baseURL拼接
	if cdn == "" {
		baseURL := beego.AppConfig.DefaultString("baseurl", "/")

		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseURL, "/") {
			return baseURL + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseURL, "/") {
			return baseURL + "/" + p
		}
		return baseURL + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
		return cdn + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
		return cdn + "/" + p
	}
	return cdn + p
}

// 注册静态文件模板方法
func URLForWithStatic(p string, v ...string) string {
	static := beego.AppConfig.DefaultString("static", "")
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "//") {
		return p
	}
	filePath := WorkingDir(p)
	if f, err := os.Stat(filePath); err == nil && !strings.Contains(p, "?") && len(v) > 0 && v[0] == "version" {
		p = p + fmt.Sprintf("?v=%s", f.ModTime().Format("20060102150405"))
	}
	//如果没有设置cdn，则使用baseURL拼接
	if static == "" {
		baseURL := beego.AppConfig.DefaultString("baseurl", "/")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(baseURL, "/") {
			return baseURL + p[1:]
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(baseURL, "/") {
			return baseURL + "/" + p
		}
		return baseURL + p
	}
	if strings.HasPrefix(p, "/") && strings.HasSuffix(static, "/") {
		return static + string(p[1:])
	}
	if !strings.HasPrefix(p, "/") && !strings.HasSuffix(static, "/") {
		return static + "/" + p
	}
	return static + p
}

// 获取当前工作目录
func WorkingDir(elem ...string) string {
	elems := append([]string{WorkingDirectory}, elem...)
	return filepath.Join(elems...)
}

func init() {
	if p, err := filepath.Abs("./conf/app.conf"); err == nil {
		ConfigurationFile = p
	}
	if p, err := filepath.Abs("./"); err == nil {
		WorkingDirectory = p
	}
	if p, err := filepath.Abs("./runtime/logs"); err == nil {
		LogFile = p
	}
}
