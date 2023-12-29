package generatetoken

import (
	"log"
	"os"
	"time"

	"github.com/hienphan0111/movie-review-api/database"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type JwtSignedDetails struct {
	Email     string
	Name      string
	Username  string
	Uid       string
	User_type string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(
	email string,
	name string,
	username string,
	uid string,
	userType string,
) (
	signedToken string,
	signedRefreshToken string,
	err error,
) {
	claims := &JwtSignedDetails{
		Email:     email,
		Name:      name,
		Username:  username,
		Uid:       uid,
		User_type: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(12)).Unix(),
		},
	}

	refreshClaims := &JwtSignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}
