package Define

import "gopkg.in/mgo.v2/bson"

var AdminUserLoginLogTableName = "admin_user_login_log"

type AdminUserLoginLogDefine struct {
	Id             bson.ObjectId `bson:"_id"`
	Name           string        `bson:"Name"`
	LoginIp        string        `bson:LoginIp`
	LoginTime      string        `bson:"LoginTime"`
	LoginTimeStamp int64         `bson:"LoginTimeStamp"`
}
