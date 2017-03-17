# Devfeel/Task
简约大方的go-task组件
<br>支持cron、loop两种模式

## 安装：

```
go get -u github.com/devfeel/task
```

## 快速开始：

```go
func Job_Test(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Job_Test")
	return nil
}

func Loop_Test(ctx *TaskContext) error {
	fmt.Println(time.Now().String(), " => Loop_Test")
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

	for true {
	}
}

```
## 特性
* 支持cron、loop两种模式
* cron模式支持“秒 分 时 日 月 周”配置
* loop模式支持毫秒级别
* 上次任务没有停止的情况下下次任务顺延


## 如何联系
QQ群：193409346
