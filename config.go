package task

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	AppConfig struct {
		XMLName xml.Name `xml:"config"`

		Global struct {
			LogPath string `xml:"logpath,attr"  yaml:"logpath"`
			IsRun   bool   `xml:"isrun,attr" yaml:"isrun"`
			Timeout int64  `xml:"timeout,attr" yaml:"timeout"` //全局超时配置，单位为毫秒
		} `xml:"global" yaml:"global"`

		Tasks []struct {
			TaskID      string `xml:"taskid,attr" yaml:"taskid"`           //task编号，需唯一
			IsRun       bool   `xml:"isrun,attr" yaml:"isrun"`             //标识是否允许task执行,默认为false,如设为flash,则启动后不执行task
			TaskType    string `xml:"type,attr" yaml:"type"`               //Task类型,目前支持loop、cron、queue
			DueTime     int64  `xml:"duetime,attr" yaml:"duetime"`         //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
			Interval    int64  `xml:"interval,attr" yaml:"interval"`       //loop类型下,两次Task执行之间的间隔,单位为毫秒
			Express     string `xml:"express,attr" yaml:"express"`         //cron类型下,task执行的时间表达式，具体参考readme
			QueueSize   int64  `xml:"queuesize,attr" yaml:"queuesize"`     //queue类型下,queue初始长度
			HandlerName string `xml:"handlername,attr" yaml:"handlername"` //Task对应的HandlerName,需使用RegisterHandler进行统一注册
			HandlerData string `xml:"handlerdata,attr" yaml:"handlerdata"` //Task对应的自定义数据,可在配置源中设置
			Timeout     int64  `xml:"timeout,attr" yaml:"timeout"`         //全局超时配置，单位为毫秒
		} `xml:"tasks>task" yaml:"tasks"`
	}

	ConfigHandle func(configSource string) (*AppConfig, error)
)

//初始化配置文件（xml）
func XmlConfigHandler(configFile string) *AppConfig {
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
func JsonConfigHandler(configFile string) *AppConfig {
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
func YamlConfigHandler(configFile string) *AppConfig {
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
