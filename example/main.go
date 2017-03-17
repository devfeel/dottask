package main

import (
	"fmt"
	. "github.com/devfeel/task"
	"time"
)

var service *TaskService

const defaultTimeLayout = "2006-01-02 15:04:05"

type FmtLogger struct {
}

func (logger *FmtLogger) Debug(v ...interface{}) {
	fmt.Println(time.Now().Format(defaultTimeLayout), " [DEBUG] ", v[0:])
}
func (logger *FmtLogger) Info(v ...interface{}) {
	fmt.Println(time.Now().Format(defaultTimeLayout), " [INFO] ", v[0:])
}
func (logger *FmtLogger) Warn(v ...interface{}) {
	fmt.Println(time.Now().Format(defaultTimeLayout), " [WARN] ", v[0:])
}
func (logger *FmtLogger) Error(v ...interface{}) {
	fmt.Println(time.Now().Format(defaultTimeLayout), " [ERROR] ", v[0:])
}

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
	service.SetLogger(new(FmtLogger))
	_, err := service.CreateCronTask("testcron", "48-5 */2 * * * *", Job_Test, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask error! => ", err.Error())
	}
	_, err = service.CreateLoopTask("testloop", 1000, Loop_Test, nil)
	if err != nil {
		fmt.Println("service.CreateLoopTask error! => ", err.Error())
	}
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}

}
