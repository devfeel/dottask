# Devfeel/DotTask
简约大方的go-task组件
<br>支持cron、loop、queue三种模式


## 特性
* 支持配置方式(xml + json + yaml)与代码方式
* 支持cron、loop、queue三种模式
* cron模式支持“秒 分 时 日 月 周”配置
* loop模式支持毫秒级别
* queue模式支持毫秒级别
* 上次任务没有停止的情况下不触发下次任务
* 支持Exception、OnBegin、OnEnd注入点
* 支持单独执行TaskHandler
* 支持代码级重设Task的相关设置


## 安装：

```
go get -u github.com/devfeel/dottask
```

## 快速开始：

```go
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
	time.Sleep(time.Second * 3)
	return nil
}

func main() {

    //step 1: init new task service
	service = StartNewService()

	//step 2: register task handler
	_, err := service.CreateCronTask("testcron", true, "48-5 */2 * * * *", Job_Test, nil)
	if err != nil {
		fmt.Println("service.CreateCronTask error! => ", err.Error())
	}
	_, err = service.CreateLoopTask("testloop", true, 0, 1000, Loop_Test, nil)
	if err != nil {
		fmt.Println("service.CreateLoopTask error! => ", err.Error())
	}

	//step 3: start all task
	service.StartAllTask()

	fmt.Println(service.PrintAllCronTask())

	for true {
	}

}

```

#### 配置方式
```go

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
#### task.xml.conf:
```
<?xml version="1.0" encoding="UTF-8"?>
<config>
<global isrun="true" logpath="d:/"/>
<tasks>
    <task taskid="Loop_Config" type="loop" isrun="true" duetime="10000" interval="10" handlername="Loop_Config" />
    <task taskid="Job_Config" type="cron" isrun="true" express="0 */5 * * * *" handlername="Job_Config" />
</tasks>
</config>

```


## 关于表达式
* 关于CronTask的TimeExpress 简单解释：
* 基本格式：* * * * * *（6列，以空格分隔）
* f1:第1列表示秒0-59，每一秒用*或*/1 表示。
* f2:第2列表示分钟0-59。
* f3:第3列表示小时0-23。
* f4:第4列表示日期1-31。
* f5:第5列表示月份1-12。
* f6:第6列表示星期几0-7，其中0和7均表示为周日。
* 当f1为 * 时表示每秒都要执行任务，f2为 * 时表示每分钟都要执行程序，其余类推
* 当f1为 a-b 时表示从第 a 秒钟到第 b 秒钟这段时间内要执行，f2 为 a-b 时表示从第 a 到第 b 分钟都要执行，其余类推
* 当f1为 */n 时表示每 n 秒钟个时间间隔执行一次，f2 为 */n 表示每 n 小时个时间间隔执行一次，其余类推
* 当f1为 a, b, c,... 时表示第 a, b, c,... 秒钟要执行，f2 为 a, b, c,... 时表示第 a, b, c...分钟要执行，其余类推
#### 示例：
* #每天早上7点执行一次调度任务:
* 0 0 7 * * *
* #在 12 月内, 每天的早上 6 点到 12 点中，每隔3个小时执行一次调度任务:
* 0 0 6-12/3 * 12 *
* #周一到周五每天下午 5:00执行一次调度任务:
* 0 0 17 * * 1-5
* #每月每天的午夜 0 点 20 分, 2 点 20 分, 4 点 20 分....执行一次调度任务
* 0 20 0-23/2 * * *
* #每月每天的0 点 20 分, 9 点 20 分, 16 点 20 分执行一次调度任务
* 0 20 0,9,16 * * *


## 外部依赖
yaml - https://gopkg.in/yaml.v2


## 如何联系
QQ群：193409346
