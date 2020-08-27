package myemail

import (
	"net/smtp"
	"strings"

	"github.com/xxjwxc/public/mylog"

	"github.com/xxjwxc/public/message"
)

// myEmail ...
type myEmail struct {
	user     string
	password string
	host     string
	title    string
}

// New 新建一个
func New(user, password, host, title string) *myEmail {
	return &myEmail{
		user:     user,
		password: password,
		host:     host,
		title:    title,
	}
}

// SendMail 发送邮件
/*
to: 目的人
subject: 标题
body: 邮件内容
*/
func (e *myEmail) SendMail(to []string, subject, body string) (state bool, code int) {
	// now := time.Now()
	// year, mon, day := now.Date()
	// hour, min, sec := now.Clock()
	// dayString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, mon, day, hour, min, sec)
	// body = fmt.Sprintf(MailBodyCn, verification, dayString)

	// index := strings.LastIndex(to, "@")
	// if checkGbkMail(to[index+1 : len(to)]) {
	// 	subject = mahonia.NewEncoder("gbk").ConvertString(fmt.Sprintf(MailSubjectCn, verification))
	// } else {
	// 	subject = fmt.Sprintf(MailSubjectCn, verification)
	// }

	err := SendToMail(e.user, e.password, e.host, e.title, subject, body, "html", to)
	if err != nil {
		mylog.Errorf("Send mail error:%s", err)
		return false, int(message.MailSendFaild)
	}

	return true, int(message.NormalMessageID)
}

// SendToMail 发送邮件
/*
 *    user : example@example.com login smtp server user
 *    password: xxxxx login smtp server password
 *    host: smtp.example.com:port   smtp.163.com:25
 *    to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *    subject:The subject of mail
 *    body: The content of mail
 *    mailtyoe: mail type html or text
 */
func SendToMail(user, password, host, title, subject, body, mailtype string, to []string) error {
	//d := mahonia.NewDecoder("UTF-8")
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("xxj", user, password, hp[0])
	contentType := "Content-Type:text/plain;charset=UTF-8"
	if mailtype == "html" {
		contentType = "Content-Type:text/html;charset=UTF-8"
	}

	msg := []byte("To:" + strings.Join(to, ";") + "\r\nFrom:" + title + "<" + user + ">\r\nSubject:" + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	return smtp.SendMail(host, auth, user, to, msg)
}
