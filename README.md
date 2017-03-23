# Devfeel/Task
简约大方的go-task组件
<br>支持cron、loop两种模式


## 特性
* 支持配置方式与代码方式
* 支持cron、loop两种模式
* cron模式支持“秒 分 时 日 月 周”配置
* loop模式支持毫秒级别
* 上次任务没有停止的情况下下次任务顺延


## 安装：

```
go get -u github.com/devfeel/task
```

## 快速开始：

#### 配置方式
```go

package main

import (
	"fmt"
	. "github.com/devfeel/task"
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
	time.Sleep(time.Second * 3)
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
	service.LoadConfig("d:\\task.conf")

	//step 4: start all task
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}
}

```
task.conf:
```
<?xml version="1.0" encoding="UTF-8"?>
<config>
<global isrun="true" logpath="d:/"/>
<tasks>
    <task taskid="Loop_Config" type="loop" isrun="true" interval="10" handlername="Loop_Config"></task>
    <task taskid="Job_Config" type="cron" isrun="true" express="0 */5 * * * *" handlername="Job_Config"></task>
</tasks>
</config>

```

#### 代码方式

```go
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
	_, err = service.CreateLoopTask("testloop", true, 1000, Loop_Test, nil)
	if err != nil {
		fmt.Println("service.CreateLoopTask error! => ", err.Error())
	}
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}

}

```


## 如何联系
QQ群：193409346
