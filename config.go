package task

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v2"
)

type (
	AppConfig struct {
		XMLName xml.Name `xml:"config"`

		Global struct {
			LogPath string `xml:"logpath,attr"  yaml:"logpath"`
			IsRun   bool   `xml:"isrun,attr" yaml:"isrun"`
		} `xml:"global" yaml:"global"`

		Tasks []struct {
			TaskID      string `xml:"taskid,attr" yaml:"taskid"`
			IsRun       bool   `xml:"isrun,attr" yaml:"isrun"`
			TaskType    string `xml:"type,attr" yaml:"type"`
			DueTime     int64  `xml:"duetime,attr" yaml:"duetime"` //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
			Interval    int64  `xml:"interval,attr" yaml:"interval"`
			Express     string `xml:"express,attr" yaml:"express"`
			HandlerName string `xml:"handlername,attr" yaml:"handlername"`
			HandlerData string `xml:"handlerdata,attr" yaml:"handlerdata"`
		} `xml:"tasks>task" yaml:"tasks"`
	}
)

//初始化配置文件（xml）
func InitConfig(configFile string) *AppConfig {
	realFile, exists := lookupFile(configFile)
	if !exists {
		panic("Task:Config:InitConfig 配置文件[" + configFile + "]无法解析 - 无法寻找到指定配置文件")
		os.Exit(1)
	}

	content, err := ioutil.ReadFile(realFile)
	if err != nil {
		panic("Task:Config:InitConfig 配置文件[" + realFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config AppConfig
	err = xml.Unmarshal(content, &config)
	if err != nil {
		panic("Task:Config:InitConfig 配置文件[" + realFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	return &config
}

//初始化配置文件（json）
func InitJsonConfig(configFile string) *AppConfig {
	realFile, exists := lookupFile(configFile)
	if !exists {
		panic("Task:Config:InitConfig 配置文件[" + configFile + "]无法解析 - 无法寻找到指定配置文件")
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(realFile)
	if err != nil {
		panic("Task:Config:InitJsonConfig 配置文件[" + realFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config AppConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		panic("Task:Config:InitJsonConfig 配置文件[" + realFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	return &config
}

//初始化配置文件（yaml）
func InitYamlConfig(configFile string) *AppConfig {
	realFile, exists := lookupFile(configFile)
	if !exists {
		panic("Task:Config:InitConfig 配置文件[" + configFile + "]无法解析 - 无法寻找到指定配置文件")
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(realFile)
	if err != nil {
		panic("Task:Config:InitYamlConfig 配置文件[" + realFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config AppConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		panic("Task:Config:InitYamlConfig 配置文件[" + realFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	return &config
}

func lookupFile(configFile string) (realFile string, exists bool) {
	//add default file lookup
	//1、按绝对路径检查
	//2、尝试在当前进程根目录下寻找
	//3、尝试在当前进程根目录/config/ 下寻找
	//fixed for (#3 当使用json配置的时候，运行会抛出panic)
	realFile = configFile
	exists = true
	if !fileExists(realFile) {
		realFile = getCurrentDirectory() + "/" + configFile
		exists = false
	}
	if !exists && !fileExists(realFile) {
		realFile = getCurrentDirectory() + "/config/" + configFile
	} else {
		exists = true
	}
	if !exists && fileExists(realFile) {
		exists = true
	}
	return realFile, exists
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func getCurrentDirectory() string {
	return filepath.Clean(filepath.Dir(os.Args[0])) + "/"
}
