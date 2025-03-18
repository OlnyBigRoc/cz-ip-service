package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// IsFileExist 检查文件是否存在
func IsFileExist(filePath string) (bool, error) {

	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil // 文件存在
	}
	if os.IsNotExist(err) {
		return false, nil // 文件不存在
	}
	return false, err // 其他错误（如权限问题）
}

// IsDirExist 检查目录是否存在
func IsDirExist(dirPath string) (bool, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 目录不存在
		}
		return false, err // 其他错误（如权限问题）
	}
	return fileInfo.IsDir(), nil
}

// CreateDir 创建目录
func CreateDir(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil { // 0755 表示目录的读写权限
		return fmt.Errorf("create directory failed: %v", err)
	}
	return nil
}

// DownloadFile 实现文件下载逻辑
func DownloadFile(url, filePath string) error {
	// 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file failed: %v", err)
	}
	defer out.Close()

	// 发起 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 流式写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("write file failed: %v", err)
	}

	// 返回文件名称
	return nil
}

// DecompressZip 减压缩
func DecompressZip(zipPath, destDir string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open zip failed: %v", err)
	}
	defer r.Close()

	// 创建目标目录
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %v", err)
	}

	// 遍历 ZIP 文件内容
	for _, f := range r.File {
		// 防止 ZIP 路径遍历攻击
		targetPath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(targetPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		// 创建目录结构
		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, f.Mode())
			continue
		}

		// 创建文件
		if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("create subdir failed: %v", err)
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("open file failed: %v", err)
		}

		// 复制文件内容
		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("open zip entry failed: %v", err)
		}

		if _, err = io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return fmt.Errorf("write file failed: %v", err)
		}

		outFile.Close()
		rc.Close()
	}
	return nil
}

func ListDir(path string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var fileList []os.DirEntry
	for _, file := range files {
		fileList = append(fileList, file)
	}

	return fileList, nil
}

func DeleteDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}
