package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	TaskType_Loop  = "loop"
	TaskType_Cron  = "cron"
	TaskType_Queue = "queue"
)

const (
	ConfigType_Xml  = "xml"
	ConfigType_Json = "json"
	ConfigType_Yaml = "yaml"
)

const (
	defaultCounterTaskName = "TASK_DEFAULT_COUNTERINFO"
	fullTimeLayout         = "2006-01-02 15:04:05.9999"
)

var ErrNotSupportTaskType = errors.New("not support task type")

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
		CounterInfo() *CounterInfo
	}

	ExceptionHandleFunc func(*TaskContext, error)

	//task 容器
	TaskService struct {
		Config           *AppConfig
		taskMap          map[string]Task
		taskMutex        *sync.RWMutex
		counterTask      Task
		logger           Logger
		handlerMap       map[string]TaskHandle
		handlerMutex     *sync.RWMutex
		ExceptionHandler ExceptionHandleFunc
		OnBeforeHandler  TaskHandle
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

// UseDefaultLogCounterTask use default LogCounterTask in TaskService
func (service *TaskService) UseDefaultLogCounterTask() {
	counterTask, err := NewCronTask(defaultCounterTaskName, true, "0 * * * * *", service.defaultLogCounterInfo, nil)
	if err != nil {
		fmt.Println("Init Counter Task error", err)
	}
	counterTask.SetTaskService(service)
	counterTask.Start()
	service.counterTask = counterTask
}

// SetExceptionHandler 设置自定义异常处理方法
func (service *TaskService) SetExceptionHandler(handler ExceptionHandleFunc) {
	service.ExceptionHandler = handler
}

// SetOnBeforHandler set handler which exec before task run
func (service *TaskService) SetOnBeforeHandler(handler TaskHandle) {
	service.OnBeforeHandler = handler
}

// SetOnEndHandler set handler which exec after task run
func (service *TaskService) SetOnEndHandler(handler TaskHandle) {
	service.OnEndHandler = handler
}

// LoadConfig 如果指定配置文件，初始化配置
// Deprecated: Use the LoadFileConfig instead
func (service *TaskService) LoadConfig(configFile string, confType ...interface{}) *TaskService {
	return service.LoadFileConfig(configFile, confType...)
}

// LoadFileConfig 如果指定配置文件，初始化配置
func (service *TaskService) LoadFileConfig(configFile string, confType ...interface{}) *TaskService {
	cType := ConfigType_Xml
	if len(confType) > 0 && confType[0] == ConfigType_Json {
		cType = ConfigType_Json
	}
	if len(confType) > 0 && confType[0] == ConfigType_Yaml {
		cType = ConfigType_Yaml
	}
	var config *AppConfig
	if cType == ConfigType_Json {
		config = JsonConfigHandler(configFile)
	} else if cType == ConfigType_Yaml {
		config = YamlConfigHandler(configFile)
	} else {
		config = XmlConfigHandler(configFile)
	}
	service.applyConfig(config)
	return service
}

// applyConfig apply task config with AppConfig
func (service *TaskService) applyConfig(config *AppConfig) *TaskService {
	service.Config = config
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
			} else if v.TaskType == TaskType_Queue && v.Interval > 0 {
				_, err := service.CreateQueueTask(v.TaskID, v.IsRun, v.Interval, handler, v, v.QueueSize)
				if err != nil {
					service.Logger().Warn("CreateQueueTask failed [" + err.Error() + "] [" + fmt.Sprint(v) + "]")
				} else {
					service.Logger().Debug("CreateQueueTask success [" + fmt.Sprint(v) + "]")
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

// LoadConfigHandler load config handler and init task config
func (service *TaskService) LoadConfigHandler(configHandler ConfigHandle, configSource string) *TaskService {
	config, err := configHandler(configSource)
	if err != nil {
		panic("Task:LoadConfigHandler 配置源[" + configSource + "]解析失败: " + err.Error())
		os.Exit(1)
	} else {
		service.applyConfig(config)
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

// CreateTask create new task with TaskConfig and register to task service
func (service *TaskService) CreateTask(config TaskConfig) (Task, error) {
	if config.TaskType == TaskType_Cron {
		return service.CreateCronTask(config.TaskID, config.IsRun, config.Express, config.Handler, config.TaskData)
	} else if config.TaskType == TaskType_Loop {
		return service.CreateLoopTask(config.TaskID, config.IsRun, config.DueTime, config.Interval, config.Handler, config.TaskData)
	} else if config.TaskType == TaskType_Loop {
		return service.CreateLoopTask(config.TaskID, config.IsRun, config.DueTime, config.Interval, config.Handler, config.TaskData)
	}
	return nil, ErrNotSupportTaskType
}

// CreateCronTask create new cron task and register to task service
func (service *TaskService) CreateCronTask(taskID string, isRun bool, express string, handler TaskHandle, taskData interface{}) (Task, error) {
	task, err := NewCronTask(taskID, isRun, express, handler, taskData)
	if err != nil {
		return task, err
	}

	//service.debugExpress(task.(*CronTask).time_WeekDay)
	//service.debugExpress(task.(*CronTask).time_Month)
	//service.debugExpress(task.(*CronTask).time_Day)
	//service.debugExpress(task.(*CronTask).time_Hour)
	//service.debugExpress(task.(*CronTask).time_Minute)
	//service.debugExpress(task.(*CronTask).time_Second)

	task.SetTaskService(service)
	service.AddTask(task)
	return task, nil
}

// CreateLoopTask create new loop task and register to task service
func (service *TaskService) CreateLoopTask(taskID string, isRun bool, dueTime int64, interval int64, handler TaskHandle, taskData interface{}) (Task, error) {
	task, err := NewLoopTask(taskID, isRun, dueTime, interval, handler, taskData)
	if err != nil {
		return task, err
	}
	task.SetTaskService(service)
	service.AddTask(task)
	return task, nil
}

// CreateQueueTask create new queue task and register to task service
func (service *TaskService) CreateQueueTask(taskID string, isRun bool, interval int64, handler TaskHandle, taskData interface{}, queueSize int64) (Task, error) {
	task, err := NewQueueTask(taskID, isRun, interval, handler, taskData, queueSize)
	if err != nil {
		return task, err
	}

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

	service.Logger().Info(fmt.Sprint("Task:AddTask => ", t.TaskID(), t.GetConfig()))
}

// RemoveTask remove task by taskid
func (service *TaskService) RemoveTask(taskID string) {
	service.taskMutex.Lock()
	delete(service.taskMap, taskID)
	service.taskMutex.Unlock()
	service.Logger().Info(fmt.Sprint("Task:RemoveTask => ", taskID))
}

// Count get all task's count
func (service *TaskService) Count() int {
	return len(service.taskMap)
}

// PrintAllCronTask print all task info
// Deprecated: Use the PrintAllTasks instead
func (service *TaskService) PrintAllCronTask() string {
	body := ""
	for _, v := range service.taskMap {
		str, _ := json.Marshal(v)
		body += string(str) + "\r\n"
	}
	return body
}

// PrintAllCronTask print all task info
func (service *TaskService) PrintAllTasks() string {
	body := ""
	for _, v := range service.taskMap {
		str, _ := json.Marshal(v)
		body += string(str) + "\r\n"
	}
	return body
}

// PrintTaskCountData print all task counter data
func (service *TaskService) PrintAllTaskCounterInfo() string {
	body := ""
	for _, v := range service.taskMap {
		body += fmt.Sprintln(v.TaskID(), "Run", v.CounterInfo().RunCounter.Count())
		body += fmt.Sprintln(v.TaskID(), "Error", v.CounterInfo().ErrorCounter.Count())
	}
	return body
}

// GetAllTaskCounterInfo return all show count info
func (service *TaskService) GetAllTaskCountInfo() []ShowCountInfo {
	showInfos := []ShowCountInfo{}
	for _, t := range service.GetAllTasks() {
		showInfos = append(showInfos, ShowCountInfo{TaskID: t.TaskID(), Lable: "RUN", Count: t.CounterInfo().RunCounter.Count()})
		showInfos = append(showInfos, ShowCountInfo{TaskID: t.TaskID(), Lable: "ERROR", Count: t.CounterInfo().ErrorCounter.Count()})
	}
	return showInfos
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
		service.Logger().Info(fmt.Sprint("Task:StopAllTask::StopTask => ", v.TaskID()))
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
	service.Logger().Debug(fmt.Sprint("parseExpress(", set.rawExpress, " , ", set.expressType, ") => ", set.timeMap))
}

func (service *TaskService) defaultLogCounterInfo(ctx *TaskContext) error {
	showInfos := []*ShowCountInfo{}
	for _, t := range service.GetAllTasks() {
		//service.Logger().Debug(t.TaskID(), " Start:", t.CounterInfo().StartTime.Format(fullTimeLayout), " RUN:", t.CounterInfo().RunCounter.Count(), " ERROR:", t.CounterInfo().ErrorCounter.Count())
		showInfos = append(showInfos, &ShowCountInfo{TaskID: t.TaskID(), Lable: "RUN", Count: t.CounterInfo().RunCounter.Count()})
		showInfos = append(showInfos, &ShowCountInfo{TaskID: t.TaskID(), Lable: "ERROR", Count: t.CounterInfo().ErrorCounter.Count()})
	}
	if len(showInfos) > 0 {
		str, _ := json.Marshal(showInfos)
		service.Logger().Debug(string(str))
	}
	return nil
}

// ValidateTaskType validate the TaskType is supported
func ValidateTaskType(taskType string) bool {
	if taskType == "" {
		return false
	}
	checkType := strings.ToLower(taskType)
	if checkType != TaskType_Cron && checkType != TaskType_Loop && checkType != TaskType_Queue {
		return false
	}
	return true
}
