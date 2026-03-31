package media

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 支持的图片和音频格式
var (
	imageExts = map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".webp": "image/webp",
	}
	audioExts = map[string]string{
		".mp3": "audio/mpeg",
		".wav": "audio/wav",
		".ogg": "audio/ogg",
	}
)

// IsImageFile 判断文件路径是否为支持的图片格式
func IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := imageExts[ext]
	return ok
}

// IsAudioFile 判断文件路径是否为支持的音频格式
func IsAudioFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := audioExts[ext]
	return ok
}

// EncodeFileBase64 读取文件并返回 base64 编码内容
func EncodeFileBase64(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败 \"%s\": %w", path, err)
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// BuildDataURI 根据文件扩展名和 base64 内容构造 data URI（如 data:audio/mpeg;base64,...）
func BuildDataURI(path, base64Data string) string {
	ext := strings.ToLower(filepath.Ext(path))
	// 优先匹配音频，再匹配图片
	if mime, ok := audioExts[ext]; ok {
		return fmt.Sprintf("data:%s;base64,%s", mime, base64Data)
	}
	if mime, ok := imageExts[ext]; ok {
		return fmt.Sprintf("data:%s;base64,%s", mime, base64Data)
	}
	return fmt.Sprintf("data:application/octet-stream;base64,%s", base64Data)
}

// GetMimeType 获取文件的 MIME 类型
func GetMimeType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if mime, ok := audioExts[ext]; ok {
		return mime
	}
	if mime, ok := imageExts[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

// FileExists 检查文件是否存在且非目录
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
