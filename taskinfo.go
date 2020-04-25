package task

import (
	"context"
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
		taskID      string `json:"taskid"`
		IsRun       bool   `json:"isrun"`
		taskService *TaskService
		mutex       sync.RWMutex
		TimeTicker  *time.Ticker `json:"-"`
		TaskType    string       `json:"tasktype"`
		handler     TaskHandle
		TaskData    interface{}
		State       string `json:"state"`   //匹配 TskState_Init、TaskState_Run、TaskState_Stop
		DueTime     int64  `json:"duetime"` //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
		Timeout     int64

		counters *CounterInfo
	}

	//TaskConfig task config
	TaskConfig struct {
		TaskID   string
		TaskType string
		IsRun    bool
		Handler  TaskHandle `json:"-"`
		DueTime  int64
		Interval int64
		Express  string
		TaskData interface{}
	}

	CounterInfo struct {
		StartTime    time.Time
		RunCounter   Counter
		ErrorCounter Counter
	}

	ShowCountInfo struct {
		TaskID string
		Lable  string
		Count  int64
	}
)

//Stop stop task
func (task *TaskInfo) Stop() {
	if !task.IsRun {
		return
	}
	if task.State == TaskState_Run {
		task.TimeTicker.Stop()
		task.State = TaskState_Stop
	}
}

func (task *TaskInfo) SetTimeout(timeout int64) {
	task.Timeout = timeout
}

//TaskID return taskID
func (task *TaskInfo) TaskID() string {
	return task.taskID
}

//SetTaskService Set up the associated service
func (task *TaskInfo) SetTaskService(service *TaskService) {
	task.taskService = service
}

// RunOnce do task only once
// no match Express or Interval
// no recover panic
// support for #6 新增RunOnce方法建议
func (task *TaskInfo) RunOnce() error {
	err := task.handler(task.getTaskContext())
	return err
}

func (task *TaskInfo) getTaskContext() *TaskContext {
	ctx := task.taskService.contextPool.Get().(*TaskContext)
	ctx.TaskID = task.taskID
	ctx.TaskData = task.TaskData
	ctx.Header = make(map[string]interface{})
	if task.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(task.Timeout))
		ctx.TimeoutContext = timeoutCtx
		ctx.TimeoutCancel = cancel
		ctx.doneChan = make(chan struct{})
	}
	return ctx
}

func (task *TaskInfo) putTaskContext(ctx *TaskContext) {
	ctx.reset()
	task.taskService.contextPool.Put(ctx)
}

func (task *TaskInfo) initCounters() {
	counterInfo := new(CounterInfo)
	counterInfo.StartTime = time.Now()
	counterInfo.RunCounter = NewCounter()
	counterInfo.ErrorCounter = NewCounter()
	task.counters = counterInfo
}

func (task *TaskInfo) CounterInfo() *CounterInfo {
	return task.counters
}
