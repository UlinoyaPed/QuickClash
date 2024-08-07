package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/dablelv/go-huge-util/zip"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/i18n"
)

func init() {
	cancelProxy()
	// 设置选项支持 ENV 解析
	config.WithOptions(config.ParseEnv)
	// 添加驱动程序以支持yaml内容解析（除了JSON是默认支持，其他的则是按需使用）
	config.AddDriver(yaml.Driver)

	//检查配置文件是否存在
	_, err := os.Stat(QuickClashYml)
	// 如果文件不存在则下载
	if err != nil {
		os.Mkdir("config", os.ModePerm)
		color.FgLightBlue.Println("Downloading config file!")
		Download(RepoBaseUrl+QuickClashYml, QuickClashYml)
	}

	// 加载配置，可以同时传入多个文件
	err = config.LoadFiles(QuickClashYml)
	if err != nil {
		panic(err)
	}

	//多语言
	DefaultLang := config.String("quickclash.lang")
	Languages := map[string]string{
		"ara":   "العربية",
		"de":    "Deutsch",
		"en":    "English",
		"fra":   "français",
		"jp":    "日本語",
		"kor":   "한국어",
		"ru":    "русский",
		"th":    "ภาษาไทย",
		"zh-CN": "简体中文",
		"zh-TW": "繁體中文",
	}
	i18n.Init("lang/", DefaultLang, Languages)

}

func cancelProxy() {
	//取消代理
	if err := SetProxy(""); err == nil {
		color.BgLightBlue.Println(i18n.Dtr("cancelProxySuccess"))
	} else {
		color.BgRed.Println(i18n.Dtr("cancelProxyFail", err))
	}
	//等待2秒，防止按下关闭后看不清取消代理的状态
	time.Sleep(2 * time.Second)
}

func main() {
	var err error

	_, err = os.Stat("web/")
	if err != nil {
		color.BgLightBlue.Println(i18n.Dtr("downloading", "WebUI"))
		Download(WebUIUrl, "web.zip")
		color.BgLightBlue.Println(i18n.Dtr("unziping"))
		zip.Unzip("web.zip", ".")
	}
	//设置Gin为发布模式
	gin.SetMode(gin.ReleaseMode)
	// 创建一个Gin实例
	r := gin.Default()

	// 定义路由
	//静态文件路由，作为主页
	r.Static("/", "./web")

	//设置端口
	GinPort := fmt.Sprintf(":%s", config.String("quickclash.port"))
	// 使用协程启动Gin服务
	go func() {
		if err = r.Run(GinPort); err != nil {
			panic(err)
		}
	}()

	//		主协程可以继续执行其他操作

	//检查配置文件是否存在
	_, errStat := os.Stat(QuickClashSubYml)
	var duration time.Duration
	if errStat == nil {
		//检查文件修改时间
		fileInfo, _ := os.Stat(QuickClashSubYml)
		modTime := fileInfo.ModTime()
		duration = time.Since(modTime)
	}

	// 如果文件不存在或超过多少小时则下载
	if errStat != nil || duration.Hours() > config.Float("quickclash.duration") {
		SubUrl = config.String("quickclash.sublink")
		if SubUrl == "" {
			color.BgLightBlue.Println(i18n.Dtr("inputSublink"))
			fmt.Scanln(&SubUrl)
		}
		SubUrl = AddHTTPSPrefix(SubUrl)
		fmt.Println(SubUrl)
		os.Mkdir("yaml", os.ModePerm)
		color.BgLightBlue.Println(i18n.Dtr("downloading", i18n.Dtr("configFile")))
		err = Download(SubUrl, QuickClashSubYml)
		if err != nil {
			panic(err)
		}
		color.BgLightBlue.Println(i18n.Dtr("dlComplete"))
	} else {
		color.BgLightBlue.Println(i18n.Dtr("noneedUpdateC"))
	}

	//加载配置文件
	err = config.LoadFiles(QuickClashSubYml)
	if err != nil {
		os.Remove(QuickClashSubYml)
		panic(err)
	}

	//检查内核
	if runtime.GOARCH == "386" || config.String("quickclash.core") == "clash.metax86" {
		Core = ClashMetaCorex86
		color.BgLightBlue.Println(i18n.Dtr("currentCore", "ClashMeta X86"))
	} else if runtime.GOARCH == "amd64" {
		if config.String("quickclash.core") == "clash" {
			Core = ClashCore
			color.BgLightBlue.Println(i18n.Dtr("currentCore", "Clash"))
		} else if config.String("quickclash.core") == "clash.meta" {
			Core = ClashMetaCore
			color.BgLightBlue.Println(i18n.Dtr("currentCore", "ClashMeta"))
		}
	} else {
		panic(color.BgRed.Sprintf(i18n.Dtr("coreSetWrong", QuickClashYml)))
	}

	//检查内核是否存在
	_, err = os.Stat(Core)
	// 如果文件不存在则下载
	if err != nil {
		color.FgLightBlue.Println(i18n.Dtr("downloading", Core+" Core"))
		err = DownloadWithBar(ReleaseBaseUrl+Core, Core)
		if err != nil {
			panic(err)
		}
	}

	//启动Clash
	color.BgLightBlue.Println(i18n.Dtr("clashOpening"))
	command := "./" + Core
	params := []string{"-f", QuickClashSubYml, "-secret", config.String("quickclash.secret")}
	go ExecCommand(command, params)

	//等待2秒
	time.Sleep(2 * time.Second)

	//关闭时自动取消代理
	defer cancelProxy()
	//设置系统代理
	proxyport := config.String("port")
	if proxyport == "" {
		proxyport = config.String("mixed-port")
	}
	proxy := fmt.Sprintf("127.0.0.1:%s", proxyport)
	if err := SetProxy(proxy); err == nil {
		color.BgLightBlue.Println(i18n.Dtr("setProxySuccess", proxy))
	} else {
		color.BgRed.Println(i18n.Dtr("setProxyFail", proxy, err))
	}

	//提示信息
	color.BgLightBlue.Println(i18n.Dtr("baseUrl", "127.0.0.1"+config.String("external-controller"), config.String("quickclash.secret")))
	color.BgLightRed.Println(i18n.Dtr("webUIStarted", GinPort))

	//		阻止主协程退出，以保持Gin服务的运行
	Clog()
}

func Clog() {
	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("get signal: ", s)
}
