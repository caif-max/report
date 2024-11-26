package mail

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"report/common/log"
)

type MailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var config = MailConfig{
	Host:     "smtpdm.aliyun.com",
	Port:     465,
	User:     "noreply@notice.dustess.com",
	Password: "z76LIS9bFnIA",
}

type SendMailAgentApplyData struct {
	Name             string `json:"name" bson:"name"`                         // 代理商名称
	RegisterManName  string `json:"registerManName" bson:"registerManName"`   //申请人
	RegisterManPhone string `json:"registerManPhone" bson:"registerManPhone"` //申请人电话
	Area             string `json:"area" bson:"area"`                         //区域
	Location         string `json:"location" bson:"location"`                 //公司位置
	Situation        string `json:"situation" bson:"situation"`               //概况
}

func SendMailAgentApply(data SendMailAgentApplyData) error {
	log.GetLogger().Info("Send <AgentApply> Email ")
	toAddress := "fangf@dustess.com"
	toName := "樊国富"
	return sendMail(toAddress, toName, "代理商申请提醒", "common/mail/resources/agentApply.html", data)
}

func sendMail(toAddress string, toName string, subject string, templatePath string, data interface{}) error {
	//render html
	var htmlBuffer bytes.Buffer
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	if err := t.Execute(&htmlBuffer, data); err != nil {
		return err
	}
	html := htmlBuffer.String()
	log.GetLogger().Debug(html)

	m := gomail.NewMessage()
	m.SetHeader("From", config.User)
	m.SetAddressHeader("To", toAddress, toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	d := gomail.NewDialer(config.Host, config.Port, config.User, config.Password)

	return d.DialAndSend(m)
}
