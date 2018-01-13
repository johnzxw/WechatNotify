package mongoCache

import (
	"../Tool"
	"gopkg.in/mgo.v2"
)

var mgoSession *mgo.Session

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GetSession() *mgo.Session {
	mongoUrl := Tool.GetConfig("mongoUrl")
	URL := mongoUrl["name"]
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic("mongodb connect faild! " + err.Error()) //直接终止程序运行
		}
	}
	//mgoSession.SetMode(mgo.Monotonic, true)
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	dataBaseTmp := Tool.GetConfig("mongoUrl")
	dataBase := dataBaseTmp["dataBase"]
	c := session.DB(dataBase).C(collection)
	return s(c)
}
