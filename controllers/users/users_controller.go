package users

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imaarov/bookstore_microservice/domain/users"
	"github.com/imaarov/bookstore_microservice/services"
	"github.com/imaarov/bookstore_microservice/utils/errors"
)

func TestServiceInterface() {
}

func getUserId(userIdParam string) (int64, *errors.RestErr) {
	userId, userErr := strconv.ParseInt(userIdParam, 10, 64)
	if userErr != nil {
		return 0, errors.NewBadRequestError("Invalid User Id")
	}
	return userId, nil
}

func Create(c *gin.Context) {
	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		//TODO: Handle Json Error -> RETURN BAD REQUEST TO THE CALLER
		restErr := errors.NewBadRequestError("Invalid Json Body")
		c.JSON(restErr.Status, restErr)
		fmt.Println("Json Err:", err.Error())
		return
	}

	result, saveErr := services.UsersService.CreateUser(user)
	if saveErr != nil {
		//TODO:Handle User Creation Error
		c.JSON(saveErr.Status, saveErr)
		return
	}
	publicRequest := c.GetHeader("X-public") == "true"
	c.JSON(http.StatusCreated, result.Marshall(publicRequest))
}

func Get(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	user, getErr := services.UsersService.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		return
	}
	publicRequest := c.GetHeader("X-public") == "true"
	c.JSON(http.StatusOK, user.Marshall(publicRequest))
}

func Update(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId

	isPartial := c.Request.Method == http.MethodPatch

	result, updateErr := services.UsersService.UpdateUser(isPartial, user)
	if updateErr != nil {
		c.JSON(updateErr.Status, updateErr)
		return
	}
	publicRequest := c.GetHeader("X-public") == "true"
	c.JSON(http.StatusOK, result.Marshall(publicRequest))
}

func Delete(c *gin.Context) {
	userId, idErr := getUserId(c.Param("user_id"))
	if idErr != nil {
		c.JSON(idErr.Status, idErr)
		return
	}

	if err := services.UsersService.DeleteUser(userId); err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{"status": "Deleted"})
}

func FindUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "<h1>implement me! ???</h1>")
}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.UsersService.Search(status)
	fmt.Println(err)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	publicRequest := c.GetHeader("X-public") == "true"
	c.JSON(http.StatusOK, users.Marshall(publicRequest))
}
