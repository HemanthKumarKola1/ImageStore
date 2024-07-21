package user

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"-"`
}

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (r *DbUser) RegisterHandler(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.Contains(newUser.Password, ":") {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"sorry": "we can't support : in the password"})
	}
	err, val := r.UserExists(newUser.Username)
	if err == redis.Nil && len(val) != 0 {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "User already Exists"})
		return
	} else if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate OTP
	otp := GenerateOTP()
	fmt.Printf("otp for the %s is %s", newUser.Username, otp)

	r.SetUser(newUser.Username, newUser.Password+":"+otp)
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your mobile number", "otp": otp})
}

// VerifyOTPHandler handles OTP verification
func (r *DbUser) VerifyOTPHandler(c *gin.Context) {
	username := c.Param("username")
	otp := c.Query("otp")

	err, user := r.UserExists(username)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "OTP EXPIRED"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	parts := strings.SplitN(user, ":", 2)

	if otp != parts[1] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	r.RClient.Set(username, user, 365*24*time.Hour).Result()

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully", "username": username, "password": user[0]})
}

func (r *DbUser) PostData(c *gin.Context) {
	if !r.verify(c) {
		return
	}
	username := c.Param("username")
	dirname := c.Query("dirname")
	imageUrl := c.Query("imageurl")

	insertQuery := `
        INSERT INTO user_files (username, directory, imageurl) VALUES (?, ?, ?);
    `

	if err := r.CassConnection.Query(insertQuery, username, dirname, imageUrl).Exec(); err != nil {
		log.Fatalf("Unable to insert data: %v", err)
	}

	log.Println("Data inserted successfully")
}

func (r *DbUser) GetData(c *gin.Context) {
	if !r.verify(c) {
		return
	}
	username := c.Param("username")
	dirname := c.Query("dirname")

	insertQuery := `
        SELECT * FROM user_files WHERE username = ? and dirname = ?;
    `

	if err := r.CassConnection.Query(insertQuery, username, dirname).Exec(); err != nil {
		log.Fatalf("Unable to insert data: %v", err)
	}

	log.Println("Data inserted successfully")

}

func (r *DbUser) verify(c *gin.Context) bool {
	username := c.Param("username")
	pwd := c.Query("password")

	err, rPwd := r.UserExists(username)
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Incorrect username/password"})
			return false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if pwd != rPwd {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Incorrect username/password"})
		return false
	}
	return true
}
