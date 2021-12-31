package message

import (
	"gopkg.in/gomail.v2"
)

// emailDiale .
type emailDiale struct {
	*gomail.Dialer
}

// NewEmailDialer .
func NewEmailDialer(smtp string, port int, username string, password string) *emailDiale {
	return &emailDiale{Dialer: gomail.NewDialer(smtp, port, username, password)}
}

func (d *emailDiale) Send(em *EmailMessage) error {
	gm := gomail.NewMessage()
	gm.SetAddressHeader("From", d.Username, "")
	gm.SetHeader("To", em.To...)
	if len(em.Cc) != 0 {
		gm.SetHeader("Cc", em.Cc...)
	}
	gm.SetHeader("Subject", em.Subject)
	gm.SetBody(em.ContentType.String(), em.Content)
	if em.Attach != "" {
		gm.Attach(em.Attach)
	}
	return d.DialAndSend(gm)
}

type EmailMessageType string

const (
	EmailMessageTypePlain EmailMessageType = "text/plain"
	EmailMessageTypeHTML  EmailMessageType = "text/html"
)

// String .
func (mt *EmailMessageType) String() string {
	return string(*mt)
}

// EmailMessage .
type EmailMessage struct {
	To          []string         // 收件人
	Cc          []string         // 抄送人
	Subject     string           // 标题
	ContentType EmailMessageType // 内容类型 text/plain text/html
	Content     string           // 内容
	Attach      string           // 附件
}
