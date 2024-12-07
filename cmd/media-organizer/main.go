package main

import (
	"flag"
	"log"
	"os"

	"media-organizer/internal/core"
	"media-organizer/internal/utils"
)

func main() {
	srcPath := flag.String("src", "", "源文件夹路径")
	destPath := flag.String("dest", "", "目标文件夹路径")
	flag.Parse()

	if *srcPath == "" || *destPath == "" {
		log.Fatal("请提供源文件夹和目标文件夹路径")
	}

	// 初始化日志
	logger := utils.NewLogger()
	defer logger.Close()

	// 创建目标文件夹
	err := os.MkdirAll(*destPath, 0755)
	if err != nil {
		logger.Fatal("创建目标文件夹失败:", err)
	}

	// 创建组织器实例
	organizer := core.NewOrganizer(*srcPath, *destPath, logger)

	// 开始处理
	err = organizer.Process()
	if err != nil {
		logger.Fatal("处理失败:", err)
	}
}