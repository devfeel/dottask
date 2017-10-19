package main

import (
	"fmt"
	"time"

	. "github.com/devfeel/dottask"
)

var service *TaskService

func Job_Config(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Job_Config")
	return nil
}

func Loop_Config(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Loop_Config")
	time.Sleep(time.Second * 3)
	return nil
}

func main() {
	//step 1: init new task service
	service = StartNewService()

	service.
		//step 2: register all task handler
		RegisterHandler("Job_Config", Job_Config).
		RegisterHandler("Loop_Config", Loop_Config).

		//step 3: load config file
		LoadConfig("d:/gotmp/task/task.json.conf", ConfigType_Json).

		//step 4: start all task
		StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}
}
