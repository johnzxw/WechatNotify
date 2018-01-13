package mongoCache

import (
	"../Define"
	"../Tool"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/**
 * 添加AdminUserData对象
 */
func SaveAdminUserData(p Define.AdminUserDefine) string {
	p.Id = bson.NewObjectId()
	if p.RegisterTime == "" {
		p.RegisterTime = Tool.GetDateString()
	}
	if p.LastLoginTime == "" {
		p.LastLoginTime = Tool.GetDateString()
	}
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(Define.AdminUserTableName, query)
	if err != nil {
		return "false"
	}
	return p.Id.Hex()
}

/**
 *关注用户组
 */
func SaveWechatUserGroup(p Define.WechatUserGroupStruct) bool {
	p.Id = bson.NewObjectId()
	if p.AddTime == 0 {
		p.AddTime = Tool.GetTimeStamp()
	}
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(Define.WechatUserGroupTable, query)
	if err != nil {
		return false
	}
	return true
}

/**
 *获取所有的微信用户用户组
 */
func GetAllWechatGroupData() []Define.WechatUserGroupStruct {
	var persons []Define.WechatUserGroupStruct
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&persons)
	}
	witchCollection(Define.WechatUserGroupTable, query)
	return persons

}

/**
 * 获取一条记录通过objectid
 */
func GetAdminUserDataByObjectId(id string) *Define.AdminUserDefine {
	objid := bson.ObjectIdHex(id)
	person := new(Define.AdminUserDefine)
	query := func(c *mgo.Collection) error {
		return c.FindId(objid).One(&person)
	}
	witchCollection(Define.AdminUserTableName, query)
	return person
}

/**
 *获取一条记录通过username
 */
func GetAdminUserDataByUsername(username string) *Define.AdminUserDefine {
	AdminUserQuery := bson.M{"Name": username}
	result := new(Define.AdminUserDefine)
	query := func(c *mgo.Collection) error {
		return c.Find(AdminUserQuery).One(&result)
	}
	witchCollection(Define.AdminUserTableName, query)
	return result
}

/**
 *获取一条记录通过mobile
 */
func GetAdminUserDataByMobile(mobile string) *Define.AdminUserDefine {
	AdminUserQuery := bson.M{"Mobile": mobile}
	result := new(Define.AdminUserDefine)
	query := func(c *mgo.Collection) error {
		return c.Find(AdminUserQuery).One(&result)
	}
	witchCollection(Define.AdminUserTableName, query)
	return result
}

//更新数据
func UpdateAdminUserDataForWhere(query bson.M, change bson.M) bool {
	exop := func(c *mgo.Collection) error {
		return c.Update(query, change)
	}
	err := witchCollection(Define.AdminUserTableName, exop)
	if err != nil {
		return true
	}
	return false
}

/**
 *
 */
func DeleteAdminUserDataByObjectId(id string) bool {
	deleteWhere := bson.ObjectIdHex(id)
	query := func(c *mgo.Collection) error {
		return c.Remove(deleteWhere)
	}
	err := witchCollection(Define.AdminUserTableName, query)
	if err != nil {
		return false
	}
	return true
}

/**
 *用户登录日志
 */
func SetAdminUserLoginLog(data Define.AdminUserLoginLogDefine) bool {
	query := func(c *mgo.Collection) error {
		return c.Insert(data)
	}
	err := witchCollection(Define.AdminUserLoginLogTableName, query)
	if err != nil {
		return false
	}
	return true
}

/**
 *用户属于哪个组
 */
func SaveWechatUserGroupBelong(data Define.WechatUserGroupBelongStruct) bool {
	query := func(c *mgo.Collection) error {
		return c.Insert(data)
	}
	err := witchCollection(Define.WechatUserGroupBelongTable, query)
	if err != nil {
		return false
	}
	return true
}

func DeleteAdminUserGroupBelong(groupId string) bool {
	AdminUserQuery := bson.M{"GroupId": groupId}
	query := func(c *mgo.Collection) error {
		return c.Remove(AdminUserQuery)
	}
	err := witchCollection(Define.WechatUserGroupBelongTable, query)
	if err != nil {
		return false
	}
	return true
}

func GetAllAdminUserGroupBelong(groupId string) []Define.WechatUserGroupBelongStruct {
	AdminUserQuery := bson.M{"GroupId": groupId}
	var persons []Define.WechatUserGroupBelongStruct
	query := func(c *mgo.Collection) error {
		return c.Find(AdminUserQuery).All(&persons)
	}
	err := witchCollection(Define.WechatUserGroupBelongTable, query)
	if err != nil {
		return persons
	}
	return persons
}
