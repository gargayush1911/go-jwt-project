package middleware

import(
	"fmt"
	"net/http"
	helper "go-jwt-project/helpers"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc{
	return func(c *gin.Context) {
		clientToken :=  c.Request.Header.Get("token")
		if clientToken == ""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":fmt.Sprintf("no authorized header provided")})
			c.Abort()
			return
		}

		claims,err:=helper.ValidateToken(clientToken)
		if err!=""{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name",claims.First_Name)
		c.Set("last_name",claims.Last_Name)
		c.Set("uid",claims.Uid)
		c.Set("user_type",claims.User_type)
		c.Next()

	}
}