package file_utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ReadJSON 从给定路径读取 JSON 文件并反序列化为指定的结构体
func ReadJSON(filePath string, out interface{}) error {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// 反序列化 JSON 数据
	err = json.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// SaveJSON 将指定结构体序列化为 JSON 并写入到文件
func SaveJSON(filePath string, data interface{}) error {
	// 打开文件，如果文件不存在则创建
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 序列化数据为 JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 设置缩进格式
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
