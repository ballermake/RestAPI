package user

import (
	"RestAPI/internal/domain"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		db: database,
	}
}

func (us UserService) Register(ctx *gin.Context) {
	var user domain.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	if user.Name == "" {
		ctx.JSON(400, gin.H{
			"message": "field name required",
		})
		return
	}

	if user.Email == "" {
		ctx.JSON(400, gin.H{
			"message": "field email should not be blank",
		})
		return
	}

	if user.Password == "" {
		ctx.JSON(400, gin.H{
			"message": "field password cannot be empty",
		})
		return
	}

	if len(user.Password) < 6 {
		ctx.JSON(400, gin.H{
			"message": "password has to be minimal 6 characters",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	//err = us.db.Create(&user).Error
	if err := us.db.Create(&user).Error; err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when creating user",
		})
		return
	}

	token, err := generateJWT(user.ID)
	//fmt.Println("===", err)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when creating user",
		})
		return
	}
	ctx.JSON(201, gin.H{
		"token": token,
	})

}

func (us UserService) Login(ctx *gin.Context) {
	var currentUser domain.User
	err := ctx.ShouldBind(&currentUser)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	var user domain.User
	err = us.db.Where("email = ?", currentUser.Email).Take(&user).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "maaf email/password salah",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentUser.Password))
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "fail when retrieving user",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
	})

}

var signatureKey = []byte("mySuperSecretSignature")

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iss":     "edspert",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(signatureKey)
	if err != nil {
		return "", err
	}
	return stringToken, nil
}

func (us UserService) DecriptJWT(token string) (map[string]interface{}, error) {
	//fmt.Println("====", token)
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth invalid")
		}
		return signatureKey, nil
	})

	//fmt.Println("===", err)
	data := make(map[string]interface{})
	if err != nil {
		return data, err
	}
	fmt.Println("===", parsedToken.Valid)
	if !parsedToken.Valid {
		return data, errors.New("token invalid")
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}
