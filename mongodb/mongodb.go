package mongodb

import (
	"gopkg.in/mgo.v2"
	"sync"
)

type Mongodb struct {
	Connect string //连接字符串
}

var (
	m          *Mongodb
	once       sync.Once
	mgoSession *mgo.Session
)

/**
 * 返回单例实例
 * @method New
 */
func New(connect string) *Mongodb {
	once.Do(func() { //只执行一次
		m = &Mongodb{Connect: connect}
	})
	return m
}

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(m.Connect)
		if err != nil {
			panic(err) //直接终止程序运行
		}
		mgoSession.SetMode(mgo.Monotonic, true) //设置一致性模式
	}
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
// func (m *Mongodb) SwitchC(collName string, f func(*mgo.Collection)) {
// 	session := getSession()
// 	defer session.Close()
// 	c := session.DB("").C(collName)
// 	f(c)
// }
func (m *Mongodb) SwitchC(collName string) (*mgo.Session, *mgo.Collection) {
	session := getSession()
	c := session.DB("").C(collName)
	return session, c
}

// /**
//  * [SearchPerson description]
//  * @param {[type]} collectionName string [集合名称]
//  * @param {[type]} query          bson.M [查询对象，不筛选为nil]
//  * @param {[type]} sort           string [排序字符串]
//  * @param {[type]} fields         bson.M [返回字段，全部为nil]
//  * @param {[type]} skip           int    [description]
//  * @param {[type]} limit          int
//  * @return   (results      []interface{}, err error [description]
//  */
// func Search(collectionName string, query bson.M, sort string, fields bson.M, skip int, limit int) (results []interface{}, err error) {
// 	exop := func(c *mgo.Collection) error {
// 		return c.Find(query).Sort(sort).Select(fields).Skip(skip).Limit(limit).All(&results)
// 	}
// 	err = SwitchCollection(collectionName, exop)
// 	return
// }
