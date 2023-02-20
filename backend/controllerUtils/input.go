package controllerUtils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type SignUpInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func InputSignUp(c *gin.Context) (SignUpInput, error) {
	var in SignUpInput
	if err := c.ShouldBindJSON(&in); err != nil {
		return in, err
	}
	if in.Name == "" || in.Password == "" {
		return in, fmt.Errorf("not found name or password")
	}
	return in, nil
}


