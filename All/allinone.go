package All

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetWindowsDrives 获取Windows系统所有盘符
func GetWindowsDrives() []string {
	var drives []string
	if runtime.GOOS != "windows" {
		return []string{"/"}
	}

	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + ":\\"
		_, err := os.Stat(drivePath)
		if err == nil {
			drives = append(drives, drivePath)
		}
	}
	return drives
}

// AppendBrowserResultsToFile 将浏览器结果追加到输出文件
func AppendBrowserResultsToFile(outputFile string) error {
	resultsDir := "results"
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {
		return nil
	}

	file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入浏览器数据分隔头
	file.WriteString("\n\n")
	file.WriteString("================================================================================\n")
	file.WriteString("                          Browser Data Results\n")
	file.WriteString("================================================================================\n\n")

	err = filepath.Walk(resultsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// 只处理 csv 文件
		if !strings.HasSuffix(strings.ToLower(path), ".csv") {
			return nil
		}

		csvFile, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer csvFile.Close()

		// 写入文件名作为标题
		file.WriteString(fmt.Sprintf("\n--- %s ---\n", filepath.Base(path)))

		scanner := bufio.NewScanner(csvFile)
		for scanner.Scan() {
			file.WriteString(scanner.Text() + "\n")
		}
		file.WriteString("\n")

		return nil
	})

	return err
}
