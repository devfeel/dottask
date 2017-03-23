package task

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type (
	AppConfig struct {
		XMLName xml.Name     `xml:"config"`
		Global  GlobalConfig `xml:"global"`
		Tasks   []TaskConfig `xml:"tasks>task"`
	}

	GlobalConfig struct {
		LogPath string `xml:"logpath,attr"`
		IsRun   bool   `xml:"isrun,attr"`
	}

	//default task config
	TaskConfig struct {
		TaskID      string `xml:"taskid,attr"`
		IsRun       bool   `xml:"isrun,attr"`
		TaskType    string `xml:"type,attr"`
		DueTime     int64  `xml:"duetime,attr"` //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
		Interval    int64  `xml:"interval,attr"`
		Express     string `xml:"express,attr"`
		HandlerName string `xml:"handlername,attr"`
		HandlerData string `xml:"handlerdata,attr"`
	}
)

//初始化配置文件
func InitConfig(configFile string) *AppConfig {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("Task:Config:InitConfig 配置文件[" + configFile + "]无法解析 - " + err.Error())
		os.Exit(1)
	}

	var config AppConfig
	err = xml.Unmarshal(content, &config)
	if err != nil {
		panic("Task:Config:InitConfig 配置文件[" + configFile + "]解析失败 - " + err.Error())
		os.Exit(1)
	}
	return &config
}
