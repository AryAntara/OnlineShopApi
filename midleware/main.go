package token
import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
	"RESTAPI/config"
	"fmt"
)
func log(str interface{}){
	fmt.Println("::",str)
}
func Verify(c *fiber.Ctx) error {
	Token := c.Query("key")
	if Token == "" {
		c.SendString("> cannot verify")
		log("token empty")
		return nil
	}
	type custom struct {
		jwt.RegisteredClaims
	}
	token,err:= jwt.ParseWithClaims(Token,&custom{},func(t *jwt.Token)(interface{},error){
		return []byte(config.KEY),nil
	})
	if err != nil {
		c.Send([]byte("> secret key Expired"))
		log("error on parse")
		log(err)
		return nil
	}
	if claims,ok := token.Claims.(*custom); ok && token.Valid {
		config.ID = claims.RegisteredClaims.ID
	} else {
		log(err)
		c.Send([]byte("> secret key Expired"))
		return nil
	}
	log("succes verify for user "+config.ID)
	c.Next()
	return nil
}
