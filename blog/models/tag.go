package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/JonSnow47/Graduation-Project/blog/consts"
	"github.com/JonSnow47/Graduation-Project/blog/mongo"
)

type tagServiceProvider struct{}

var TagService *tagServiceProvider

/*func init() (session *mgo.Session) {
	session, err := mgo.Dial(conf.MongoURL + "/" + consts.Database)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	session.DB(consts.Database).C(consts.CollectionTag).EnsureIndex(mgo.Index{
		Key:        []string{"tag"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	return
}*/

func CollectionTag() mongo.Mongodb {
	m := mongo.ConnectMongo(consts.CollectionTag)
	m.C.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})
	return m
}

// Tag represents the tag of article.
type Tag struct {
	Id    bson.ObjectId `bson:"_id,omitempty"`
	Name  string        `bson:"name"`
	Count int           `bson:"count"`
	State bool          `bson:"state"`
}

// New insert a new article.
func (*tagServiceProvider) New(tag string) (string, error) {
	m := CollectionTag()
	defer m.S.Close()

	t := &Tag{
		Id:   bson.NewObjectId(),
		Name: tag,
	}
	err := m.C.Insert(t)
	return t.Id.Hex(), err
}

// Delete modify tag's status.
func (*tagServiceProvider) Delete(id string) error {
	m := CollectionTag()
	defer m.S.Close()

	err := m.C.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"state": false}})
	return err
}

// Enable modify tag's status.
func (*tagServiceProvider) Enable(id string) error {
	m := CollectionTag()
	defer m.S.Close()

	err := m.C.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"state": true}})
	return err
}

// Get one tag's info.
func (*tagServiceProvider) Get(id string) (t Tag, err error) {
	m := CollectionTag()
	defer m.S.Close()

	err = m.C.FindId(bson.ObjectIdHex(id)).One(&t)
	return
}

// All get all tags list.
func (*tagServiceProvider) All(page int) (tags []Tag, err error) {
	m := CollectionTag()
	defer m.S.Close()

	err = m.C.Find(nil).Limit(20).Skip(20 * page).All(&tags)
	return tags, err
}

func (*tagServiceProvider) Count(tagsid []bson.ObjectId) error {
	m := CollectionTag()
	defer m.S.Close()

	for i, _ := range tagsid {
		err := m.C.UpdateId(tagsid[i], bson.M{"$inc": bson.M{"count": +1}})
		if err != nil {
			return err
		}
	}
	return nil
}