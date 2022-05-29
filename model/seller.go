package model
import (
	"time"
)
type Toko struct{
	Name string `json:"seller" bson:"seller"`
	Store string `json:"store_name" bson:"store_name"`
	Since time.Time `json:"since" bson:"since"`
}
