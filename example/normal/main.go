package main

import (
	"fmt"
	. "github.com/devfeel/dottask"
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
	return nil
}

func beginHandler(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => OnBegin")
	return nil
}

func endHandler(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => OnEnd")
	return nil
}

func errorHandler(ctx *TaskContext, err error) {
	fmt.Println(time.Now().String(), " => Error ", ctx.TaskID, err.Error())
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

	service.SetExceptionHandler(errorHandler)
	//service.SetOnBeforHandler(beginHandler)
	//service.SetOnEndHandler(endHandler)

	t, exists := service.GetTask("testloop")
	if exists {
		err = t.RunOnce()
		if err != nil {
			fmt.Println(t.Context(), "RunOnce error =>", err)
		}
	}

	fmt.Println(service.PrintAllCronTask())

	for _, t := range service.GetAllTasks(){
		fmt.Println("GetAllTasks", t.TaskID(), t.GetConfig().TaskType, t.GetConfig().IsRun, t.GetConfig().Interval, t.GetConfig().Express)
	}


	t, exists = service.GetTask("testcron")
	if exists {
		conf := &TaskConfig{
			IsRun:   true,
			Express: "0 */1 * * * *",
		}
		err = t.Reset(conf)
		if err != nil {
			fmt.Println(t, "Reset error =>", err)
		}
	}

	for {
		time.Sleep(time.Hour)
	}

}
