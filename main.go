package main

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

var (
	mysecret = "secret"
)

func main() {
	r := gin.Default()

	public := r.Group("/public")
	public.GET("/", publicHandler)

	private := r.Group("/private")
	private.Use(jwt.Auth(mysecret))
	private.GET("/", privateHandler)

	r.Run()
}

func publicHandler(c *gin.Context) {
	token := jwtgo.New(jwtgo.GetSigningMethod("HS256"))
	token.Claims = jwtgo.MapClaims{
		"username": "moeabdol",
	}

	tokenString, err := token.SignedString([]byte(mysecret))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Could not create JWT token!",
		})
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

func privateHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Successfully authorized",
	})
}
