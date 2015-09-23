package mongodb

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	mgoSession      *mgo.Session
	readOnlySession *mgo.Session
	//MgoLogSession   *mgo.Session
)

// Exec execute a function with given collectionName. It will close the session after the fuction returning
func Exec(collectionName string, f func(*mgo.Collection) error) (recordExist bool) {
	return SessionExec(mgoSession, collectionName, f)
}

func Read(collectionName string, f func(*mgo.Collection) error) (recordExist bool) {
	return SessionExec(readOnlySession, collectionName, f)
}

func SessionExec(mSession *mgo.Session, collectionName string, f func(*mgo.Collection) error) (recordExist bool) {
	session := mSession.Clone()
	defer session.Close()
	err := f(session.DB("").C(collectionName))
	if err != nil {
		if err == mgo.ErrNotFound {
			return false
		}
	}
	return true

}

func ExecBulk(mSession *mgo.Session, collectionName string, f func(*mgo.Bulk)) {
	SessionExec(mgoSession, collectionName, func(c *mgo.Collection) error {
		bulk := c.Bulk()
		f(bulk)
		_, err := bulk.Run()
		return err
	})

}

func GetMgoSession() *mgo.Session {
	return mgoSession
}

func Cmd(f func(*mgo.Database) error) {
	session := mgoSession.Clone()
	defer session.Close()
	err := f(session.DB(""))
	if err != nil {
	}

}

func connect(url string) *mgo.Session {
	session, err := mgo.Dial(url)
	if err != nil {
	}
	return session

}

func AppendAttrs(data interface{}, attrs bson.M) bson.RawD {
	rawID, err := ToRaw(attrs)
	if err != nil {
		panic(err)
	}

	rawData, err := ToRaw(data)
	if err != nil {
		panic(err)
	}

	return append(rawID, rawData...)

}

func ToRaw(data interface{}) (bson.RawD, error) {

	rawResult := make(bson.RawD, 0)
	out, err := bson.Marshal(data)

	if err != nil {
		return nil, err
	}
	err = bson.Unmarshal(out, &rawResult)
	if err != nil {
		return nil, err
	}

	return rawResult, nil

}

// Init initialize mongodb connection
func Init(url string) {
	mgoSession = connect(url)
	readOnlySession = mgoSession.Clone()
	readOnlySession.SetMode(mgo.Eventual, false)

}

// Helper function for selecting specified fields
// c.Find(xxx).Select(sel("fieldA", "fieldB"))
func Selc(q ...string) (r bson.M) {
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func Update(collectionName string, query, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		return c.Update(query, bson.M{"$set": updateData})
	})
}

func UpdateAll(collectionName string, query, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		_, err := c.UpdateAll(query, bson.M{"$set": updateData})
		return err
	})
}

func UpdateUnset(collectionName string, query, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		return c.Update(query, bson.M{"$unset": updateData})
	})
}

func FindByID(collectionName string, id interface{}, result interface{}) bool {
	return Read(collectionName, func(c *mgo.Collection) error {
		return c.FindId(id).One(result)
	})
}

func UpdateIDUnset(collectionName string, id interface{}, updateData interface{}) bool {
	return UpdateUnset(collectionName, bson.M{"_id": id}, updateData)
}

func UpdateID(collectionName string, id interface{}, updateData interface{}) bool {
	return Update(collectionName, bson.M{"_id": id}, updateData)
}

func UpsertID(collectionName string, id interface{}, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		_, err := c.UpsertId(id, bson.M{"$set": updateData})
		return err
	})
}

func Delta(collectionName string, query, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		return c.Update(query, bson.M{"$inc": updateData})
	})
}

func DeltaID(collectionName string, id interface{}, updateData interface{}) bool {
	return Exec(collectionName, func(c *mgo.Collection) error {
		return c.Update(bson.M{"_id": id}, bson.M{"$inc": updateData})
	})
}
