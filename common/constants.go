package common

import (
	"time"
)

const (
	MINUTE    = 60
	HOUR      = MINUTE * 60
	HALF_HOUR = MINUTE * 30
	DAY       = HOUR * 24

	// REDIS 前缀
	SessionPrefix = "anti-ban_session."
)

const (
	MongoQueryDefaultTimeBill time.Duration = 5
	// MongoQueryLongTimeOut 20s超时时间
	MongoQueryLongTimeBill time.Duration = 20
)
