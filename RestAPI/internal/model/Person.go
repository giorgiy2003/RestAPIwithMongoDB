package Model

type Person struct {
	Id        interface{} `json:"id" bson:"_id"`
	Email     string      `json:"email" bson:"email"`
	Phone     string      `json:"phone" bson:"phone"`
	FirstName string      `json:"firstName" bson:"firstName"`
	LastName  string      `json:"lastName" bson:"lastName"`
}
