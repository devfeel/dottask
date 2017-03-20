package task

import (
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
	taskService  *TaskService
	mutex        sync.RWMutex
	TimeTicker   *time.Ticker `json:"-"`
	TaskType     string       `json:"tasktype"`
	handler      TaskHandle
	Context      *TaskContext `json:"context"`
	State        string       `json:"state"`    //匹配 TskState_Init、TaskState_Run、TaskState_Stop
	Interval     int64        `json:"interval"` //运行间隔时间，单位毫秒，当TaskType==TaskType_Loop时有效
	RawExpress   string       `json:"express"`  //运行周期表达式，当TaskType==TaskType_Cron时有效
	time_WeekDay *ExpressSet
	time_Month   *ExpressSet
	time_Day     *ExpressSet
	time_Hour    *ExpressSet
	time_Minute  *ExpressSet
	time_Second  *ExpressSet
}

//start crontask
func (task *TaskInfo) Start() {
	task.mutex.Lock()
	defer task.mutex.Unlock()

	if task.State == TaskState_Init || task.State == TaskState_Stop {
		task.State = TaskState_Run
		switch task.TaskType {
		case TaskType_Cron:
			now := time.Now()
			nowsecond := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
			afterTime := nowsecond.Add(time.Second).Sub(time.Now().Local())
			dofunc := func() {
				task.TimeTicker = time.NewTicker(DefaultPeriod)
				go func() {
					for {
						select {
						case <-task.TimeTicker.C:
							defer func() {
								if err := recover(); err != nil {
									task.taskService.Debug(task.TaskID, " cron handler recover error => ", err)
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
								//task.taskService.logger.Debug(task.TaskID, " begin dohandler")
								err := task.handler(task.Context)
								if err != nil {
									task.taskService.Debug(task.TaskID, " cron handler failed => "+err.Error())
								} else {
									//task.taskService.logger.Debug(task.TaskID, " cron handler end success")
								}
							}
						}
					}
				}()
			}
			time.AfterFunc(afterTime, dofunc)
		case TaskType_Loop:
			task.TimeTicker = time.NewTicker(time.Duration(task.Interval) * time.Millisecond)
			go func() {
				for {
					select {
					case <-task.TimeTicker.C:
						defer func() {
							if err := recover(); err != nil {
								task.taskService.Debug(task.TaskID, " loop handler recover error => ", err)
							}
						}()
						//do log
						//task.taskService.logger.Debug(task.TaskID, " loop handler begin")
						err := task.handler(task.Context)
						if err != nil {
							task.taskService.Debug(task.TaskID, " loop handler failed => "+err.Error())
						} else {
							//task.taskService.logger.Debug(task.TaskID, " loop handler end success")
						}
					}
				}
			}()
		default:
			panic("not support task_type => " + task.TaskType)
		}
	}
}

//stop taskinfo
func (task *TaskInfo) Stop() {
	if task.State == TaskState_Stop {
		task.TimeTicker.Stop()
		task.State = TaskState_Stop
	}
}
