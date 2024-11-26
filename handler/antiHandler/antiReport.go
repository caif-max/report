package antiHandler

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"report/common/db/mgClient"
	"report/common/log"
	"report/model/antiModel"
	"time"
)

func Start() {
	log.GetLogger().Info("开始统计anti-ban的报表")
	distinct, err := mgClient.Distinct(antiModel.DBName, antiModel.AntiUserNum, "account", &bson.M{})
	if err != nil {
		log.GetLogger().Error(err.Error())
		return
	}
	fmt.Println(distinct)
	for _, v := range distinct {
		countForUser(v.(string))
		countForAccount(v.(string))
	}
}

// 不支持员工数大于10000的企业
func countForUser(account string) {

	log.GetLogger().Info("统计今日员工数据: " + account)
	q := bson.M{}
	q["account"] = account

	op := &options.FindOptions{}
	op.SetLimit(10000)
	op.SetSkip(0)

	result := &mgClient.PageResult{}
	result.List = make([]interface{}, 1)
	result.List[0] = &antiModel.UserNumber{}
	err := mgClient.Find(antiModel.DBName, antiModel.AntiUserNum, &q, op, result)
	if err != nil {
		log.GetLogger().Error(err)
		return
	}
	user := antiModel.UserNumber{}
	for u := range result.List {
		user = result.List[u].(antiModel.UserNumber)

		//统计每天的数据
		today := time.Now().Format("2006-01-02")
		thisMonth := time.Now().Format("2006-01")
		thisYear := time.Now().Format("2006")
		q1 := bson.M{}
		q1["userId"] = user.UserId
		q1["createTime"] = bson.M{"$gte": today}
		count, err1 := mgClient.Count(account, antiModel.AntiComplaint, &q1)
		if err1 != nil {
			log.GetLogger().Error("统计员工数据出错，终止统计")
			return
		}
		r := antiModel.UserReport{}
		r.UserId = user.UserId
		r.Id = today + user.UserId
		r.Count = count
		r.Date = today

		//添加到数据库中
		_ = mgClient.InsertOne(account, antiModel.AntiUserReport, &r)

		//获取本月数据，并更新
		m := antiModel.UserReport{}
		err = mgClient.FindOne(account, antiModel.AntiUserReport, &bson.M{"_id": thisMonth + user.UserId}, &m)
		if err != nil {
			//没有本月数据，初始化一条
			m.Count = 0
			m.UserId = user.UserId
			m.Id = thisMonth + user.UserId
			m.Date = thisMonth
		} else {
			m.Count += count
		}
		_ = mgClient.InsertOne(account, antiModel.AntiUserReport, &m)

		//获取今年数据，并更新
		y := antiModel.UserReport{}
		err = mgClient.FindOne(account, antiModel.AntiUserReport, &bson.M{"_id": thisYear + user.UserId}, &y)
		if err != nil {
			y.Count = 0
			y.UserId = user.UserId
			y.Id = thisYear + user.UserId
			y.Date = thisYear
		} else {
			y.Count += count
		}
		_ = mgClient.InsertOne(account, antiModel.AntiUserReport, &y)
	}
}

func countForAccount(account string) {

	log.GetLogger().Info("开始统计今日账户数据")
	//获取所有类型，按类型统计
	op := &options.FindOptions{}
	op.SetLimit(10000)
	op.SetSkip(0)
	result := &antiModel.ComplaintStruct{}
	_ = mgClient.FindOne(antiModel.DBName, antiModel.AntiComplaintStruct, &bson.M{}, result)
	for _, v := range result.List {

		//统计每个类型今日的数据
		today := time.Now().Format("2006-01-02")
		thisMonth := time.Now().Format("2006-01")
		thisYear := time.Now().Format("2006")
		q := bson.M{}
		q["type"] = v.Id
		q["createTime"] = bson.M{"$gte": today}
		count, err1 := mgClient.Count(account, antiModel.AntiComplaint, &q)
		if err1 != nil {
			log.GetLogger().Error("统计账户数据出错，终止统计")
			return
		}
		t := antiModel.AccountReport{}
		t.Id = today + v.Id
		t.Count = count
		t.Date = today
		t.Type = v.Id
		t.Name = v.Name

		_ = mgClient.InsertOne(account, antiModel.AntiAccountReport, &t)

		//统计本月数据
		m := antiModel.AccountReport{}
		err := mgClient.FindOne(account, antiModel.AntiAccountReport, &bson.M{"_id": thisMonth + v.Id}, &m)
		if err != nil {
			//没有本月数据，初始化一条
			m.Count = 0
			m.Date = thisMonth
			m.Type = v.Id
			m.Name = v.Name
			m.Id = thisMonth + v.Id
		} else {
			m.Count += count
		}
		_ = mgClient.InsertOne(account, antiModel.AntiAccountReport, &m)

		//统计今年数据
		y := antiModel.AccountReport{}
		err = mgClient.FindOne(account, antiModel.AntiAccountReport, &bson.M{"_id": thisYear + v.Id}, &m)
		if err != nil {
			//没有本月数据，初始化一条
			y.Count = 0
			y.Date = thisYear
			y.Type = v.Id
			y.Name = v.Name
			y.Id = thisYear + v.Id
		} else {
			y.Count += count
		}
		_ = mgClient.InsertOne(account, antiModel.AntiAccountReport, &y)
	}
}
