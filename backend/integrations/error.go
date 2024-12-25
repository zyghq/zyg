package integrations

type integrationErr string

func (err integrationErr) Error() string {
	return string(err)
}

const (
	ErrPostmarkSendMail = integrationErr("postmark send mail error")
)
