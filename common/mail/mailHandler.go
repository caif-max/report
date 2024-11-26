package mail

import (
	"errors"
	"gopkg.in/gomail.v2"
	"strconv"
)

// 产品使用申请邮件
//type TrialAppMail struct {
//	TimeRange string   // 时间范围		 exp: 2019-05-05 17:30 ~ 2019-05-05 18:00
//	SmsNum    int      // 敏感短信数量
//	SmsAgent  []string // (短信)责任人	 exp: ["张三"]
//
//	MailTo  []string // 邮件接收人
//	Subject string   // 邮件标题
//	Url     string   // 详情地址
//}

// 联系客户信息
type TrialAppMail struct {
	Id        string   `json:"_id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`           // 姓名
	Phone     string   `json:"phone" bson:"phone"`         // 电话
	Company   string   `json:"company" bson:"company"`     // 公司
	Intention []string `json:"intention" bson:"intention"` // 意向
	LogTime   string   `json:"logTime" bson:"logTime"`     // 记录时间
	MailTo    []string `json:"mailTo" bson:"mailTo"`       // 邮件接收人						   // 邮件接收放
	Subject   string   `json:"subject" bson:"subject"`     // 邮件标题
	Title     string   `json:"title" bson:"title"`         // 邮件自定义标题
}

/**
 * 发送产品使用申请邮件
 * @return * {error} 错误响应
 */
func (rm *TrialAppMail) SendMail() error {
	if len(rm.MailTo) < 1 {
		return errors.New("MailTo is empty")
	}

	// 当前使用阿里邮箱
	// 定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "noreply@notice.dustess.com", // from
		"pass": "z76LIS9bFnIA",               // password
		"host": "smtpdm.aliyun.com",          // 发送服务器
		"port": "465",                        // 服务区端口
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	//这种方式可以添加别名，即 7Phone， 也可以直接用 m.SetHeader("From",mailConn["user"])
	//m.SetHeader("From","7Phone" + "<" + mailConn["user"] + ">")
	m.SetAddressHeader("From", mailConn["user"], rm.Subject)

	m.SetHeader("To", rm.MailTo...) // 发送给多个用户

	m.SetHeader("Subject", rm.Subject) // 设置邮件主题

	// 内容
	// 敏感词样式
	//content := "<style>.num{color: #FF4D4F}</style>"
	content := rm.getMailContent()
	m.SetBody("text/html", content) //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	// 发送邮件
	return d.DialAndSend(m)
}

/**
 * 获取邮件内容
 * @return content {string} 邮件内容
 */
//表单数据流向：将表单内容生成邮件，发到 service@dustess.com，邮件内容格式如下：
//1）title：产品试用申请（%姓名%）
//2）内容：
//a）姓名：
//b）手机号：
//c）企业名称：
//d）感兴趣内容：%列表%
func (rm *TrialAppMail) getMailContent() (content string) {
	content += "<h2>" + rm.Subject + "</h2>"
	content += "<p>姓名: " + rm.Name + "</p>"
	content += "<p>手机号: " + rm.Phone + "</p>"
	content += "<p>企业名称: " + rm.Company + "</p>"

	intention := ""
	for i, v := range rm.Intention {
		if i == 0 {
			intention += v
		} else {
			intention += ", " + v
		}
	}
	content += "<p>感兴趣内容: " + intention + "</p>"
	return
}
