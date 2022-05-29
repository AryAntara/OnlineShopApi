package model

type Product struct {
  Name string `bson:"name_product" json:"name_product"`
	Stock int `bson:"stock_product" json:"stock_product"`
	//Category []string `bson:"category" json:"category"`
	StoreId string `bson:"store" json:"store"`
	Price string `bson:"price" json:"price"`
	Priceint int `bson:"price_int" json:"price_int"`
	Id string `json:"_id"`
}
