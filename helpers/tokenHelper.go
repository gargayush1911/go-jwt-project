package helpers

import (
	"context"
	"go-jwt-project/database"
	"log"
	"os"
	"time"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SignedDetils struct{
	Email string
	First_Name string
	Last_Name string
	Uid string
	User_type string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var SECRET_KEY string = os.Getenv("SECTRET_KEY")

func GenerateAllTokens(email string,Firstname string,Lastname string,Usertype string,Userid string) (signedToken string, signedRefreshToken string, err error){
	claims:= &SignedDetils{
		Email: email,
		First_Name: Firstname,
		Last_Name: Lastname,
		Uid: Userid,
		User_type: Usertype,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24*time.Hour)),
		},
	}
	refreshClaims:= &SignedDetils{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Hour*time.Duration(168))),
		},
	}
	token , err:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims).SignedString([]byte(SECRET_KEY))
	refreshToken , err:= jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims).SignedString([]byte(SECRET_KEY))

	if err!= nil{
		log.Panic(err)
		return 
	}
	return token,refreshToken,err
}

func UpdateAllTokens(signedtoken string, signedrefreshToken string, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var Updateobj bson.D
	Updateobj = append(Updateobj, bson.E{Key: "token", Value: signedtoken})
	Updateobj = append(Updateobj, bson.E{Key: "refresh_token", Value: signedrefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	Updateobj = append(Updateobj, bson.E{Key: "updated_at", Value: Updated_at})

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOne().SetUpsert(true)

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{{Key: "$set", Value: Updateobj}},
		opt,
	)
	if err != nil {
		log.Panic(err)
		return
	}
}

func ValidateToken(signedToken string) (claims *SignedDetils, msg string){
	token ,err:= jwt.ParseWithClaims(
		signedToken,
		&SignedDetils{},
		func(token *jwt.Token)(interface{}, error){
			return []byte(SECRET_KEY),nil
		},
	)
	if err != nil {
    msg = err.Error()
    return
}

	claims, ok := token.Claims.(*SignedDetils)
	if !ok {
	    msg = "the token is invalid"
	    return
	}

	if claims.ExpiresAt == nil || claims.ExpiresAt.Before(time.Now().Local()) {
	msg = "token is expired"
	return
	}

	return claims, msg

}