package services

type serviceErr string

func (err serviceErr) Error() string {
	return string(err)
}

const (
	ErrAccount         = serviceErr("account error")
	ErrAccountNotFound = serviceErr("account not found")

	ErrPat         = serviceErr("pat error")
	ErrPatNotFound = serviceErr("pat not found")

	ErrWorkspace         = serviceErr("workspace error")
	ErrWorkspaceNotFound = serviceErr("workspace not found")

	ErrMember         = serviceErr("member error")
	ErrMemberNotFound = serviceErr("member not found")

	ErrLabel         = serviceErr("label error")
	ErrLabelNotFound = serviceErr("label not found")

	ErrThread         = serviceErr("thread error")
	ErrThreadNotFound = serviceErr("thread not found")

	ErrThreadMetrics = serviceErr("thread chat metrics error")

	ErrThreadActivity = serviceErr("thread activity error")

	ErrCustomer         = serviceErr("customer error")
	ErrCustomerNotFound = serviceErr("customer not found")

	ErrSecretKeyNotFound = serviceErr("secret key not found")
	ErrSecretKey         = serviceErr("secret key error")

	ErrCustomerEvent = serviceErr("customer event error")

	ErrMessageAttachment         = serviceErr("message attachment error")
	ErrMessageAttachmentNotFound = serviceErr("message attachment not found")

	ErrPostmarkSettingNotFound = serviceErr("postmark setting not found")
	ErrPostmarkSetting         = serviceErr("postmark setting error")

	ErrPostmarkLog         = serviceErr("postmark log error")
	ErrPostmarkLogNotFound = serviceErr("postmark log not found")
	ErrPostmarkInbound     = serviceErr("postmark inbound error")
	ErrPostmarkOutbound    = serviceErr("postmark outbound error")
)
