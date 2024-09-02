package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func JsonDecode(jsonData []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil, err
	}

	// 递归处理可能嵌套的 JSON
	for key, value := range result {
		switch value.(type) {
		case map[string]interface{}:
			// 如果是嵌套的 JSON 对象，递归解析
			nestedJson, err := json.Marshal(value)
			if err != nil {
				fmt.Println("Error marshalling nested JSON:", err)
				continue
			}
			nestedMap, err := JsonDecode(nestedJson)
			if err != nil {
				fmt.Println("Error decoding nested JSON:", err)
				continue
			}
			result[key] = nestedMap
		case []interface{}:
			// 如果是数组，处理数组中的每个元素
			for i, item := range value.([]interface{}) {
				switch item.(type) {
				case map[string]interface{}:
					nestedJson, err := json.Marshal(item)
					if err != nil {
						fmt.Println("Error marshalling item JSON:", err)
						continue
					}
					nestedMap, err := JsonDecode(nestedJson)
					if err != nil {
						fmt.Println("Error decoding item JSON:", err)
						continue
					}
					value.([]interface{})[i] = nestedMap
				}
			}
		}
	}

	return result, nil
}

func readUntilEndMarker(conn net.Conn) (string, error) { // read throw tcp connection
	reader := bufio.NewReader(conn)
	var message string
	for {
		part, err := reader.ReadString('#') // 读取到第一个 #
		if err != nil {
			return "", err
		}
		message += part

		// 检查是否以 #END# 结束
		if len(message) >= 5 && message[len(message)-5:] == "#END#" {
			// 去掉结束符
			message = message[:len(message)-5]
			break
		}
	}

	return message, nil
}
