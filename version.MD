devfeel/task

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