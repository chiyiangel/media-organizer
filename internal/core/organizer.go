package core

import (
	"fmt"
	"io"
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
	skipDirs map[string]bool
}

func NewOrganizer(src, dest string, logger *utils.Logger, skipDirs []string) *Organizer {
	skipMap := make(map[string]bool)
	for _, dir := range skipDirs {
		skipMap[dir] = true
	}

	return &Organizer{
		srcPath:  src,
		destPath: dest,
		logger:   logger,
		skipDirs: skipMap,
	}
}

func (o *Organizer) shouldSkipDir(name string) bool {
	return o.skipDirs[name]
}

func (o *Organizer) Process() error {
	// 首先统计文件总数
	var totalFiles int
	err := filepath.Walk(o.srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否需要跳过文件夹
		if info.IsDir() && o.shouldSkipDir(info.Name()) {
			o.logger.Debug(fmt.Sprintf("跳过文件夹: %s", path))
			return filepath.SkipDir
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
		if _, err := o.program.Run(); err != nil {
			errChan <- err
		}
	}()

	// 遍历文件
	err = filepath.Walk(o.srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否需要跳过文件夹
		if info.IsDir() && o.shouldSkipDir(info.Name()) {
			o.logger.Debug(fmt.Sprintf("跳过文件夹: %s", path))
			return filepath.SkipDir
		}

		if !info.IsDir() && IsMediaFile(path) {
			filesChan <- path
		}
		return nil
	})
	if err != nil {
		return err
	}

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
	mediaType := "Photos"
	if isVideo {
		mediaType = "Videos"
	}
	return filepath.Join(o.destPath, mediaType, fmt.Sprintf("%d/%02d/%02d-%02d", t.Year(), t.Month(), t.Month(), t.Day()))
}

func (o *Organizer) copyFile(src, dest string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy data from source to destination
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Flush file contents to stable storage
	if err := destFile.Sync(); err != nil {
		return err
	}

	return nil
}
