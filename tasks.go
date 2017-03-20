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
		taskMap      map[string]*TaskInfo
		taskMutex    *sync.RWMutex
		logger       Logger
		handlerMap   map[string]TaskHandle
		handlerMutex *sync.RWMutex
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
	service.handlerMutex = new(sync.RWMutex)
	service.handlerMap = make(map[string]TaskHandle)

	return service
}

func (service *TaskService) RegisterHandler(name string, handler TaskHandle) {
	service.handlerMutex.Lock()
	service.handlerMap[name] = handler
	service.handlerMutex.Unlock()
}

func (service *TaskService) GetHandler(name string) (TaskHandle, bool) {
	service.handlerMutex.RLock()
	v, exists := service.handlerMap[name]
	service.handlerMutex.RUnlock()
	return v, exists
}

//set logger which Implements Logger interface
func (service *TaskService) SetLogger(logger Logger) {
	service.logger = logger
}

func (service *TaskService) Debug(v ...interface{}) {
	if service.logger != nil {
		service.logger.Debug(v...)
	}
}

func (service *TaskService) Info(v ...interface{}) {
	if service.logger != nil {
		service.logger.Info(v...)
	}
}

func (service *TaskService) Warn(v ...interface{}) {
	if service.logger != nil {
		service.logger.Warn(v...)
	}
}

func (service *TaskService) Error(v ...interface{}) {
	if service.logger != nil {
		service.logger.Error(v...)
	}
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
	service.Debug("Task:AddTask => ", t.TaskID)
}

//remove task by taskid
func (service *TaskService) RemoveTask(taskID string) {
	service.taskMutex.Lock()
	delete(service.taskMap, taskID)
	service.taskMutex.Unlock()
	service.Debug("Task:RemoveTask => ", taskID)
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
	service.Debug("Task:RemoveAllTask")
}

//stop all task
func (service *TaskService) StopAllTask() {
	service.Info("Task:StopAllTask begin...")
	for _, v := range service.taskMap {
		service.Debug("Task:StopAllTask::StopTask => ", v.TaskID)
		v.Stop()
	}
	service.Info("Task:StopAllTask end[" + string(len(service.taskMap)) + "]")
}

//start all task
func (service *TaskService) StartAllTask() {
	service.Info("Task:StartAllTask begin...")
	for _, v := range service.taskMap {
		service.Debug("Task:StartAllTask::StartTask => " + v.TaskID)
		v.Start()
	}
	service.Info("Task:StartAllTask end[" + strconv.Itoa(len(service.taskMap)) + "]")
}

func (service *TaskService) debugExpress(set *ExpressSet) {
	service.Debug("parseExpress(", set.rawExpress, " , ", set.expressType, ") => ", set.timeMap)
}
