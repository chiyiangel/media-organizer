package core

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Metadata struct {
	Time    time.Time
	IsVideo bool
}

var (
	// 支持的图片格式
	imageExtensions = map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".raw":  true,
		".cr2":  true,
		".nef":  true,
		".arw":  true,
	}

	// 支持的视频格式
	videoExtensions = map[string]bool{
		".mp4":  true,
		".mov":  true,
		".avi":  true,
		".mkv":  true,
		".wmv":  true,
		".flv":  true,
		".m4v":  true,
		".3gp":  true,
		".webm": true,
	}
)

func GetMetadata(filePath string) (*Metadata, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 判断是否为视频文件
	if videoExtensions[ext] {
		return getVideoMetadata(filePath)
	}

	// 判断是否为图片文件
	if imageExtensions[ext] {
		return getImageMetadata(filePath)
	}

	// 如果既不是视频也不是图片，使用文件的修改时间
	return getFallbackMetadata(filePath)
}

func getImageMetadata(filePath string) (*Metadata, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 尝试读取EXIF信息
	x, err := exif.Decode(f)
	if err == nil {
		// 尝试获取拍摄时间
		dt, err := x.DateTime()
		if err == nil {
			return &Metadata{
				Time:    dt,
				IsVideo: false,
			}, nil
		}
	}

	// 如果无法读取EXIF信息，回退到使用文件信息
	return getFallbackMetadata(filePath)
}

func getVideoMetadata(filePath string) (*Metadata, error) {
	// 对于视频文件，我们使用文件的修改时间
	return getFallbackMetadata(filePath)
}

func getFallbackMetadata(filePath string) (*Metadata, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	// 使用文件的修改时间
	modTime := fileInfo.ModTime()

	// 判断是否为视频
	ext := strings.ToLower(filepath.Ext(filePath))
	isVideo := videoExtensions[ext]

	return &Metadata{
		Time:    modTime,
		IsVideo: isVideo,
	}, nil
}

// IsMediaFile 检查文件是否为支持的媒体文件
func IsMediaFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return imageExtensions[ext] || videoExtensions[ext]
}
