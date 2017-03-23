devfeel/task

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