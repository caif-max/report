package mgClient

import (
	"context"
	"github.com/mohae/deepcopy"
	"go.mongodb.org/mongo-driver/bson"
	"report/common"
	"report/common/config"
	"report/common/log"
	"report/model/antiModel"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var once sync.Once
var client *mongo.Client

type PageResult struct {
	Count int64 `json:"count"`
	List  []any `json:"list"`
}

// Init 根据name初始化client
func Init() error {
	var err error
	once.Do(func() {
		client, err = connect()
	})
	return err
}

// client 初始化连接
func connect() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var url = "mongodb://" + config.GetConf("mongo.ip") + ":" + config.GetConf("mongo.port")

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(url),
		options.Client().SetAppName(antiModel.DBName),
		options.Client().SetConnectTimeout(10*time.Second),
		options.Client().SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    "admin",
			Username:      config.GetConf("mongo.user"),
			Password:      config.GetConf("mongo.pwd"),
			PasswordSet:   true,
		}),
	)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

func DefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), common.MongoQueryDefaultTimeBill*time.Second)
}

func FindOne(dbName string, collectionName string, query *bson.M, result interface{}) error {
	ctx, cancel := DefaultContext()
	defer cancel()
	c := client.Database(dbName).Collection(collectionName)
	return c.FindOne(ctx, query).Decode(result)
}

// Find 传result参数的时候需要注意，里面有一个任意类型的list，使用vo里面的结构体作为类型
// 调用次函数前要将result中的list赋值一个结构体对象，如:
//
//	 result := mgClient.PageResult{}
//		result.List = make([]interface{}, 1)
//		result.List[0] = &model.Account{}
//
// 这里会将查询到的bson结果转化为对应的结构体变量，这样就不需要自己去处理bson
// http模块会将结构体转化为json
// 因此我们只要找结构体中声明清楚json和bson格式即可
func Find(dbName string, collectionName string, query *bson.M,
	options *options.FindOptions, result *PageResult) error {
	ctx, cancel := DefaultContext()
	defer cancel()

	//获取list中的结构体
	origin := result.List[0]
	result.List = result.List[0:0]

	c := client.Database(dbName).Collection(collectionName)
	cursor, err := c.Find(ctx, query, options)
	if err != nil {
		log.GetLogger().Fatal(err)
		return err
	}

	for cursor.Next(ctx) {
		e := deepcopy.Copy(origin)
		if err := cursor.Decode(e); err != nil {
			return err
		}
		result.List = append(result.List, e)
	}

	if result.Count, err = c.CountDocuments(ctx, query); err != nil {
		return err
	}

	if err = cursor.Close(ctx); err != nil {
		return err
	}

	return nil
}

func Count(dbName string, collectionName string, query *bson.M) (int64, error) {
	ctx, cancel := DefaultContext()
	defer cancel()
	c := client.Database(dbName).Collection(collectionName)
	return c.CountDocuments(ctx, query)
}

func Distinct(dbName string, collectionName string, filedName string, query *bson.M) ([]interface{}, error) {
	ctx, cancel := DefaultContext()
	defer cancel()
	c := client.Database(dbName).Collection(collectionName)
	return c.Distinct(ctx, filedName, query)
}

func InsertOne(dbName string, collectionName string, document interface{}) error {
	ctx, cancel := DefaultContext()
	defer cancel()
	c := client.Database(dbName).Collection(collectionName)
	_, err := c.InsertOne(ctx, document)
	if err != nil {
		return err
	}
	return nil
}

func UpdateById(dbName string, collectionName string, id string, update interface{}) error {
	ctx, cancel := DefaultContext()
	defer cancel()
	c := client.Database(dbName).Collection(collectionName)
	_, err := c.UpdateOne(ctx, &bson.M{"_id": id}, &bson.M{"$set": update})
	if err != nil {
		return err
	}
	return nil
}
