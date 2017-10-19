package main

import (
	"fmt"
	. "github.com/devfeel/dottask"
	"time"
)

var service *TaskService

func Job_Config(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Job_Config")
	//time.Sleep(time.Second * 3)
	return nil
}

func Loop_Config(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Loop_Config")
	//time.Sleep(time.Second * 3)
	return nil
}

func RegisterTask(service *TaskService) {
	service.RegisterHandler("Job_Config", Job_Config)
	service.RegisterHandler("Loop_Config", Loop_Config)
}

func main() {
	//step 1: init new task service
	service = StartNewService()
	//step 2: register all task handler
	RegisterTask(service)

	//step 3: load config file
	service.LoadConfig("d:/gotmp/task/task.conf")

	fmt.Println(time.Now().String(), " => Begin Task")
	//step 4: start all task
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}
}
