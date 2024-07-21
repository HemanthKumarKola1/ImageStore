package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ImageStore/user"
)

func server(dbUser *user.DbUser) {
	r := gin.Default()

	r.RedirectTrailingSlash = true

	// Route to register a new user
	r.POST("/register", dbUser.RegisterHandler)

	// Route to verify OTP
	r.GET("/verify/:username", dbUser.VerifyOTPHandler)

	// Post files of a user
	r.POST("/files/:username", dbUser.PostData)

	// Get files of a user
	r.GET("/files/:username", dbUser.GetData)

	// Start server
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
