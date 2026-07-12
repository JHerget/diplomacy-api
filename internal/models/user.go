package models

type User struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	Salt        []byte `json:"salt" bson:"salt"`
	CreatedDate int64  `json:"createdDate" bson:"createdDate"`
	IsDeleted   bool   `json:"isDeleted" bson:"isDeleted"`
}
