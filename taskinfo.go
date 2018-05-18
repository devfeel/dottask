package task

import (
	"sync"
	"time"
)

const (
	DefaultPeriod     = time.Second //默认执行周期
	defaultTimeLayout = "2006-01-02 15:04:05"
)

type (
	//TaskInfo task info define
	TaskInfo struct {
		taskID       string `json:"taskid"`
		IsRun        bool   `json:"isrun"`
		taskService  *TaskService
		mutex        sync.RWMutex
		TimeTicker   *time.Ticker `json:"-"`
		TaskType     string       `json:"tasktype"`
		handler      TaskHandle
		context      *TaskContext `json:"context"`
		State        string       `json:"state"`    //匹配 TskState_Init、TaskState_Run、TaskState_Stop
		DueTime      int64        `json:"duetime"`  //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
	}

	//TaskConfig task config
	TaskConfig struct {
		TaskID   string
		TaskType     string
		IsRun    bool
		Handler  TaskHandle
		DueTime  int64
		Interval int64
		Express  string
		TaskData interface{}
	}

	//Task上下文信息
	TaskContext struct {
		TaskID   string
		TaskData interface{} //用于当前Task全局设置的数据项
		Message interface{} //用于每次Task执行上下文消息传输
		IsEnd    bool //如果设置该属性为true，则停止当次任务的后续执行，一般用在OnBegin中
	}

	TaskHandle func(*TaskContext) error
)


//Stop stop task
func (task *TaskInfo) Stop() {
	if !task.IsRun {
		return
	}
	if task.State == TaskState_Run {
		task.TimeTicker.Stop()
		task.State = TaskState_Stop
		task.taskService.Logger().Debug(task.TaskID, " Stop")
	}
}

//TaskID return taskID
func (task *TaskInfo) TaskID() string{
	return task.taskID
}

//Context return context
func (task *TaskInfo) Context() *TaskContext{
	return task.context
}

//SetTaskService Set up the associated service
func(task *TaskInfo) SetTaskService(service *TaskService){
	task.taskService = service
}

// RunOnce do task only once
// no match Express or Interval
// no recover panic
// support for #6 新增RunOnce方法建议
func (task *TaskInfo) RunOnce() error {
	err := task.handler(task.context)
	return err
}

