package control
import(
	"context"
	"fmt"
	model "RESTAPI/model"
	"RESTAPI/database"
	"github.com/gofiber/fiber/v2"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"RESTAPI/config"
	"errors"
)
var Pr = db.Collection(config.DbConn,config.Db,config.ProdColl)
func Listfilter(c *fiber.Ctx) error {
	tos := c.Query("to")
	froms := c.Query("from")
	filter := c.Query("filter")
	if tos == "" || froms == "" || filter == "" {
		c.SendString("query empty")
		return errors.New("query empty")
	}
	to,_ := strconv.Atoi(tos)
	from,_ := strconv.Atoi(froms)
	type res struct {
		FromOld int `json="from_old"`
		ToOld int `json="to_old"`
		InDb int `json="db_length"`
		Data []primitive.M `json="data"`
	}
	ctx := context.TODO()
	cursor,err := Pr.Find(ctx,bson.M{"name_product":filter})
	if err != nil {
		return err
	}
	result := []bson.M{}
	if err := cursor.All(ctx,&result);err != nil {
		return err
	}
	if to > len(result) {
		to = len(result)
	}
	if to < from {
		return errors.New("equavalent error")
	}
	resp := res{
		FromOld: from,
		ToOld: to,
		InDb: len(result),
		Data: result[from:to],
	}
	re,er := json.Marshal(resp)

	if er != nil {
		fmt.Println("error")
		c.SendString("internal server error")
		return er
	}
	c.Send(re)
	return nil
}
func WaterMark(){
	fmt.Println("THIS API CREATE BY ARYARTE AND OPEN SOURCE")
	fmt.Println("GIVE ME START ON GITHUB http://github.com/aryarte")
	return 
}
func Listed(c *fiber.Ctx) error {
	tos := c.Query("to")
	froms := c.Query("from")
	if tos == "" || froms == "" {
		c.SendString("query empty")
		return errors.New("query empty")
	}
	to,_ := strconv.Atoi(tos)
	from,_ := strconv.Atoi(froms)
	type res struct {
		FromOld int `json="from_old"`
		ToOld int `json="to_old"`
		InDb int `json="db_length"`
		Data []primitive.M `json="data"`
	}
	ctx := context.TODO()
	cursor,err := Pr.Find(ctx,bson.M{})
	if err != nil {
		return err
	}
	result := []bson.M{}
	if err := cursor.All(ctx,&result);err != nil {
		return err
	}
	if to > len(result) {
		to = len(result)
	}
	if to < from {
		c.SendString("index out of range")
		return errors.New("equavalent error")
	}
	resp := res{
		FromOld: from,
		ToOld: to,
		InDb: len(result),
		Data: result[from:to],
	}
	re,er := json.Marshal(resp)
	//fmt.Println(string(re),resp)
	if er != nil {
		fmt.Println("error")
		c.SendString("internal server error")
		return er
	}
	c.Send(re)
	return nil
}
func Exist(c *fiber.Ctx) error {
	ctx := context.TODO()
	name := c.Query("name")
	newStock := c.Query("stock")
	storeid := c.Query("storeid")
	/**
		query == "" return error "empty query"
	**/
	if name == "" || newStock == "" || storeid == "" {
		c.SendString("empty query")
		return errors.New("empty query")
	}
	//find by product on database
	var result bson.M
	if err := Pr.FindOne(ctx,bson.M{"name_product":name}).Decode(&result);err != nil {
		c.SendString("product not found")
		return err
	}
	//verify the store id
	if result["name_product"] != storeid {
		c.SendString("storeid invalid")
		return errors.New("storeid invalid")
	}
	//add the stock value
	stockInInt,err := strconv.Atoi(newStock)
	if err != nil {
		c.Send([]byte("server error"))
		return err
	}
	newstock := result["stock_product"].(int) + stockInInt
	if _,err := Pr.UpdateOne(ctx,bson.M{"_id":result["_id"]},bson.M{"$set":bson.M{"stock_product":newstock}});err != nil {
		c.SendString("stock not change")
		return err
	}
	c.SendString("succes")
	return nil
}
func ListByOwner(c *fiber.Ctx) error {
	//query needed
	owner := c.Query("storeid")
	tos := c.Query("to")
	froms := c.Query("from")
	if owner == "" || tos == "" || froms == ""{
		c.Send([]byte("empty query"))
		return errors.New("empty query")
	}
	to,_ := strconv.Atoi(tos)
	from,_ := strconv.Atoi(froms)
	type res struct {
		FromOld int `json="from_old"`
		ToOld int `json="to_old"`
		InDb int `json="db_length"`
		Data []primitive.M `json="data"`
	}
	ctx := context.TODO()
	cursor,err := Pr.Find(ctx,bson.M{"store":owner})
	if err != nil {
		return err
	}
	result := []bson.M{}
	if err := cursor.All(ctx,&result);err != nil {
		return err
	}
	if to > len(result) {
		to = len(result)
	}
	if to < from {
		c.SendString("index out of range")
		return errors.New("equavalent error")
	}
	resp := res{
		FromOld: from,
		ToOld: to,
		InDb: len(result),
		Data: result[from:to],
	}
	re,er := json.Marshal(resp)
	//fmt.Println(string(re),resp)
	if er != nil {
		fmt.Println("error")
		c.SendString("internal server error")
		return er
	}
	c.Send(re)
	return nil
}

func AddToCart(c *fiber.Ctx) error {
	ctx := context.TODO()
	id := c.Query("userid")
	product := c.Query("product")
	count := c.Query("count")
	types := c.Query("type")
	fmt.Println(id,product,count)
	if id == "" || product == "" || count == "" {
		c.SendString("empty string")
		return errors.New("empty string")
	}
	//find this user 
	idUser,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.SendString("id is invalid please login first")
		return errors.New("id invalid")
	}
	//find user
	result := bson.M{}
	if err := UserDb.FindOne(ctx,bson.M{"_id":idUser}).Decode(&result);err != nil {
		return err
	}
	//parse product json string
	parseProduct := model.Product{}
	json.Unmarshal([]byte(product),&parseProduct)
	//get productname on database
	idProd,err := primitive.ObjectIDFromHex(parseProduct.Id)
	if err != nil {
		return err
	}
	resultProd := bson.M{}
	if err := Pr.FindOne(ctx,bson.M{"_id":idProd}).Decode(&resultProd);err != nil {
		return err
	}
	//declared needed var
	item := model.Cart{}
	Count,_ := strconv.Atoi(count)
	cartold := bson.A{}
	
	item.Count = Count
	item.StoreId = parseProduct.StoreId
	item.Product = parseProduct.Name
	item.Id = parseProduct.Id
	
	if types == "add" {
	/** 
		add new product
		or add to exist product
	**/
		if int(resultProd["stock_product"].(int32)) < Count {
			return errors.New("stock not available")
		}
		Count = int(resultProd["stock_product"].(int32)) - Count
		//add item to cart
		if result["cart"] != nil {
			cartold = result["cart"].(bson.A)
		}
		var found bool
		for _,v := range cartold {
			value := v.(bson.M)
			if value["name_product"] == item.Product && value["store"] == item.StoreId {
				value["count"] = int(value["count"].(int32)) + item.Count
				found = true
				break
			} 
		}
		if found == false {
			cartold = append(cartold,item)
		}
		c.SendString("succes counting")
	}else if types == "mines" {
		//code here
		if resultProd == nil {
			return errors.New("no product on cart")
		}
		Count = int(resultProd["stock_product"].(int32)) + Count
		//add item to cart
		if result["cart"] != nil {
			cartold = result["cart"].(bson.A)
		}
		var found bool
		for i,v := range cartold {
			value := v.(bson.M)
			if value["name_product"] == item.Product && value["store"] == item.StoreId {
				value["count"] = int(value["count"].(int32)) - item.Count
				if value["count"].(int) < 1 {
					value["count"] = 0
					fmt.Println("deleting...")
					cartold[i] = cartold[len(cartold)-1]
					cartold = cartold[:len(cartold)-1]
				}
				found = true
				break
			} 
		}
		if found == false {
			return errors.New("product not on cart")
		}
		c.SendString("succes counting")
	}else {
		return errors.New("nothing to do")
	}
		//update cart on db
	if _,err := UserDb.UpdateOne(ctx,bson.M{"_id":idUser},bson.M{"$set":bson.M{"cart":cartold}});err != nil {
		return err
	}
	//update Product on db
	if _,err := Pr.UpdateOne(ctx,bson.M{"_id":idProd},bson.M{"$set":bson.M{"stock_product":Count}});err != nil {
		return err
	}
	return nil
}
func CheckOut(c *fiber.Ctx) error {
	ctx := context.TODO()
	money := c.Query("money")
	userid := c.Query("id")
	if userid == "" && money == ""{
		return errors.New("empty query")
	}
	str,err := strconv.Atoi(money)
	if str == 0 {
		return errors.New("invalid structur of money")
	}
	//find user by their id
	objId,err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}
	result := bson.M{}
	if err := UserDb.FindOne(ctx,bson.M{"_id":objId}).Decode(&result);err != nil{
		return err
	}
	cart := result["cart"].(bson.A)
	if len(cart) == 0 {
		c.SendString("cart empty")
		return nil
	}
	var prices int
	for _,v := range cart {
		value := v.(map[string]interface{})
		id,err := primitive.ObjectIDFromHex(value["productid"].(string))
		if err != nil {
			return err
		}
		result := bson.M{}
		if err := Pr.FindOne(ctx,bson.M{"_id":id}).Decode(&result);err != nil {
			return nil
		}
		countInInt := value["count"].(int)
		priceInInt := result["price_int"].(int)
		prices = prices + (countInInt * priceInInt)
	}
	empty := bson.A{}
	if _,err := UserDb.UpdateOne(ctx,bson.M{"_id":objId},bson.M{"$set":bson.M{"cart":empty}});err != nil {
		return err
	}
	return nil
}
func CheckCart(c *fiber.Ctx) error{
	ctx := context.TODO()
	userid := c.Query("id")
	if userid == "" {
		return errors.New("empty query")
	}
	//find user by their id
	objId,err := primitive.ObjectIDFromHex(userid)
	if err != nil {
		return err
	}
	result := bson.M{}
	if err := UserDb.FindOne(ctx,bson.M{"_id":objId}).Decode(&result);err != nil{
		return err
	}
	jsonstr,_ := json.Marshal(result["cart"])
	c.Send(jsonstr)
	return nil
}
