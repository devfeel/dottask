package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
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
	TaskType_Queue = "queue"
)

const (
	ConfigType_Xml  = "xml"
	ConfigType_Json = "json"
	ConfigType_Yaml = "yaml"
)

type (
	Task interface {
		TaskID() string
		GetConfig() *TaskConfig
		Context() *TaskContext
		Start()
		Stop()
		RunOnce() error
		SetTaskService(service *TaskService)
		Reset(conf *TaskConfig) error
	}



	ExceptionHandleFunc func(*TaskContext, error)

	//task 容器
	TaskService struct {
		Config           *AppConfig
		taskMap          map[string]Task
		taskMutex        *sync.RWMutex
		logger           Logger
		handlerMap       map[string]TaskHandle
		handlerMutex     *sync.RWMutex
		ExceptionHandler ExceptionHandleFunc
		OnBeforHandler   TaskHandle
		OnEndHandler     TaskHandle
	}
)

func defaultExceptionHandler(ctx *TaskContext, err error) {
	stack := string(debug.Stack())
	fmt.Println("[", ctx.TaskID, "] ", ctx.TaskData, " error [", err.Error(), "] => ", stack)
}

func StartNewService() *TaskService {
	service := new(TaskService)
	service.taskMutex = new(sync.RWMutex)
	service.taskMap = make(map[string]Task)
	service.handlerMutex = new(sync.RWMutex)
	service.handlerMap = make(map[string]TaskHandle)
	service.ExceptionHandler = defaultExceptionHandler
	return service
}

// SetExceptionHandler 设置自定义异常处理方法
func (service *TaskService) SetExceptionHandler(handler ExceptionHandleFunc) {
	service.ExceptionHandler = handler
}

// SetOnBeforHandler set handler which exec before task run
func (service *TaskService) SetOnBeforHandler(handler TaskHandle) {
	service.OnBeforHandler = handler
}

// SetOnEndHandler set handler which exec after task run
func (service *TaskService) SetOnEndHandler(handler TaskHandle) {
	service.OnEndHandler = handler
}

// LoadConfig 如果指定配置文件，初始化配置
func (service *TaskService) LoadConfig(configFile string, confType ...interface{}) *TaskService {
	cType := ConfigType_Xml
	if len(confType) > 0 && confType[0] == ConfigType_Json {
		cType = ConfigType_Json
	}
	if len(confType) > 0 && confType[0] == ConfigType_Yaml {
		cType = ConfigType_Yaml
	}
	if cType == ConfigType_Json {
		service.Config = InitJsonConfig(configFile)
	} else if cType == ConfigType_Yaml {
		service.Config = InitYamlConfig(configFile)
	} else {
		service.Config = InitConfig(configFile)
	}
	if service.logger == nil {
		if service.Config.Global.LogPath != "" {
			service.SetLogger(NewFileLogger(service.Config.Global.LogPath))
		} else {
			service.SetLogger(NewFmtLogger())
		}
	}
	for _, v := range service.Config.Tasks {
		if handler, exists := service.GetHandler(v.HandlerName); exists {
			if v.TaskType == TaskType_Cron && v.Express != "" {
				_, err := service.CreateCronTask(v.TaskID, v.IsRun, v.Express, handler, v)
				if err != nil {
					service.Logger().Warn("CreateCronTask failed [" + err.Error() + "] [" + fmt.Sprint(v) + "]")
				} else {
					service.Logger().Debug("CreateCronTask success [" + fmt.Sprint(v) + "]")
				}
			} else if v.TaskType == TaskType_Loop && v.Interval > 0 {
				_, err := service.CreateLoopTask(v.TaskID, v.IsRun, v.DueTime, v.Interval, handler, v)
				if err != nil {
					service.Logger().Warn("CreateLoopTask failed [" + err.Error() + "] [" + fmt.Sprint(v) + "]")
				} else {
					service.Logger().Debug("CreateLoopTask success [" + fmt.Sprint(v) + "]")
				}
			} else {
				service.Logger().Warn("CreateTask failed not match config [" + fmt.Sprint(v) + "]")
			}

		} else {
			service.Logger().Warn("CreateTask failed not exists handler [" + fmt.Sprint(v) + "]")
		}
	}
	return service
}

// RegisterHandler register handler by name
func (service *TaskService) RegisterHandler(name string, handler TaskHandle) *TaskService {
	service.handlerMutex.Lock()
	service.handlerMap[name] = handler
	service.handlerMutex.Unlock()
	return service
}

// GetHandler get handler by handler name
func (service *TaskService) GetHandler(name string) (TaskHandle, bool) {
	service.handlerMutex.RLock()
	v, exists := service.handlerMap[name]
	service.handlerMutex.RUnlock()
	return v, exists
}

// SetLogger set logger which Implements Logger interface
func (service *TaskService) SetLogger(logger Logger) {
	service.logger = logger
}

// Logger get Logger
func (service *TaskService) Logger() Logger {
	if service.logger == nil {
		service.logger = NewFmtLogger()
	}
	return service.logger
}

// CreateCronTask create new crontask
func (service *TaskService) CreateCronTask(taskID string, isRun bool, express string, handler TaskHandle, taskData interface{}) (Task, error) {
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(CronTask)
	task.taskID = context.TaskID
	task.TaskType = TaskType_Cron
	task.IsRun = isRun
	task.handler = handler
	task.RawExpress = express
	expresslist := strings.Split(express, " ")
	if len(expresslist) != 6 {
		return nil, errors.New("express is wrong format => not 6 parts")
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
	task.context = context
	task.SetTaskService(service)

	service.AddTask(task)
	return task, nil
}

// CreateLoopTask create new looptask
func (service *TaskService) CreateLoopTask(taskID string, isRun bool, dueTime int64, interval int64, handler TaskHandle, taskData interface{}) (Task, error) {
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(LoopTask)
	task.taskID = context.TaskID
	task.TaskType = TaskType_Loop
	task.IsRun = isRun
	task.handler = handler
	task.DueTime = dueTime
	task.Interval = interval
	task.State = TaskState_Init
	task.context = context
	task.SetTaskService(service)

	service.AddTask(task)
	return task, nil
}

// CreateQueueTask create new queuetask
func (service *TaskService) CreateQueueTask(taskID string, isRun bool, interval int64, handler TaskHandle, taskData interface{}, queueSize int) (Task, error){
	context := new(TaskContext)
	context.TaskID = taskID
	context.TaskData = taskData

	task := new(QueueTask)
	task.taskID = context.TaskID
	task.TaskType = TaskType_Queue
	task.IsRun = isRun
	task.handler = handler
	task.Interval = interval
	task.State = TaskState_Init
	task.context = context
	task.MessageChan = make(chan interface{}, queueSize)
	task.SetTaskService(service)

	service.AddTask(task)
	return task, nil
}



// GetTask get TaskInfo by TaskID
func (service *TaskService) GetTask(taskID string) (t Task, exists bool) {
	service.taskMutex.RLock()
	defer service.taskMutex.RUnlock()
	t, exists = service.taskMap[taskID]
	return t, exists
}

// AddTask add new task point
func (service *TaskService) AddTask(t Task) {
	service.taskMutex.Lock()
	service.taskMap[t.TaskID()] = t
	service.taskMutex.Unlock()

	service.Logger().Info("Task:AddTask => ", t.TaskID(), t.GetConfig())
}

// RemoveTask remove task by taskid
func (service *TaskService) RemoveTask(taskID string) {
	service.taskMutex.Lock()
	delete(service.taskMap, taskID)
	service.taskMutex.Unlock()
	service.Logger().Info("Task:RemoveTask => ", taskID)
}

// Count get all task's count
func (service *TaskService) Count() int {
	return len(service.taskMap)
}

// PrintAllCronTask print all task
func (service *TaskService) PrintAllCronTask() string {
	body := ""
	for _, v := range service.taskMap {
		str, _ := json.Marshal(v)
		body += string(str) + "\r\n"
	}
	return body
}


// GetAllTasks get all tasks
func (service *TaskService) GetAllTasks() map[string]Task {
	return service.taskMap
}


// RemoveAllTask remove all task
func (service *TaskService) RemoveAllTask() {
	service.StopAllTask()
	service.taskMap = make(map[string]Task)
	service.Logger().Info("Task:RemoveAllTask")
}

// StopAllTask stop all task
func (service *TaskService) StopAllTask() {
	service.Logger().Info("Task:StopAllTask begin...")
	for _, v := range service.taskMap {
		service.Logger().Info("Task:StopAllTask::StopTask => ", v.TaskID())
		v.Stop()
	}
	service.Logger().Info("Task:StopAllTask end[" + string(len(service.taskMap)) + "]")
}

// StartAllTask start all task
func (service *TaskService) StartAllTask() {
	service.Logger().Info("Task:StartAllTask begin...")
	for _, v := range service.taskMap {
		service.Logger().Info("Task:StartAllTask::StartTask => " + v.TaskID())
		v.Start()
	}
	service.Logger().Info("Task:StartAllTask end[" + strconv.Itoa(len(service.taskMap)) + "]")
}

func (service *TaskService) debugExpress(set *ExpressSet) {
	service.Logger().Debug("parseExpress(", set.rawExpress, " , ", set.expressType, ") => ", set.timeMap)
}
