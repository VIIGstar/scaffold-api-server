package user_handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"scaffold-api-server/internal/dtos"
	"scaffold-api-server/internal/entities"
	database "scaffold-api-server/internal/services/database/mysql"
	"scaffold-api-server/pkg/auth"
	app_http "scaffold-api-server/pkg/http"
)

// @Summary  Signup create new user
// @Tags     Investor
// @Param    data body dtos.InvestorDTO true "The input struct"
// @Accept   json
// @Produce  json
// @Success  200      {object}  auth.Authentication
// @Failure  400,500  {object}  object{error=string}
// @Router   /api/v1/investors/signup [post]
func (s *userHandler) Signup(c *gin.Context) {
	user := dtos.InvestorDTO{}
	err := c.ShouldBindJSON(&user)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Invalid request, err: %v", err))
		c.JSON(http.StatusBadRequest, auth.Authentication{
			Error: app_http.HTTPBadRequestError,
		})
		return
	}

	if !user.IsValid() {
		s.logger.Error(fmt.Sprintf("Invalid dto value, dto: %v", user))
		c.JSON(http.StatusBadRequest, auth.Authentication{
			Error: app_http.HTTPBadRequestError,
		})
		return
	}

	iEntity, err := user.ToEntity()
	if err != nil {
		s.logger.Error(fmt.Sprintf("Parse to entity failed, err: %v", err))
		c.JSON(http.StatusInternalServerError, auth.Authentication{
			Error: app_http.HTTPInternalServerError,
		})
		return
	}

	entity, _ := iEntity.(*entities.Investor)
	err = s.repo.Database().User().Create(entity)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Insert into database failed, err: %v", err))
		if database.IsDuplicateErr(err) {
			c.JSON(http.StatusInternalServerError, auth.Authentication{
				Error: "Already registered!",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, app_http.HTTPInternalServerError)
		return
	}

	c.JSON(http.StatusOK, auth.Authentication{
		AccessToken: auth.New().GenerateAccessToken(entity.Address, fmt.Sprintf("%v", entity.ID), c),
		Success:     true,
	})
}
