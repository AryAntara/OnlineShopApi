package model

type Cart struct {
	StoreId string `bson:"store" json=:"store"`
	Product string `bson:"name_product" json:"name_product"`
	Count int `bson:"count" json:"count"`
	Id string `bson:"productid" json:"productid"`
}
type User struct {
	Name string `bson:"name" json:"name"`
	Password string `bson:"password" json:"password"`
	Cart []Cart `bson:"cart" json:"cart"`
	RefreshToken string `bson:"refresh_token" json:"refresh_token"`
}
