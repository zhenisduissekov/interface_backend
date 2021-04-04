package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const PATH = "./data_folder"
const BACKUP_PATH = "./backup_folder"

func readFilenamesFromDir(path string) (interface{}, error){
	filesInDir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]string,0)
	for _, f:= range filesInDir {
		result = append(result, f.Name())
	}
	return result, nil
}

func readFromFileJSON(filename string) (interface{}, error) {
	filePath := filepath.Join(PATH, filename)
	jsonFile, err := os.Open(filePath)
	if err!= nil {
		return "", err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err!= nil {
		return "", err
	}

	var result interface{}
	if err = json.Unmarshal([]byte(byteValue), &result); err !=nil {
		return "", err
	}
	return result, nil
}

func saveJSONtoFile(filename string, data interface{}) error {
	filePath := filepath.Join(PATH, filename)
	fmt.Println("data", data)
	if _, err := os.Stat(filePath); err == nil {
		backupFile(filename)
	}
	saveFile(filePath, data)
	return nil

}


func saveFile(filePath string, data interface{}) error {
	jsonData, err:= json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func deleteFile(filename string) error {
	filePath := filepath.Join(PATH, filename)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func backupFile(filename string) error {
	fmt.Println("backup file")
	filePath := filepath.Join(PATH, filename)
	dateNow := time.Now().Format("D02-01-2006T15:04:05")
	backupPath := filepath.Join(BACKUP_PATH, filename + dateNow)
	fmt.Println(backupPath)

	in, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}


func sendResetCommandToPrometheus() error {
req, err := http.NewRequest("POST", "http://localhost:9998/-/reload", nil)
if err != nil {
	return err
}
req.SetBasicAuth("promadmin", "W74!65Mp+")

resp, err := http.DefaultClient.Do(req)
if err != nil {
	return err
}
defer resp.Body.Close()
return nil
}