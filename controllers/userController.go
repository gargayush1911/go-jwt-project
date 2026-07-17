package controllers

import (
	"context"
	"fmt"
	"go-jwt-project/database"
	helper "go-jwt-project/helpers"
	"go-jwt-project/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection  *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword(password string) (string){
	passWord ,err:= bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil{
		log.Panic(err)
	}
	return string(passWord)
}

func VerifyPassword(userpassword string,providedpassword string) (bool,string){
	err:= bcrypt.CompareHashAndPassword([]byte(providedpassword),[]byte(userpassword))
	check := true
	msg := ""
	if err!=nil{
		msg = fmt.Sprintf("password of email is incorrect")
		check = false
	}
	return check,msg
}

func Signup() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var  user models.User

		if err:= c.BindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":validationErr.Error()})
		}

		count,err:= userCollection.CountDocuments(ctx,bson.M{"email":user.Email})
		if err!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking the email"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count ,err1 := userCollection.CountDocuments(ctx,bson.M{"phone":user.Phone})
		defer cancel()
		if err1!=nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while checking phone number"})
		}

		if count>0{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"this emiail or phone number already exists"})
		}

		user.Created_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))
		user.Updated_at,_= time.Parse(time.RFC3339,time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token,refreshtoken,err:= helper.GenerateAllTokens(*user.Email,*user.First_name,*user.Last_name,*user.User_type,*&user.User_id)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user already exists"})
		}
		user.Token = &token
		user.Refresh_token = &refreshtoken

		resultInsertionNumber,inserterr := userCollection.InsertOne(ctx,user)
		if inserterr!=nil{
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK,resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		ctx,cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var Founduser models.User

		if err:= c.BindJSON(&user);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}
		err:= userCollection.FindOne(ctx,bson.M{"email":user.Email}).Decode(&Founduser)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"email or password is incorrect"})
			return
		}

		passwordisValid , msg := VerifyPassword(*user.Password,*Founduser.Password)
		defer cancel()
		if passwordisValid!=true{
			c.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		if Founduser.Email == nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"user not found"})

		}
		token,refreshToken,_:=helper.GenerateAllTokens(*Founduser.Email,*Founduser.First_name,*Founduser.Last_name,*Founduser.User_type,*&Founduser.User_id)

		helper.UpdateAllTokens(token,refreshToken,Founduser.User_id)
		userCollection.FindOne(ctx,bson.M{"user_id":Founduser.User_id}).Decode(&Founduser)

		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,Founduser)
	}
}	

func Getusers() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN");err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		var ctx,cancel = context.WithTimeout(context.Background(),100*time.Second)

		recordPerPage,err:=strconv.Atoi(c.Query("recordPerPage"))
		if err!=nil || recordPerPage<1{
			recordPerPage = 10
		}
		page , err1:=strconv.Atoi(c.Query("page"))
		if err1!=nil || page<1{
			page =1
		}

		startIndex := (page-1)*recordPerPage
		startIndex,err = strconv.Atoi(c.Query("startIndex"))
		
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil},{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},{Key: "data",Value: bson.D{{Key: "$push",Value: "$$Root"}}}}}}
		projectStage := bson.D{
		{Key: "$project", Value: bson.D{
		{Key: "_id", Value: 0},
		{Key: "total_count", Value: 1},
		{Key: "user_items", Value: bson.D{
			{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
			}},
		}},}
		result, err := userCollection.Aggregate(ctx,mongo.Pipeline{
			matchStage,groupStage,projectStage,
		})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":"error occured while listing user items"})
		}
		var allUsers []bson.M
		if err = result.All(ctx,&allUsers); err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK,allUsers[0])
	}
}

func Getuser() gin.HandlerFunc{
	return func(c *gin.Context){
		userid :=  c.Param("user_id")
		
		if err:=helper.MatchUserToUid(c,userid);err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		}
		var ctx, cancel = context.WithTimeout(context.Background(),10*time.Second)

		var user models.User

		err:= userCollection.FindOne(ctx,bson.M{"user_id":userid}).Decode(&user)
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,user)
	}
}

