package tokens

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/v4"
	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/red-star25/advance-go/database"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserData *mongo.Collection = database.UserData(database.Client, "Users  ")
var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	uID       string
	jwt.StandardClaims
}

func TokenGenerator(email string, firstName string, lastName string, uID string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		uID:       uID,
		StandardClaims: jwt.StandClamis{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken
}

func ValidateToken() {}

func UpdateAllToken() {}
