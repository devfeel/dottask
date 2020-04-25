package task

import "context"

//Task上下文信息
type TaskContext struct {
	TaskID         string
	TaskData       interface{} //用于当前Task全局设置的数据项
	Message        interface{} //用于每次Task执行上下文消息传输
	IsEnd          bool        //如果设置该属性为true，则停止当次任务的后续执行，一般用在OnBegin中
	Error          error
	Header         map[string]interface{}
	TimeoutContext context.Context
	TimeoutCancel  context.CancelFunc
	doneChan       chan struct{}
}

func (c TaskContext) reset() {
	c.TaskID = ""
	c.TaskData = nil
	c.Message = nil
	c.IsEnd = false
	c.Error = nil
	c.Header = nil
	c.TimeoutContext = nil
	c.TimeoutCancel = nil
	c.doneChan = nil
}
