package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"strconv"
)

var (
	mysecret = "secret"
	counter  = 0
	infos    []Info
	conf     = aws.Config{Region: aws.String("eu-central-1")}
	sess     = session.New(&conf)
	svc      = s3manager.NewUploader(sess)
	bucket   = "lean-tobacco-bucket"
)

type Info struct {
	Id       int      `json:"id" binding:"required"`
	Username string   `json:"username" binding:"required"`
	Event    string   `json:"event" binding:"required"`
	Images   []string `json:"images"`
}

func main() {
	r := gin.Default()

	public := r.Group("/public")
	public.GET("/", publicHandler)

	private := r.Group("/private")
	private.Use(jwt.Auth(mysecret))
	private.GET("/", privateHandler)
	private.POST("/create", createHandler)
	private.POST("/upload/:id", uploadHandler)
	// private.Static("/uploads", "./uploads")

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
	counter++
	json.Id = counter

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	infos = append(infos, json)

	c.JSON(http.StatusCreated, json)
}

func uploadHandler(c *gin.Context) {
	var info Info
	var json Info

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range infos {
		if infos[i].Id == id {
			info = infos[i]
		}
	}

	json.Id = info.Id
	json.Username = info.Username
	json.Event = info.Event

	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		uuId := uuid.Must(uuid.NewV4())

		c.SaveUploadedFile(file, "uploads/"+uuId.String())

		file, err := os.Open("./uploads/" + uuId.String())
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer file.Close()

		result, err := svc.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(uuId.String()),
			Body:   file,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		json.Images = append(json.Images, result.Location)

		for i := range infos {
			if infos[i].Id == id {
				infos[i].Images = append(infos[i].Images, result.Location)
			}
		}
	}

	c.JSON(http.StatusCreated, json)
}
