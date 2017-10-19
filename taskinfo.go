package task

import (
	"fmt"
	"sync"
	"time"
)

const (
	DefaultPeriod     = time.Second //默认执行周期
	defaultTimeLayout = "2006-01-02 15:04:05"
)

//task info define
type TaskInfo struct {
	TaskID       string `json:"taskid"`
	IsRun        bool   `json:"isrun"`
	taskService  *TaskService
	mutex        sync.RWMutex
	TimeTicker   *time.Ticker `json:"-"`
	TaskType     string       `json:"tasktype"`
	handler      TaskHandle
	Context      *TaskContext `json:"context"`
	State        string       `json:"state"`    //匹配 TskState_Init、TaskState_Run、TaskState_Stop
	DueTime      int64        `json:"duetime"`  //开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
	Interval     int64        `json:"interval"` //运行间隔时间，单位毫秒，当TaskType==TaskType_Loop时有效
	RawExpress   string       `json:"express"`  //运行周期表达式，当TaskType==TaskType_Cron时有效
	time_WeekDay *ExpressSet
	time_Month   *ExpressSet
	time_Day     *ExpressSet
	time_Hour    *ExpressSet
	time_Minute  *ExpressSet
	time_Second  *ExpressSet
}

//start task
func (task *TaskInfo) Start() {
	if !task.IsRun {
		return
	}

	task.mutex.Lock()
	defer task.mutex.Unlock()

	if task.State == TaskState_Init || task.State == TaskState_Stop {
		task.State = TaskState_Run
		switch task.TaskType {
		case TaskType_Cron:
			startCronTask(task)
		case TaskType_Loop:
			startLoopTask(task)
		default:
			panic("not support task_type => " + task.TaskType)
		}
	}
}

//stop task
func (task *TaskInfo) Stop() {
	if !task.IsRun {
		return
	}
	if task.State == TaskState_Stop {
		task.TimeTicker.Stop()
		task.State = TaskState_Stop
		task.taskService.Logger().Debug(task.TaskID, " Stop")
	}
}

// RunOnce do task only once
// no match Express or Interval
// no recover panic
// support for #6 新增RunOnce方法建议
func (task *TaskInfo) RunOnce() error {
	err := task.handler(task.Context)
	return err
}

//start cron task
func startCronTask(task *TaskInfo) {
	now := time.Now()
	nowsecond := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	afterTime := nowsecond.Add(time.Second).Sub(time.Now().Local())
	task.TimeTicker = time.NewTicker(DefaultPeriod)
	go func() {
		time.Sleep(afterTime)
		for {
			select {
			case <-task.TimeTicker.C:
				defer func() {
					if err := recover(); err != nil {
						//task.taskService.Logger().Debug(task.TaskID, " cron handler recover error => ", err)
						if task.taskService.ExceptionHandler != nil {
							task.taskService.ExceptionHandler(task.Context, fmt.Errorf("%v", err))
						}
					}
				}()
				now := time.Now()
				if task.time_WeekDay.IsMatch(now) &&
					task.time_Month.IsMatch(now) &&
					task.time_Day.IsMatch(now) &&
					task.time_Hour.IsMatch(now) &&
					task.time_Minute.IsMatch(now) &&
					task.time_Second.IsMatch(now) {
					//do log
					//task.taskService.Logger().Debug(task.TaskID, " begin dohandler")
					if task.taskService.OnBeforHandler != nil {
						task.taskService.OnBeforHandler(task.Context)
					}
					var err error
					if !task.Context.IsEnd {
						err = task.handler(task.Context)
					}
					if err != nil {
						if task.taskService.ExceptionHandler != nil {
							task.taskService.ExceptionHandler(task.Context, err)
						}
					}
					if task.taskService.OnEndHandler != nil {
						task.taskService.OnEndHandler(task.Context)
					}
				}
			}
		}
	}()
}

//start loop task
func startLoopTask(task *TaskInfo) {
	handler := func() {
		defer func() {
			if err := recover(); err != nil {
				//task.taskService.Logger().Debug(task.TaskID, " loop handler recover error => ", err)
				if task.taskService.ExceptionHandler != nil {
					task.taskService.ExceptionHandler(task.Context, fmt.Errorf("%v", err))
				}
			}
		}()
		//do log
		if task.taskService.OnBeforHandler != nil {
			task.taskService.OnBeforHandler(task.Context)
		}
		var err error
		if !task.Context.IsEnd {
			err = task.handler(task.Context)
		}
		if err != nil {
			if task.taskService.ExceptionHandler != nil {
				task.taskService.ExceptionHandler(task.Context, err)
			}
		} else {
			//task.taskService.Logger().Debug(task.TaskID, " loop handler end success")
		}
		if task.taskService.OnEndHandler != nil {
			task.taskService.OnEndHandler(task.Context)
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
