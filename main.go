package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/dablelv/go-huge-util/zip"
	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/gookit/i18n"
)

func init() {
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

	DefaultLang := config.String("quickclash.lang")
	Languages := map[string]string{
		"en":    "English",
		"zh-CN": "简体中文",
	}
	i18n.Init("lang/", DefaultLang, Languages)

}

func cancelProxy() {
	if err := SetProxy(""); err == nil {
		color.BgLightBlue.Println(i18n.Dtr("cancelProxySuccess"))
	} else {
		color.BgRed.Println(i18n.Dtr("cancelProxyFail", err))
	}
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

	// 如果文件不存在或超过6小时则下载
	if errStat != nil || duration.Hours() > 6 {
		SubUrl = config.String("quickclash.sublink")
		if SubUrl == "" {
			color.BgLightBlue.Println(i18n.Dtr("inputSublink"))
			fmt.Scanln(&SubUrl)
		}
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
		panic(err)
	}

	//检查内核
	if config.String("quickclash.core") == "clash" {
		Core = ClashCore
		color.BgLightBlue.Println(i18n.Dtr("currentCore", "Clash"))
	} else if config.String("quickclash.core") == "clash.meta" {
		Core = ClashMetaCore
		color.BgLightBlue.Println(i18n.Dtr("currentCore", "ClashMeta"))
	} else {
		panic(color.BgRed.Sprintf(i18n.Dtr("coreSetWrong", QuickClashYml)))
	}

	//检查内核是否存在
	_, err = os.Stat(Core)
	// 如果文件不存在则下载
	if err != nil {
		color.FgLightBlue.Println(i18n.Dtr("downloading", Core+" Core"))
		err = Download(ReleaseBaseUrl+Core, Core)
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

	//设置系统代理
	proxy := fmt.Sprintf("127.0.0.1:%s", config.String("port"))
	defer cancelProxy()
	if err := SetProxy(proxy); err == nil {
		color.BgLightBlue.Println(i18n.Dtr("setProxySuccess", proxy))
	} else {
		color.BgRed.Println(i18n.Dtr("setProxyFail", proxy, err))
	}

	//提示信息
	color.BgLightBlue.Println(i18n.Dtr("baseUrl", "127.0.0.1"+config.String("external-controller"), config.String("quickclash.secret")))
	color.BgLightRed.Println(i18n.Dtr("webUIStarted", GinPort))

	//		阻止主协程退出，以保持Gin服务的运行
	//select {}
	c := make(chan os.Signal)
	signal.Notify(c)
	s := <-c
	fmt.Println("get signal: ", s)
}
