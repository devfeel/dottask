package task

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
)

const (
	TaskState_Init = "0"
	TaskState_Run  = "1"
	TaskState_Stop = "2"
)

const (
	TaskType_Loop = "loop"
	TaskType_Cron = "cron"
)

type (
	Task interface {
		Start() error
		Stop() error
	}

	//Task上下文信息
	TaskContext struct {
		TaskID   string
		TaskData interface{}
	}

	TaskHandle func(*TaskContext) error

	//task 容器
	TaskService struct {
		taskMap   map[string]*TaskInfo
		taskMutex *sync.RWMutex
		logger    Logger
	}

	Logger interface {
		Error(v ...interface{})
		Warn(v ...interface{})
		Info(v ...interface{})
		Debug(v ...interface{})
	}
)

func StartNewService() *TaskService {
	service := new(TaskService)
	service.taskMutex = new(sync.RWMutex)
	service.taskMap = make(map[string]*TaskInfo)

	return service
}

//set logger which Implements Logger interface
func (service *TaskService) SetLogger(logger Logger) {
	service.logger = logger
}

//create new crontask
func (service *TaskService) CreateCronTask(taskID string, express string, handler TaskHandle, taskData interface{}) (*TaskInfo, error) {
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(TaskInfo)
	task.TaskID = context.TaskID
	task.TaskType = TaskType_Cron
	task.handler = handler
	task.RawExpress = express
	expresslist := strings.Split(express, " ")
	if len(expresslist) != 6 {
		return nil, errors.New("express is wrong format => not 6 part")
	}
	task.time_WeekDay = parseExpress(expresslist[5], ExpressType_WeekDay)
	service.debugExpress(task.time_WeekDay)
	task.time_Month = parseExpress(expresslist[4], ExpressType_Month)
	service.debugExpress(task.time_Month)
	task.time_Day = parseExpress(expresslist[3], ExpressType_Day)
	service.debugExpress(task.time_Day)
	task.time_Hour = parseExpress(expresslist[2], ExpressType_Hour)
	service.debugExpress(task.time_Hour)
	task.time_Minute = parseExpress(expresslist[1], ExpressType_Minute)
	service.debugExpress(task.time_Minute)
	task.time_Second = parseExpress(expresslist[0], ExpressType_Second)
	service.debugExpress(task.time_Second)

	task.State = TaskState_Init
	task.Context = context

	service.AddTask(task)

	return task, nil
}

//create new looptask
func (service *TaskService) CreateLoopTask(taskID string, interval int64, handler TaskHandle, taskData interface{}) (*TaskInfo, error) {
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(TaskInfo)
	task.TaskID = context.TaskID
	task.TaskType = TaskType_Loop
	task.handler = handler
	task.Interval = interval
	task.State = TaskState_Init
	task.Context = context

	service.AddTask(task)
	return task, nil
}

//add new task point
func (service *TaskService) AddTask(t *TaskInfo) {
	service.taskMutex.Lock()
	service.taskMap[t.TaskID] = t
	service.taskMutex.Unlock()
	t.taskService = service
	service.logger.Debug("Task:AddTask => ", t.TaskID)
}

//remove task by taskid
func (service *TaskService) RemoveTask(taskID string) {
	service.taskMutex.Lock()
	delete(service.taskMap, taskID)
	service.taskMutex.Unlock()
	service.logger.Debug("Task:RemoveTask => ", taskID)
}

//get all task's count
func (service *TaskService) Count() int {
	return len(service.taskMap)
}

//print all crontask
func (service *TaskService) PrintAllCronTask() string {
	body := ""
	for _, v := range service.taskMap {
		str, _ := json.Marshal(v)
		body += string(str) + "\r\n"
	}
	return body

}

//remove all task
func (service *TaskService) RemoveAllTask() {
	service.StopAllTask()
	service.taskMap = make(map[string]*TaskInfo)
	service.logger.Debug("Task:RemoveAllTask")
}

//stop all task
func (service *TaskService) StopAllTask() {
	service.logger.Info("Task:StopAllTask begin...")
	for _, v := range service.taskMap {
		service.logger.Debug("Task:StopAllTask::StopTask => ", v.TaskID)
		v.Stop()
	}
	service.logger.Info("Task:StopAllTask end[" + string(len(service.taskMap)) + "]")
}

//start all task
func (service *TaskService) StartAllTask() {
	service.logger.Info("Task:StartAllTask begin...")
	for _, v := range service.taskMap {
		service.logger.Debug("Task:StartAllTask::StartTask => " + v.TaskID)
		v.Start()
	}
	service.logger.Info("Task:StartAllTask end[" + strconv.Itoa(len(service.taskMap)) + "]")
}

func (service *TaskService) debugExpress(set *ExpressSet) {
	service.logger.Debug("parseExpress(", set.rawExpress, " , ", set.expressType, ") => ", set.timeMap)
}
