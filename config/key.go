package config
import (
	"github.com/gofiber/fiber/v2/middleware/session"
	"RESTAPI/database"
)
// secret key fot jwt token
var KEY = "AkuGaPunyaAyank"
var ID string
// cookie setup
var Session = session.New(session.Config{
	CookieName: "refresh_token",
	CookieHTTPOnly: true,
})
// database configuration
var DbConn = db.Init()
var Db = "tookie" // database name
var UserColl = "user"  // user collection
var ProdColl = "product" // product collection
var StoreColl = "store" // store collection
