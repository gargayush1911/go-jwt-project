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
	"github.com/playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection  *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func Getusers()

func Getuser()

