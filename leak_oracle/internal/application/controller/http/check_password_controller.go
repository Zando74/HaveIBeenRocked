package http_controller

import (
	"github.com/Zando74/IHaveBeenRocked/leak_oracle/internal/domain/usecase"
	"github.com/gin-gonic/gin"
)

type CheckPasswordController struct {
	CheckPasswordUseCase *usecase.CheckPasswordPresence
}

func (controller *CheckPasswordController) CheckPasswordPresence(c *gin.Context) {
	var requestPayload CheckPasswordValidator
	if err := c.ShouldBindJSON(&requestPayload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	isPresent, err := controller.CheckPasswordUseCase.Execute(requestPayload.Password)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !isPresent {
		c.JSON(200, CheckPasswordResponse{Found: false})
		return
	}

	c.JSON(200, CheckPasswordResponse{Found: true})

}
