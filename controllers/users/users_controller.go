package users

import (
	"github.com/nicoletafratila/bookstore_oauth-go/oauth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicoletafratila/bookstore_users-api/domain/users"
	"github.com/nicoletafratila/bookstore_users-api/services"
	"github.com/nicoletafratila/bookstore_utils-go/rest_errors"
)

func getUserId(userIdParam string) (int64, rest_errors.RestErr) {
	userId, userIdErr := strconv.ParseInt(userIdParam, 10, 64)
	if userIdErr != nil {
		return 0, rest_errors.NewBadRequestError("user id should be a number")
	}
	return userId, nil
}

func Create(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, saveErr := services.UsersService.Create(user)
	if saveErr != nil {
		c.JSON(saveErr.Status(), saveErr)
		return
	}
	c.JSON(http.StatusCreated, result.Marshall(oauth.IsPublic(c.Request)))
}

func Get(c *gin.Context) {
	if err := oauth.AuthenticateRequest(c.Request); err != nil {
		c.JSON(err.Status(), err)
		return
	}

	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	user, getErr := services.UsersService.Get(userId)
	if getErr != nil {
		c.JSON(getErr.Status(), getErr)
		return
	}

	if oauth.GetCallerId(c.Request) == user.Id {
		c.JSON(http.StatusOK, user.Marshall(false))
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.Update(isPartial, user)
	if updateErr != nil {
		c.JSON(updateErr.Status(), updateErr)
		return
	}
	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status(), idErr)
		return
	}

	deleteErr := services.UsersService.Delete(userId)
	if deleteErr != nil {
		c.JSON(deleteErr.Status(), deleteErr)
		return
	}
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func Search(c *gin.Context) {
	status := c.Query("status")

	result, err := services.UsersService.Search(status)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, result.Marshall(oauth.IsPublic(c.Request)))
}

func Login(c *gin.Context) {
	var request users.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		restErr := rest_errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status(), restErr)
		return
	}
	user, err := services.UsersService.Login(request)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, user.Marshall(oauth.IsPublic(c.Request)))
}
