package main

import (
	"fmt"
	. "github.com/devfeel/dottask"
	"time"
)

var service *TaskService

var firstLoopTimeout = 0
var firstCronTimeout = 0

var patchLoop = 0
var patchCron = 0

func Job_Timeout_Test(ctx *TaskContext) error {
	patchCron += 1
	patch := patchLoop

	if firstLoopTimeout <= 0 {
		firstLoopTimeout = 1
		time.Sleep(time.Second * 15)
	} else {
		time.Sleep(time.Second)
	}

	fmt.Println(time.Now().String(), " => Job_Timeout_Test", patch)
	return nil
}

func Loop_Timeout_Test(ctx *TaskContext) error {
	patchLoop += 1
	patch := patchLoop
	if firstCronTimeout <= 0 {
		firstCronTimeout = 1
		time.Sleep(time.Second * 20)
	}

	fmt.Println(time.Now().String(), " => Loop_Timeout_Test", patch)
	return nil
}

func main() {
	service = StartNewService()

	t1, err := service.CreateCronTask("test-timeout-cron", true, "*/10 * * * * *", Job_Timeout_Test, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask error! => ", err.Error())
	}
	t1.SetTimeout(5)
	t2, err := service.CreateLoopTask("test-timeout-loop", true, 0, 10000, Loop_Timeout_Test, nil)
	if err != nil {
		fmt.Println("service.CreateLoopTask error! => ", err.Error())
	}
	t2.SetTimeout(5)

	service.StartAllTask()

	for {
		time.Sleep(time.Hour)
	}

}
