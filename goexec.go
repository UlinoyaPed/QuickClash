package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

func GoClashExec() {
	cmd := exec.Command("clash.meta.exe", "-f yaml/QuickClashSub.yaml") // 可以改为你需要执行的命令和参数

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	// 创建一个sync.WaitGroup以等待命令执行完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 在Goroutine中读取命令的标准输出
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		wg.Done()
	}()

	// 在另一个Goroutine中读取命令的标准错误
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		wg.Done()
	}()

	// 等待命令执行完成
	cmd.Wait()
	wg.Wait()
}
