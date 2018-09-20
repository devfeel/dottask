package task

import (
	"time"
	"fmt"
	"strings"
	"errors"
)

type (
	//CronTask cron task info define
	CronTask struct {
		TaskInfo
		RawExpress   string       `json:"express"`  //运行周期表达式，当TaskType==TaskType_Cron时有效
		time_WeekDay *ExpressSet
		time_Month   *ExpressSet
		time_Day     *ExpressSet
		time_Hour    *ExpressSet
		time_Minute  *ExpressSet
		time_Second  *ExpressSet
	}
)

// GetConfig get task config info
func (task *CronTask) GetConfig() *TaskConfig{
		return &TaskConfig{
			TaskID:task.taskID,
			TaskType:task.TaskType,
			IsRun : task.IsRun,
			Handler:task.handler,
			DueTime:task.DueTime,
			Interval:0,
			Express:task.RawExpress,
			TaskData:task.Context().TaskData,
		}
}


//Reset first check conf, then reload conf & restart task
//special, TaskID can not be reset
//special, if TaskData is nil, it can not be reset
//special, if Handler is nil, it can not be reset
func (task *CronTask) Reset(conf *TaskConfig) error {
	expresslist := strings.Split(conf.Express, " ")

	//basic check
	if conf.Express == "" {
		errmsg := "express is empty"
		task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "error", errmsg)
		return errors.New(errmsg)
	}
	if len(expresslist) != 6 {
		errmsg := "express is wrong format => not 6 parts"
		task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "error", errmsg)
		return errors.New("express is wrong format => not 6 parts")
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
	task.RawExpress = conf.Express
	if task.TaskType == TaskType_Cron {
		task.time_WeekDay = parseExpress(expresslist[5], ExpressType_WeekDay)
		task.taskService.debugExpress(task.time_WeekDay)
		task.time_Month = parseExpress(expresslist[4], ExpressType_Month)
		task.taskService.debugExpress(task.time_Month)
		task.time_Day = parseExpress(expresslist[3], ExpressType_Day)
		task.taskService.debugExpress(task.time_Day)
		task.time_Hour = parseExpress(expresslist[2], ExpressType_Hour)
		task.taskService.debugExpress(task.time_Hour)
		task.time_Minute = parseExpress(expresslist[1], ExpressType_Minute)
		task.taskService.debugExpress(task.time_Minute)
		task.time_Second = parseExpress(expresslist[0], ExpressType_Second)
		task.taskService.debugExpress(task.time_Second)
	}
	task.Start()
	task.taskService.Logger().Debug("TaskInfo:Reset ", task, conf, "success")
	return nil
}


//Start start task
func (task *CronTask) Start() {
	if !task.IsRun {
		return
	}

	task.mutex.Lock()
	defer task.mutex.Unlock()

	if task.State == TaskState_Init || task.State == TaskState_Stop {
		task.State = TaskState_Run
		startCronTask(task)
	}
}


// RunOnce do task only once
// no match Express or Interval
// no recover panic
// support for #6 新增RunOnce方法建议
func (task *CronTask) RunOnce() error {
	err := task.handler(task.context)
	return err
}





//start cron task
func startCronTask(task *CronTask) {
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
						task.taskService.Logger().Debug(task.TaskID, " cron handler recover error => ", err)
						if task.taskService.ExceptionHandler != nil {
							task.taskService.ExceptionHandler(task.Context(), fmt.Errorf("%v", err))
						}
						//goroutine panic, restart cron task
						startCronTask(task)
						task.taskService.Logger().Debug(task.TaskID, " goroutine panic, restart CronTask")
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
					}
					if task.taskService.OnEndHandler != nil {
						task.taskService.OnEndHandler(task.Context())
					}
				}
			}
		}
	}()
}

