package main

import (
	"fmt"
	. "github.com/devfeel/task"
	"time"
)

var service *TaskService

func Job_Test(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Job_Test")
	//time.Sleep(time.Second * 3)
	return nil
}

func Loop_Test(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Loop_Test")
	time.Sleep(time.Second * 3)
	return nil
}

func main() {
	service = StartNewService()
	_, err := service.CreateCronTask("testcron", true, "48-5 */2 * * * *", Job_Test, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask error! => ", err.Error())
	}
	_, err = service.CreateLoopTask("testloop", true, 0, 1000, Loop_Test, nil)
	if err != nil {
		fmt.Println("service.CreateLoopTask error! => ", err.Error())
	}
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}

}
