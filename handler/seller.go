package control
import(
	"context"
	"fmt"
	model "RESTAPI/model"
	"RESTAPI/database"
	"github.com/gofiber/fiber/v2"
	//bcrypt "golang.org/x/crypto/bcrypt"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/golang-jwt/jwt/v4"
	"strings"
	"RESTAPI/config"
	"time"
	"errors"
)
var SellerDb = db.Collection(config.DbConn,config.Db,config.StoreColl)
func MakeStore(c *fiber.Ctx) error {
	//using cookie
	idd := c.Query("id")
	store := c.Query("store")
	if store == "" || idd == "" {
		c.Send([]byte("error query"))
		return errors.New("empty query")
	}
	ids := strings.ReplaceAll(idd, `"`,"")
	id,err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		fmt.Println(":: err >",err)
		c.Send([]byte(":: login first to get user id"))
		return err
	}
	var result bson.M
	err = UserDb.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(&result)
	if err != nil {
		c.SendString(":: user not found")
		return err
	}
	var results bson.M
	err = SellerDb.FindOne(context.TODO(),bson.M{"seller":result["name"]}).Decode(&results)
	if err == nil {
		c.SendString("user has been creating store")
		return errors.New("user has created store")
	}
	fmt.Println("::",err)
	toko := model.Toko{
		Name: result["name"].(string),
		Store: store,
		Since: time.Now(),
	}
	SellerDb.InsertOne(context.TODO(),toko)
	c.Send([]byte("store has added"))
	return nil
}
func StoreName(c *fiber.Ctx) error {
	UserId := c.Query("id")
	id,err := primitive.ObjectIDFromHex(UserId)
	if err != nil {
		fmt.Println(":: err >",err)
		c.Send([]byte(":: login first to get user id"))
		return err
	}
	var result bson.M
	err = UserDb.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(&result)
	if err != nil {
		fmt.Println(":: user not found")
		c.SendString(":: user not found")
		return err
	}
	var results bson.M
	err = SellerDb.FindOne(context.TODO(),bson.M{"seller":result["name"]}).Decode(&results)
	if err != nil {
		c.SendString("user has not created a store")
		return err
	}
	data,_ := json.Marshal(results["_id"])
	c.SendString(strings.ReplaceAll(string(data),`"`,"")+" "+results["store_name"].(string))
	return nil
}
func AddProduct(c *fiber.Ctx) error {
	StoreID := c.Query("store_id")
	Product := c.Query("product")
	if StoreID == "" || Product == ""{
		c.Send([]byte("empty string"))
		return errors.New("empty query")
	}
	//validate store
	DbResult := bson.M{}
	id,_ := primitive.ObjectIDFromHex(StoreID)
	if err := SellerDb.FindOne(context.TODO(),bson.M{"_id":id}).Decode(&DbResult);err != nil {
		fmt.Println(err)
		c.Send([]byte("error"))
		return nil
	}
	resultProduct := model.Product{}
	if err := json.Unmarshal([]byte(Product),&resultProduct);err != nil {
		fmt.Println(err)
		c.Send([]byte("error"))
		return err
	}
	if res,err := Pr.InsertOne(context.TODO(),resultProduct);err != nil {
		fmt.Println(err)
		c.Send([]byte("error"))
		return err
	}else {
		fmt.Println(res)
	}
	return nil
}
