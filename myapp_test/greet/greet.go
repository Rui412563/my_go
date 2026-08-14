// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package main

import (
	"bufio"
	"flag"
	"github.com/zeromicro/go-zero/core/logx"
	"greet/logwriter"
	"io"
	"log"
	"os"

	"greet/internal/config"
	"greet/internal/handler"
	"greet/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/greet-api.yaml", "the config file")

func main() {
	flag.Parse()

	// 1. 第一步：先加载配置，获取日志路径
	var c config.Config
	conf.MustLoad(*configFile, &c)

	logPath := c.Log.Path
	if logPath == "" {
		logPath = "logs"
	}

	// 2. 第二步：【极其重要】立即初始化 Writer 并重定向所有输出
	// 必须在 svc.NewServiceContext(c) 之前执行
	allInOneWriter := logwriter.NewSmartRotationWriter(logPath, 5)
	setupGlobalRedirection(allInOneWriter)

	// 3. 第三步：配置 go-zero 的 logx
	logx.MustSetup(c.Log)
	logx.SetWriter(logx.NewWriter(allInOneWriter))

	// 4. 第四步：初始化服务上下文（此时内部的 DB 连接日志就会进入文件了）
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	server.Start()
}

// setupGlobalRedirection 统一处理所有输出源
func setupGlobalRedirection(w io.Writer) {
	// A. 重定向标准库 log 包 (解决你终端看到的 2026/08/11... 格式日志)
	log.SetOutput(w)

	// B. 重定向 os.Stdout 和 os.Stderr (解决 fmt.Print 和 panic)
	r, pw, _ := os.Pipe()
	os.Stdout = pw
	os.Stderr = pw

	go func() {
		// 使用 scanner 逐行读取并写入自定义 Writer，确保触发 5MB 切分逻辑
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			_, _ = w.Write(append(scanner.Bytes(), '\n'))
		}
	}()
}
