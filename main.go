package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
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
		color.FgLightBlue.Println("配置文件不存在，正在下载...")
		Download(RepoBaseUrl+QuickClashYml, QuickClashYml)
	}

	// 加载配置，可以同时传入多个文件
	err = config.LoadFiles(QuickClashYml)
	if err != nil {
		panic(err)
	}
}

func cancelProxy() {
	if err := SetProxy(""); err == nil {
		color.BgLightBlue.Println("取消代理设置成功!")
	} else {
		color.BgRed.Printf("取消代理设置失败: %s\n", err)
	}
}

func main() {
	var err error
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
			color.BgLightBlue.Println("您没有填写 quickclash.sublink 字段，请在此处输入")
			fmt.Scanln(&SubUrl)
		}
		color.BgLightBlue.Println("正在下载配置文件，请稍等")
		err = Download(SubUrl, QuickClashSubYml)
		if err != nil {
			panic(err)
		}
		color.BgLightBlue.Println("文件下载完成！")
	} else {
		color.BgLightBlue.Println("配置文件无需更新")
	}

	//加载配置文件
	err = config.LoadFiles(QuickClashSubYml)
	if err != nil {
		panic(err)
	}

	//检查内核
	if config.String("quickclash.core") == "clash" {
		Core = ClashCore
		color.BgLightBlue.Println("当前为 Clash 内核")
	} else if config.String("quickclash.core") == "clash.meta" {
		Core = ClashMetaCore
		color.BgLightBlue.Println("当前为 ClashMeta 内核")
	} else {
		panic(color.BgRed.Sprintf("内核设置有误，请检查 %s", QuickClashYml))
	}

	//检查内核是否存在
	_, err = os.Stat(Core)
	// 如果文件不存在则下载
	if err != nil {
		color.FgLightBlue.Println("内核不存在，正在下载...")
		err = Download(ReleaseBaseUrl+Core, Core)
		if err != nil {
			panic(err)
		}
	}

	//启动Clash
	color.BgLightBlue.Println("正在开启Clash")
	command := "./" + Core
	params := []string{"-f", QuickClashSubYml, "-secret", config.String("quickclash.secret")}
	go ExecCommand(command, params)

	//等待2秒
	time.Sleep(2 * time.Second)

	//设置系统代理
	proxy := fmt.Sprintf("127.0.0.1:%s", config.String("port"))
	defer cancelProxy()
	if err := SetProxy(proxy); err == nil {
		color.BgLightBlue.Printf("设置代理服务器: %s 成功!\n", proxy)
	} else {
		color.BgRed.Printf("设置代理服务器: %s 失败, : %s\n", proxy, err)
	}

	//提示信息
	color.BgLightBlue.Printf("管理面板 Base url：http://%s Secret：%s\n", "127.0.0.1"+config.String("external-controller"), config.String("quickclash.secret"))
	color.BgLightRed.Printf("管理面板已启动，请访问 localhost%s\n", GinPort)

	//		阻止主协程退出，以保持Gin服务的运行
	select {}
}
