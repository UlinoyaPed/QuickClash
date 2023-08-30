package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var contentArray = make([]string, 0, 5)
var IsExecNext = make(chan bool, 1)

func ExecCommand(commandName string, params []string) bool {
	contentArray = contentArray[0:0]
	cmd := exec.Command(commandName, params...)
	// 显示运行的命令
	fmt.Printf("执行命令: %s\n", strings.Join(cmd.Args, " "))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error=>", err.Error())
		return false
	}
	// Start开始执行包含的命令,但并不会等待该命令完成即返回
	// wait方法会返回命令的返回状态码并在命令返回后释放相关的资源
	cmd.Start()
	reader := bufio.NewReader(stdout)
	var index, CU int
	var timebefore, timeafter int64
	ISMached := true
	timebefore = time.Now().UnixNano()
	// 实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		index++
		timeafter = time.Now().UnixNano()
		if timeafter-timebefore >= 1e8 && ISMached {
			CU = index + 1
			ISMached = false
		}
		if CU == index {
			IsExecNext <- true

		}
		fmt.Printf(line)
		contentArray = append(contentArray, line)
		timebefore = timeafter
	}
	cmd.Wait()
	return true
}

// func execNext(ch chan bool) {
// 	for {
// 		if iSNext {
// 			ch <- 1
// 		}
// 	}
// }
