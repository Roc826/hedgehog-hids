package config

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"regexp"
)

type Config struct {
	BlackExt    []string `json:"BlackExt"`
	WhiteExt    []string `json:"WhiteExt"`
	//MonitorUser []string `json:"MonitorUser"`
	//CheckFd		[]int	 `json:"CheckFd"`
	//DangerBinary	[]string `json:"DangerBinary"`
	MonitorDirs []string `json:"MonitorDirs"`
}

var ConfigPath string="./config.json"
//func GetConfig() *Config {
//
//}
func LoadConfig(path string) Config {
	var config Config
	config_file, err := os.Open(path)
	if err != nil {
		ErrorLog("Failed to open config file '%s': %s\n", path, err)
		return config
	}

	fi, _ := config_file.Stat()
	if fi.Size() == 0 {
		ErrorLog("config file(%q) is empty,skipping\n", path)
		return config
	}
	buffer := make([]byte, fi.Size())
	_, err = config_file.Read(buffer)
	if err != nil {
		log.Println()
	}
	buffer, err = StripComments(buffer) //去掉json里的注释
	if err != nil {
		ErrorLog("Failed to strip comments from json: %s\n", err)
		return config
	}

	buffer = []byte(os.ExpandEnv(string(buffer))) //特殊,处理系统的特殊变量
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		ErrorLog("Failed unmarshalling json: %s\n", err)
		return config
	}
	return config
}

func StripComments(data []byte) ([]byte, error) {
	data = bytes.Replace(data, []byte("\r"), []byte(""), 0) // Windows
	lines := bytes.Split(data, []byte("\n"))                //split to muli lines
	filtered := make([][]byte, 0)

	for _, line := range lines {
		match, err := regexp.Match(`^\s*#`, line)
		if err != nil {
			return nil, err
		}
		if !match {
			filtered = append(filtered, line)
		}
	}

	return bytes.Join(filtered, []byte("\n")), nil
}

func ErrorLog(str string, args ...interface{}) {
	log.Printf(str, args...)
}
