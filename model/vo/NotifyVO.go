package vo

type NotifyVO struct {
	ProcessId     int      `json:"processId"`
	TaskId        int      `json:"taskId"`
	NotifyUserIds []string `json:"notifyUserIds"`
}
