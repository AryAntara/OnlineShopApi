package control
import(
	"context"
	"fmt"
	model "RESTAPI/model"
	"RESTAPI/database"
	"github.com/gofiber/fiber/v2"
	bcrypt "golang.org/x/crypto/bcrypt"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"RESTAPI/config"
	"time"
	"errors"
)
var UserDb = db.Collection(config.DbConn,config.Db,config.UserColl)
func valid(username string) bool {
	var result bson.M
	//check username has usage or not
	err := UserDb.FindOne(context.TODO(),bson.D{{"name",username}}).Decode(&result)
	if err != nil {
		return true
	}
	return false
}
func recovery(){
	recover()
}
func SignUp(c *fiber.Ctx) error {
	defer recovery()
	User := model.User{}
	Username := c.Query("username")
	Password := c.Query("password")
	Hashing,err := bcrypt.GenerateFromPassword([]byte(Password),bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	User.Name = Username
	User.Password = string(Hashing)
	if valid(Username) == false {
		c.SendString("> user has been used")
		return errors.New("user has been used")
	}
	result,err := UserDb.InsertOne(context.Background(),User)
	if err != nil{
		return err
	}
	fmt.Println(":: user added with id",*result)
	c.SendString("> user has been added")
	return nil
}
func User(c *fiber.Ctx) error {
	defer recovery()
	ids := strings.ReplaceAll(config.ID, `"`,"")
	id,err := primitive.ObjectIDFromHex(ids)
	if err != nil {
		fmt.Println(":: err >",err)
		return err
	}
	var result bson.M
	err = UserDb.FindOne(context.TODO(),bson.D{{"_id",id}}).Decode(&result)
	if err != nil {
		fmt.Println(":: err >",err)
		c.SendString("> username not found ")
		return err
	}
	type StrictUser struct {
		Name string `json:"name"`
		Id interface{} `json:"uuid"`
	}
	user := StrictUser{}
	user.Name = "annonymus"
	user.Id = result["_id"]
	JSONStr,_ := json.Marshal(user)
	c.SendString(string(JSONStr))
	return nil
}
func Login(c *fiber.Ctx) error {
	pass := c.Query("password")
	username := c.Query("username")
	if pass == "" && username == "" {
		c.Send([]byte("> empty string is thrown error"))
		return errors.New("empty string on query")
	}
	if valid(username) != false {
		c.Send([]byte("> user not found"))
		return errors.New("user not signup")
	}
	var result bson.M
	err := UserDb.FindOne(context.Background(),bson.D{{"name",username}}).Decode(&result)
	if err != nil{
		c.Send([]byte("> user not found"))
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(result["password"].(string)),[]byte(pass));err != nil {
		c.Send([]byte("> password wrong"))
		return err
	}
	uuid := result["_id"]
	Uuid,_ := json.Marshal(uuid)
	KeyScreet := config.KEY
	claims := &jwt.RegisteredClaims{
		ID: string(Uuid),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 45)),
	}
	ClaimForRefreshToken := &jwt.RegisteredClaims{
		ID: string(Uuid),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
	}
	RefreshToken,_ := jwt.NewWithClaims(jwt.SigningMethodHS256,ClaimForRefreshToken).SignedString([]byte(KeyScreet))
	Sign := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	token,err := Sign.SignedString([]byte(KeyScreet))
	if err != nil {
		c.Send([]byte("> Internal Server Error"))
		return err
	}
	c.SendString(token)

	//update db
	UserDb.ReplaceOne(
		context.TODO(),
		bson.M{"_id":result["_id"]},
		bson.M{
			"RefreshToken":RefreshToken,
			"name":username,
			"password":result["password"],
			"cart":result["cart"],
		})
	//set session for this user
	sess,err := config.Session.Get(c)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	sess.Set("key",RefreshToken)
	if err := sess.Save();err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func Token(c *fiber.Ctx) error {
	sess,_ := config.Session.Get(c)
	Refresh := sess.Get("key")
	//Name := sess.Get("name")
	if Refresh == nil {
		c.Send([]byte("> you not login"))
		return errors.New("user not login first")
	}
	type custom struct {
		jwt.RegisteredClaims
	}
	T,err:= jwt.ParseWithClaims(Refresh.(string),&custom{},func(t *jwt.Token)(interface{},error){
		return []byte(config.KEY),nil
	})
	if err != nil {
		c.SendString("token is invalid")
		return err
	}
	//verify cookie from  db
	if claims,ok := T.Claims.(*custom); ok && T.Valid {
		ID := claims.RegisteredClaims.ID
		Id := strings.ReplaceAll(ID,`"`,"")
		id,_ := primitive.ObjectIDFromHex(Id)
		result := bson.M{}
		err := UserDb.FindOne(context.TODO(),bson.M{"_id":id}).Decode(&result)
		if err != nil {
			fmt.Println(err)
			c.Send([]byte("> user not found on database"))
			return err
		}
		KeyScreet := config.KEY
		claims := &jwt.RegisteredClaims{
			ID: strings.ReplaceAll(claims.RegisteredClaims.ID,`"`,""),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 45)),
		}	
		Sign := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
		token,err := Sign.SignedString([]byte(KeyScreet))
		if err != nil {
			c.Send([]byte("> Internal Server Error"))
			return err
		}
		c.SendString(token)
	}
	return nil

}
func Logout(c *fiber.Ctx) error {
	s,_ := config.Session.Get(c)
	s.Destroy()
	return nil
}
