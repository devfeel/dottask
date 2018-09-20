## devfeel/task

#### Version 0.8.4
* Bug Fixed: Auto restart task when task's app code panic error, if no auto restart, the goroutine will be stop
* 2018-09-20 22:00

#### Version 0.8.3
* New feature: add TaskService.GetAllTasks() used to return all task info which type is map[string]Task
* Example:
  ```golang
  for _, t := range service.GetAllTasks(){
  		fmt.Println("GetAllTasks", t.TaskID(), t.GetConfig().TaskType, t.GetConfig().IsRun, t.GetConfig().Interval, t.GetConfig().Express)
  }
  ```
* 2018-09-19 12:00

#### Version 0.8.2
* 调整AddTask\RemoveTask\StartTask\StopTask\StopAllTask日志级别为info级别
* 调整默认日志文件为dottask_{time}.log
* Task新增GetConfig接口，用于返回当前task的设置信息，使用TaskConfig包装
* 2018-05-18 16:00

#### Version 0.8.1
* 修复BUG #8：ExpressType_WeekDay下未设置周日情况下自动添加周日运行
* 2018-03-13 13:00

#### Version 0.8
* 重构TaskInfo,分拆CronTask\LoopTask,增加Task interface
* 重构TaskService,所有*TaskInfo调整为Task interface
* 新增QueueTask,适用于内存队列任务的简单模型
* 特别的,新增TaskContext.Message：用于每次Task执行上下文消息传输
* 特别的,TaskContext.TaskData：用于当前Task全局设置的数据项
* 新增example/queuetask
* 修复Task.Stop, Task.Reset隐藏bug
* 2017-11-27 10:00

#### Version 0.7
* TaskInfo增加Reset方法，用于重设Task相关参数并重启Task
* 新增TaskConfig结构，用于描述Task基本设置
* 特别的,TaskID无法被覆盖
* 特别的,TaskConfig.Handler为nil时，不会覆盖Task老的Handler
* 特别的,TaskConfig.TaskData为nil时，不会覆盖Task老的TaskData
* Support for #7 建议增加一个Reset方法，用于重置某个task
* 2017-10-30 12:00

#### Version 0.6
* TaskInfo增加GetTask、RunOnce方法
* GetTask：根据TaskID获取*TaskInfo
* RunOnce：仅执行一次TaskInfo对应的Handler，不匹配任何表达式，不捕获未处理异常
* Support for #6 新增RunOnce方法建议
* 2017-10-19 09:00

#### Version 0.5
* 配置增加yaml文件支持
* 新增 example\loadfromyaml 示例
* 贡献者：https://github.com/yulibaozi
* 2017-08-20 09:00

#### Version 0.4
* 增加ExceptionHandler\OnBeforHandler\OnEndHandler三个注入点
* 新增ExceptionHandleFunc func(*TaskContext, error) 定义
* 内建defaultExceptionHandler支持
* 更新 example\normal
* 2017-05-17 19:00

#### Version 0.3
* 配置增加json文件支持
* 支持链式调用
* 更新 example\normal、example\loadfromconfig 示例, 新增 example\loadfromjson 示例
* 贡献者：https://github.com/esap
* 2017-04-19 09:00

#### Version 0.2.1
* LoopTask增加DueTime设置，表示开始任务的延迟时间（以毫秒为单位），如果<=0则不延迟
* 更新 example\normal、example\loadfromconfig 示例
* 2017-03-23 16:00

#### Version 0.2
* 新增配置文件加载方式
* 新增默认FileLogger、FmtLogger内建支持，如果未显式设置logger，自动使用FmtLogger
* 新增 example\normal、example\loadfromconfig 示例
* 2017-03-23 10:00


#### Version 0.1.1
* 新增支持RegisterHandler\GetHandler，便于配置工作
* 2017-03-19 16:00

#### Version 0.1
* Init version
* 支持cron、loop两种模式
* cron模式支持“秒 分 时 日 月 周”配置
* loop模式支持毫秒级别
* 上次任务没有停止的情况下下次任务顺延
* 2017-03-01