package task

import (
	"errors"
	"fmt"
	"time"
)

const (
	DefaultQueueSize = 1000
)

type (
	QueueTask struct {
		TaskInfo
		Interval    int64 //运行间隔时间，单位毫秒，当TaskType==TaskType_Loop||TaskType_Queue时有效
		MessageChan chan interface{}
	}
)

//EnQueue enqueue value into message queue
func (task *QueueTask) EnQueue(value interface{}) {
	task.MessageChan <- value
}

//Start start task
func (task *QueueTask) Start() {
	if !task.IsRun {
		return
	}

	if task.State == TaskState_Init || task.State == TaskState_Stop {
		task.State = TaskState_Run
		startQueueTask(task)
	}
}

// RunOnce do task only once
func (task *QueueTask) RunOnce() error {
	err := task.handler(task.Context())
	return err
}

// GetConfig get task config info
func (task *QueueTask) GetConfig() *TaskConfig {
	return &TaskConfig{
		TaskID:   task.taskID,
		TaskType: task.TaskType,
		IsRun:    task.IsRun,
		Handler:  task.handler,
		DueTime:  task.DueTime,
		Interval: 0,
		Express:  "",
		TaskData: task.Context().TaskData,
	}
}

//Reset first check conf, then reload conf & restart task
func (task *QueueTask) Reset(conf *TaskConfig) error {
	if conf.Interval <= 0 {
		errmsg := "interval is wrong format => must bigger then zero"
		task.taskService.Logger().Debug(fmt.Sprint("TaskInfo:Reset ", task, conf, "error", errmsg))
		return errors.New(errmsg)
	}

	//restart task
	task.Stop()
	task.IsRun = conf.IsRun
	if conf.TaskData != nil {
		task.Context().TaskData = conf.TaskData
	}
	if conf.Handler != nil {
		task.handler = conf.Handler
	}
	task.Interval = conf.Interval
	task.Start()
	task.taskService.Logger().Debug(fmt.Sprint("TaskInfo:Reset ", task, conf, "success"))
	return nil
}

// NewQueueTask create new queue task
func NewQueueTask(taskID string, isRun bool, interval int64, handler TaskHandle, taskData interface{}, queueSize int64) (Task, error) {
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(QueueTask)
	task.initCounters()
	task.taskID = context.TaskID
	task.TaskType = TaskType_Queue
	task.IsRun = isRun
	task.handler = handler
	task.Interval = interval
	task.State = TaskState_Init
	task.context = context
	task.MessageChan = make(chan interface{}, queueSize)
	return task, nil
}

//start queue task
func startQueueTask(task *QueueTask) {
	handler := func() {
		defer func() {
			task.Context().Header = nil
			task.Context().Error = nil
			task.CounterInfo().RunCounter.Inc(1)
			if err := recover(); err != nil {
				task.CounterInfo().ErrorCounter.Inc(1)
				if task.taskService.ExceptionHandler != nil {
					task.taskService.ExceptionHandler(task.Context(), fmt.Errorf("%v", err))
				}
			}
		}()
		task.Context().Header = make(map[string]interface{})
		//get value from message chan
		message := <-task.MessageChan
		task.Context().Message = message

		if task.taskService != nil && task.taskService.OnBeforeHandler != nil {
			task.taskService.OnBeforeHandler(task.Context())
		}

		var err error
		if !task.Context().IsEnd {
			err = task.handler(task.Context())
		}

		if err != nil {
			task.Context().Error = err
			task.CounterInfo().ErrorCounter.Inc(1)
			if task.taskService != nil && task.taskService.ExceptionHandler != nil {
				task.taskService.ExceptionHandler(task.Context(), err)
			}
		}

		if task.taskService != nil && task.taskService.OnEndHandler != nil {
			task.taskService.OnEndHandler(task.Context())
		}
	}
	dofunc := func() {
		task.TimeTicker = time.NewTicker(time.Duration(task.Interval) * time.Millisecond)
		handler()
		for {
			select {
			case <-task.TimeTicker.C:
				handler()
			}
		}
	}
	go dofunc()
}
