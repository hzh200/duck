package constance

type TaskStatus int64

const (
	_ TaskStatus = iota
    Waiting
    Running
    Paused
    Stopped
    Successed
    Failed
)
