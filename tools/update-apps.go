package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	TARGET_JSON_LOCATION = "settingApi/index.json"
	BUNDLE_FOLDER       = "bundle"
)

func main() {
	// 取得當前工作目錄
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	targetJsonPath := filepath.Join(cwd, TARGET_JSON_LOCATION)
	bundleFolderPath := filepath.Join(cwd, BUNDLE_FOLDER)

	// 讀取 JSON 檔案
	jsonFile, err := ioutil.ReadFile(targetJsonPath)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// 解析 JSON
	var targetJson map[string]string
	if err := json.Unmarshal(jsonFile, &targetJson); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// 讀取 bundle 目錄
	folders, err := ioutil.ReadDir(bundleFolderPath)
	if err != nil {
		fmt.Println("Error reading bundle folder:", err)
		return
	}

	for _, folder := range folders {
		if folder.IsDir() {
			folderPath := filepath.Join(bundleFolderPath, folder.Name())
			files, err := ioutil.ReadDir(folderPath)
			if err != nil {
				fmt.Println("Error reading folder:", folderPath, err)
				continue
			}

			parts := strings.Split(folder.Name(), "-")
			gameType := parts[len(parts)-1]

			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".js") {
					fmt.Println("Found JS file:", file.Name())

					nameParts := strings.Split(file.Name(), ".")
					if len(nameParts) == 3 {
						hash := nameParts[1]
						if _, exists := targetJson[gameType]; !exists {
							targetJson[gameType] = hash
						}
					}
				}
			}
		}
	}

	// 寫回 JSON 檔案
	updatedJson, err := json.MarshalIndent(targetJson, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	if err := ioutil.WriteFile(targetJsonPath, updatedJson, 0644); err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}

	fmt.Println("Updated JSON saved successfully!")
}
