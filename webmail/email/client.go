package email

import (
	"fmt"

	"gopkg.in/gomail.v2"

	"github.com/charlesbases/library/logger"
	"github.com/charlesbases/library/webmail"
)

type ContentType string

const (
	ContentTypeHTML  ContentType = "text/html"
	ContentTypePlain ContentType = "text/plain"
)

var _ webmail.Message = (*Message)(nil)

// Message .
type Message struct {
	To          []string    // 收件人
	Cc          []string    // 抄送人
	Subject     string      // 标题
	Content     string      // 内容
	ContentType ContentType // 内容类型
	Attach      string      // 附件
}

// client .
type client struct {
	*gomail.Dialer
}

func (c *client) Send(v ...webmail.Message) error {
	var mess = make([]*gomail.Message, 0, len(v))
	for _, data := range v {
		if x, ok := data.(*Message); ok {
			mess = append(mess, gomail.NewMessage(func(m *gomail.Message) {
				m.SetAddressHeader("From", c.Username, "")
				m.SetHeader("To", x.To...)
				m.SetHeader("Subject", x.Subject)
				m.SetBody(string(x.ContentType), x.Content)

				if len(x.Cc) != 0 {
					m.SetHeader("Cc", x.Cc...)
				}
				if len(x.Attach) != 0 {
					m.Attach(x.Attach)
				}
			}))
		} else {
			err := fmt.Errorf("unsupported email message type of %T.", x)
			logger.Errorf(`[webmail.email] %s`, err.Error())
			return err
		}
	}

	if err := c.DialAndSend(mess...); err != nil {
		logger.Errorf(`[webmail.email] send failed. %s`, err.Error())
		return err
	}
	return nil
}

// NewClient .
func NewClient(smtp string, port int, username string, password string) (webmail.Client, error) {
	d := gomail.NewDialer(smtp, port, username, password)
	if closer, err := d.Dial(); err != nil {
		return nil, err
	} else {
		closer.Close()
	}

	return &client{Dialer: d}, nil
}
