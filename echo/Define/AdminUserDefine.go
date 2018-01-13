package Define

import "gopkg.in/mgo.v2/bson"

var AdminUserTableName = "admin_user"

var WechatUserGroupTable = "wechat_user_group"

var WechatUserGroupBelongTable = "wechat_user_group_belong"

//Name unique index
type AdminUserDefine struct {
	Id            bson.ObjectId `bson:"_id"`
	Name          string        `bson:"Name"`
	Pass          string        `bson:"Pass"`
	Salt          string        `bson:"Salt"`
	Mobile        string        `bson:"Mobile"`
	RegisterTime  string        `bson:"RegisterTime"`
	LastLoginTime string        `bson:"LastLoginTime"`
	RegisterIp    string        `bson:"RegisterIp"`
}

/**
 *用户组结构
 * name unique
 */
type WechatUserGroupStruct struct {
	Id      bson.ObjectId `bson:"_id"`
	Name    string        `bson:"Name"`
	AddTime int64         `bson:"AddTime"`
	ForName string        `bson:"ForName"`
}

/**
 *用户属于哪个组
 */
type WechatUserGroupBelongStruct struct {
	Id      bson.ObjectId `bson:"_id"`
	Uid     string        `bson:"Uid"`
	GroupId string        `bson:"GroupId"`
	AddTime int64         `bson:"AddTime"`
}
