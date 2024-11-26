package antiModel

const (
	DBName              = "anti-ban"
	AntiUserNum         = "user_num"
	AntiComplaintStruct = "complaint_struct"

	//以下表都在account库里面
	AntiComplaint     = "complaint"
	AntiUserReport    = "user_report"
	AntiAccountReport = "account_report"
)

type UserNumber struct {
	Id      string `json:"id" bson:"_id"`
	Account string `json:"account" bson:"account"`
	UserId  string `json:"userId" bson:"userId"`
}

type Complaint struct {
	Id         string   `json:"id" bson:"_id"`
	Type       string   `json:"type"`
	SecondType string   `json:"secondType"`
	Pic        []string `json:"pic"`
	Note       string   `json:"note"`
	Phone      string   `json:"phone"`
	UserId     string   `json:"userId" bson:"userId"`
	CreateTime string   `json:"createTime" bson:"createTime"`
}

type ComplaintType struct {
	Id         string                `json:"id" bson:"_id"`
	Name       string                `json:"name"`
	SecondType []ComplaintSecondType `json:"secondType" bson:"secondType"`
}

type ComplaintSecondType struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name"`
}

type ComplaintStruct struct {
	Id   string          `json:"id" bson:"_id"`
	List []ComplaintType `json:"list"`
}

type UserReport struct {
	Id     string `json:"id" bson:"_id"`
	Date   string `json:"date" bson:"date"`
	UserId string `json:"userId"`
	Count  int64  `json:"count"`
}

type AccountReport struct {
	Id    string `json:"id" bson:"_id"`
	Date  string `json:"date"`
	Count int64  `json:"count"`
	Type  string `json:"type"`
	Name  string `json:"name"`
}
