package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"media-organizer/internal/tui"
	"media-organizer/internal/utils"

	tea "github.com/charmbracelet/bubbletea"
)

type Organizer struct {
	srcPath  string
	destPath string
	logger   *utils.Logger
	progress *tui.ProgressModel
	program  *tea.Program
}

func NewOrganizer(src, dest string, logger *utils.Logger) *Organizer {
	return &Organizer{
		srcPath:  src,
		destPath: dest,
		logger:   logger,
	}
}

func (o *Organizer) Process() error {
	// 首先统计文件总数
	var totalFiles int
	err := filepath.Walk(o.srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && IsMediaFile(path) {
			totalFiles++
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 初始化进度显示
	o.progress = tui.NewProgress(totalFiles)
	o.program = tea.NewProgram(o.progress)

	var wg sync.WaitGroup
	workerCount := 4
	filesChan := make(chan string, 100)
	errChan := make(chan error, 1)
	processedCount := atomic.Int32{}

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go o.worker(&wg, filesChan, &processedCount)
	}

	// 启动进度显示
	go func() {
		if err := o.program.Start(); err != nil {
			errChan <- err
		}
	}()

	// 遍历文件
	err = filepath.Walk(o.srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && IsMediaFile(path) {
			filesChan <- path
		}
		return nil
	})

	close(filesChan)
	wg.Wait()

	// 标记完成
	o.program.Send(o.progress.Done())

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

func (o *Organizer) worker(wg *sync.WaitGroup, files <-chan string, count *atomic.Int32) {
	defer wg.Done()

	for file := range files {
		if err := o.processFile(file); err != nil {
			switch e := err.(type) {
			case *utils.MediaError:
				switch e.Type {
				case utils.ErrMetadata:
					o.logger.Warn(e.Error())
				case utils.ErrFileAccess:
					o.logger.Error(e.Error())
				case utils.ErrCopy:
					o.logger.Error(e.Error())
				case utils.ErrInvalidFormat:
					o.logger.Debug(e.Error())
				}
			default:
				o.logger.Error(fmt.Sprintf("处理文件 %s 时发生未知错误: %v", file, err))
			}
			continue
		}

		current := int(count.Add(1))
		o.program.Send(o.progress.UpdateProgress(current, file))
	}
}

func (o *Organizer) processFile(filePath string) error {
	// 检查是否为支持的媒体文件
	if !IsMediaFile(filePath) {
		o.logger.Debug(fmt.Sprintf("跳过非媒体文件: %s", filePath))
		return nil
	}

	// 获取文件元数据
	metadata, err := GetMetadata(filePath)
	if err != nil {
		return utils.NewMetadataError(filePath, err)
	}

	// 构建目标路径
	destDir := o.buildDestPath(metadata.Time, metadata.IsVideo)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return utils.NewFileAccessError(destDir, err)
	}

	// 复制文件
	destPath := filepath.Join(destDir, filepath.Base(filePath))
	if _, err := os.Stat(destPath); err == nil {
		o.logger.Info(fmt.Sprintf("文件已存在，跳过: %s", destPath))
		return nil
	}

	if err := o.copyFile(filePath, destPath); err != nil {
		return utils.NewCopyError(filePath, err)
	}

	o.logger.Info(fmt.Sprintf("成功处理文件: %s -> %s", filePath, destPath))
	return nil
}

func (o *Organizer) buildDestPath(t time.Time, isVideo bool) string {
	mediaType := "photos"
	if isVideo {
		mediaType = "videos"
	}
	return filepath.Join(o.destPath, mediaType, fmt.Sprintf("%d/%02d", t.Year(), t.Month()))
}

func (o *Organizer) copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dest, input, 0644)
}