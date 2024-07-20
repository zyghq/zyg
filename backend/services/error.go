package services

type serviceErr string

func (err serviceErr) Error() string {
	return string(err)
}

const (
	ErrAccount               = serviceErr("account error")
	ErrAccountNotFound       = serviceErr("account not found")
	ErrPat                   = serviceErr("pat error")
	ErrPatNotFound           = serviceErr("pat not found")
	ErrWorkspace             = serviceErr("workspace error")
	ErrWorkspaceNotFound     = serviceErr("workspace not found")
	ErrMember                = serviceErr("member error")
	ErrMemberNotFound        = serviceErr("member not found")
	ErrLabel                 = serviceErr("label error")
	ErrLabelNotFound         = serviceErr("label not found")
	ErrThreadChat            = serviceErr("thread chat error")
	ErrThreadChatNotFound    = serviceErr("thread chat not found")
	ErrThreadChatAssign      = serviceErr("thread chat assignee error")
	ErrThreadChatReplied     = serviceErr("thread chat reply error")
	ErrThreadChatMetrics     = serviceErr("thread chat metrics error")
	ErrThChatMessage         = serviceErr("thread chat message error")
	ErrThChatMessageNotFound = serviceErr("thread chat message not found")
	ErrThChatLabelNotFound   = serviceErr("thread chat label not found")
	ErrThChatLabel           = serviceErr("thread chat label error")
	ErrCustomer              = serviceErr("customer error")
	ErrCustomerNotFound      = serviceErr("customer not found")
	ErrSecretKeyNotFound     = serviceErr("secret key not found")
	ErrSecretKey             = serviceErr("secret key error")
	ErrWidget                = serviceErr("widget error")
	ErrWidgetNotFound        = serviceErr("widget not found")
)
