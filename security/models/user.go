package models

import "github.com/goincremental/dal"

// OAuthID allows a user to be matched to different OAuth provider accounts
// (i.e. google, facebook, twitter, linked in) using the id from that provider.
type OAuthID struct {
	Provider string `bson:"provider"`
	ID       string `bson:"providerId"`
}

//User is a struct that represents a user of the application
type User struct {
	ID        dal.ObjectID `bson:"_id"`
	SystemID  dal.ObjectID `bson:"systemId"`
	FirstName string       `bson:"firstName"`
	LastName  string       `bson:"lastName"`
	Email     string       `bson:"email"`
	Password  string       `bson:"password"`
	APISecret string       `bson:"apiSecret"`
	OAuthID   []OAuthID    `bson:"userIds"`
	Roles     []string     `bson:"roles"`
}

const userCollection string = "users"

//GetUserByOAuthID allows a single user to be found by provider name and id
func GetUserByOAuthID(db dal.Database, systemID *dal.ObjectID,
	id string) (result User, err error) {
	users := db.C(userCollection)

	err = users.Find(dal.BSON{
		"systemId":           systemID,
		"userIds.providerId": id},
	).One(&result)

	return
}

//GetUserByID allows retrieval of a user by its id and systemID
func GetUserByID(db dal.Database, systemID *dal.ObjectID,
	id *dal.ObjectID) (result User, err error) {
	users := db.C(userCollection)

	err = users.Find(dal.BSON{
		"systemId": systemID,
		"_id":      id},
	).One(&result)

	return
}

//Save persists the user to the database through the dal
func (u *User) Save(db dal.Database) (err error) {
	col := db.C(userCollection)
	if !u.ID.Valid() {
		u.ID = dal.NewObjectId()
	}
	_, err = col.UpsertID(u.ID, u)

	return err
}
