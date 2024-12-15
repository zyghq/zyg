package services

type serviceErr string

func (err serviceErr) Error() string {
	return string(err)
}

const (
	ErrAccount                   = serviceErr("account error")
	ErrAccountNotFound           = serviceErr("account not found")
	ErrPat                       = serviceErr("pat error")
	ErrPatNotFound               = serviceErr("pat not found")
	ErrWorkspace                 = serviceErr("workspace error")
	ErrWorkspaceNotFound         = serviceErr("workspace not found")
	ErrMember                    = serviceErr("member error")
	ErrMemberNotFound            = serviceErr("member not found")
	ErrLabel                     = serviceErr("label error")
	ErrLabelNotFound             = serviceErr("label not found")
	ErrThreadChat                = serviceErr("thread chat error")
	ErrThread                    = serviceErr("thread error")
	ErrThreadNotFound            = serviceErr("thread not found")
	ErrThreadMetrics             = serviceErr("thread chat metrics error")
	ErrThreadMessage             = serviceErr("thread message error")
	ErrThreadLabel               = serviceErr("thread label error")
	ErrCustomer                  = serviceErr("customer error")
	ErrCustomerNotFound          = serviceErr("customer not found")
	ErrSecretKeyNotFound         = serviceErr("secret key not found")
	ErrSecretKey                 = serviceErr("secret key error")
	ErrWidget                    = serviceErr("widget error")
	ErrWidgetNotFound            = serviceErr("widget not found")
	ErrWidgetSession             = serviceErr("widget session error")
	ErrWidgetSessionInvalid      = serviceErr("widget session invalid")
	ErrClaimedMail               = serviceErr("claimed mail error")
	ErrClaimedMailNotFound       = serviceErr("claimed mail not found")
	ErrClaimedMailExpired        = serviceErr("claimed mail expired")
	ErrCustomerEvent             = serviceErr("customer event error")
	ErrPostmarkInbound           = serviceErr("postmark inbound error")
	ErrMessageAttachment         = serviceErr("message attachment error")
	ErrMessageAttachmentNotFound = serviceErr("message attachment not found")
	ErrPostmarkSettingNotFound   = serviceErr("postmark setting not found")
	ErrPostmarkSetting           = serviceErr("postmark setting error")
)
