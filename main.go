package main

import (
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	mysecret = "secret"
	counter  = 0
)

type Info struct {
	Id       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Event    string `json:"event" binding:"required"`
	Image1   string `json:"image1"`
	Image2   string `json:"image2"`
	Image3   string `json:"image3"`
	Image4   string `json:"image4"`
	Image5   string `json:"image5"`
}

func main() {
	r := gin.Default()

	public := r.Group("/public")
	public.GET("/", publicHandler)

	private := r.Group("/private")
	private.Use(jwt.Auth(mysecret))
	private.GET("/", privateHandler)
	private.POST("/create", createHandler)
	// private.POST("/upload", uploadHandler)

	r.Run()
}

func publicHandler(c *gin.Context) {
	token := jwtgo.New(jwtgo.GetSigningMethod("HS256"))
	token.Claims = jwtgo.MapClaims{
		"iss": "golang photo upload app",
	}

	tokenString, err := token.SignedString([]byte(mysecret))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Could not create JWT token!",
		})
		return
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

func createHandler(c *gin.Context) {
	var json Info
	json.Id = counter + 1

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, json)
}
