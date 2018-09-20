package task

import (
	"time"
	"fmt"
	"errors"
)

type (
	//LoopTask loop task info define
	LoopTask struct {
		TaskInfo
		Interval     int64        `json:"interval"` //运行间隔时间，单位毫秒，当TaskType==TaskType_Loop||TaskType_Queue时有效
	}

)

// GetConfig get task config info
func (task *LoopTask) GetConfig() *TaskConfig{
	return &TaskConfig{
		TaskID:task.taskID,
		TaskType:task.TaskType,
		IsRun : task.IsRun,
		Handler:task.handler,
		DueTime:task.DueTime,
		Interval:task.Interval,
		Express:"",
		TaskData:task.Context().TaskData,
	}
}


//Reset first check conf, then reload conf & restart task
//special, TaskID can not be reset
//special, if TaskData is nil, it can not be reset
//special, if Handler is nil, it can not be reset
func (task *LoopTask) Reset(conf *TaskConfig) error {
	if conf.DueTime < 0 {
		errmsg := "DueTime is wrong format => must bigger or equal then zero"
		task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "error", errmsg)
		return errors.New(errmsg)
	}

	if conf.Interval <= 0 {
		errmsg := "interval is wrong format => must bigger then zero"
		task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "error", errmsg)
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
	task.DueTime = conf.DueTime
	task.Interval = conf.Interval
	task.Start()
	task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "success")
	return nil
}


//Start start task
func (task *LoopTask) Start() {
	if !task.IsRun {
		return
	}

	task.mutex.Lock()
	defer task.mutex.Unlock()

	if task.State == TaskState_Init || task.State == TaskState_Stop {
		task.State = TaskState_Run
		startLoopTask(task)
	}
}


// RunOnce do task only once
// no match Express or Interval
// no recover panic
// support for #6 新增RunOnce方法建议
func (task *LoopTask) RunOnce() error {
	err := task.handler(task.context)
	return err
}


//start loop task
func startLoopTask(task *LoopTask) {
	handler := func() {
		defer func() {
			if err := recover(); err != nil {
				//task.taskService.Logger().Debug(task.TaskID, " loop handler recover error => ", err)
				if task.taskService.ExceptionHandler != nil {
					task.taskService.ExceptionHandler(task.Context(), fmt.Errorf("%v", err))
				}
				//goroutine panic, restart cron task
				startLoopTask(task)
				task.taskService.Logger().Debug(task.TaskID, " goroutine panic, restart LoopTask")
			}
		}()
		//do log
		if task.taskService.OnBeforHandler != nil {
			task.taskService.OnBeforHandler(task.Context())
		}
		var err error
		if !task.Context().IsEnd {
			err = task.handler(task.Context())
		}
		if err != nil {
			if task.taskService.ExceptionHandler != nil {
				task.taskService.ExceptionHandler(task.Context(), err)
			}
		} else {
			//task.taskService.Logger().Debug(task.TaskID, " loop handler end success")
		}
		if task.taskService.OnEndHandler != nil {
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
	//等待设定的延时毫秒
	if task.DueTime > 0 {
		go time.AfterFunc(time.Duration(task.DueTime)*time.Millisecond, dofunc)
	} else {
		go dofunc()
	}

}
