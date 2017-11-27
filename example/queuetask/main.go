package main

import (
	"fmt"
	. "github.com/devfeel/dottask"
	"time"
)

var service *TaskService
const(
	taskName = "TestQueue"
)

func errorHandler(ctx *TaskContext, err error) {
	fmt.Println(time.Now().String(), " => Error ", ctx.TaskID, err.Error())
}

func Job_DealMessage(ctx *TaskContext) error {
	fmt.Println("Job_DealMessage -", ctx.Message)
	return nil
}

func enqueeMessage() {
	for i:=0;i<100;i++ {
		t, exists := service.GetTask(taskName)
		if exists {
			qTask := t.(*QueueTask)
			qTask.EnQueue(time.Now().String())
			fmt.Println(i, t.TaskID(), "enqueeMessage success")
		} else {
			fmt.Println("enqueeMessage not exists ", taskName)
		}
		time.Sleep(time.Second * 1)
	}
}

func main() {
	service = StartNewService()

	qTask, err := service.CreateQueueTask(taskName, true, 1, Job_DealMessage, nil, DefaultQueueSize)
	if err != nil {
		fmt.Println("service.CreateQueueTask error! => ", err.Error())
	}else{
		fmt.Println("service.CreateQueueTask success! => ", qTask.TaskID())
	}

	service.StartAllTask()

	service.SetExceptionHandler(errorHandler)

	t, exists := service.GetTask(taskName)
	if exists {
		err = t.RunOnce()
		if err != nil {
			fmt.Println(t.Context(), "RunOnce error =>", err)
		}else{
			fmt.Println(t.Context(), "RunOnce success")
		}
	}

	fmt.Println(service.PrintAllCronTask())

	t, exists = service.GetTask(taskName)
	if exists {
		conf := &TaskConfig{
			IsRun:   true,
			Interval:1000,
		}
		err = t.Reset(conf)
		if err != nil {
			fmt.Println(t, "Reset error =>", err)
		}else{
			fmt.Println(t, "Reset success ")
		}
	}

	go enqueeMessage()

	for true {
	}

}
