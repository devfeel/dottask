package task

const (
	packageVersion = "0.9.10"
)

const (
	TaskState_Init = "0"
	TaskState_Run  = "1"
	TaskState_Stop = "2"

	TaskType_Loop  = "loop"
	TaskType_Cron  = "cron"
	TaskType_Queue = "queue"

	ConfigType_Xml  = "xml"
	ConfigType_Json = "json"
	ConfigType_Yaml = "yaml"
)

const (
	defaultCounterTaskName = "TASK_DEFAULT_COUNTERINFO"
	fullTimeLayout         = "2006-01-02 15:04:05.9999"
)
