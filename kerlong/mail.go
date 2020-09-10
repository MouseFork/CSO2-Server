package kerlong

import (
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	SenderMail    string
	SenderCode    string
	SenderSMTP    string
	TargetMail    string
	MailSubHeader string
	MailHeader    string
	Content       string
}

func SendEmailTO(mail *EmailData) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", mail.SenderMail, mail.MailSubHeader) // 发件人
	m.SetHeader("To",                                               // 收件人
		m.FormatAddress(mail.TargetMail, "111"),
	)
	m.SetHeader("Subject", mail.MailHeader) // 主题
	body := mail.Content
	m.SetBody("text/html", body) // 正文

	d := gomail.NewPlainDialer(mail.SenderSMTP, 465, mail.SenderMail, mail.SenderCode) // 发送邮件服务器、端口、发件人账号、发件人密码

	err := d.DialAndSend(m)

	return err
}
