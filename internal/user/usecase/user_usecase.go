package usecase

import (
	"course-api/internal/domain"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	db *gorm.DB
}

func NewUserUseCase(db *gorm.DB) *UserUseCase {
	return &UserUseCase{
		db: db,
	}
}

func (userUsecase *UserUseCase) Register(c *gin.Context) {
	type RegisterRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var registerRequest RegisterRequest
	if err := c.BindJSON(&registerRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
		return
	}

	if registerRequest.Name == "" || registerRequest.Email == "" || registerRequest.Password == "" {
		c.JSON(400, gin.H{
			"message": "Isi semua field",
		})
		return
	}

	if len(registerRequest.Password) < 6 {
		c.JSON(400, gin.H{
			"message": "Password minimal 6 karakter",
		})
		return
	}

	user := domain.NewUser(registerRequest.Name, registerRequest.Email, registerRequest.Password)
	if err := userUsecase.db.Create(user).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "cannot create user",
		})
		return
	}

	token, err := user.GenerateToken()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "cannot create token",
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "user created successfully",
		"token":   token,
	})
}

func (userUsecase *UserUseCase) DecryptJWT(token string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return domain.PrivateKey, nil
	})

	if err != nil {
		return map[string]interface{}{}, err
	}

	if !parsedToken.Valid {
		return map[string]interface{}{}, err
	}
	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (userUsecase *UserUseCase) Login(c *gin.Context) {
	var userRequest domain.User
	err := c.ShouldBind(&userRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid input",
		})
		return
	}
	if userRequest.Email == "" || userRequest.Password == "" {
		c.JSON(400, gin.H{
			"message": "email/password required",
		})
		return
	}

	var user domain.User
	err = userUsecase.db.Where("email = ?", userRequest.Email).Take(&user).Error
	if err != nil || user.ID == 0 {
		c.JSON(400, gin.H{
			"message": "wrong email/password",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "wrong email/password",
		})
		return
	}
	token, _ := user.GenerateToken()
	c.JSON(200, gin.H{
		"token": token,
	})
}
