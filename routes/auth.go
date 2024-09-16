package routes

import (
	"fmt"
	"net/http"

	"github.com/Atheer-Ganayem/Go-social-media-api/models"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func signup(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBind(&user)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data. Invalid input."})
		return
	}

	filePath, code, err := utils.HandleImage(ctx.Request, "avatar")
	if err != nil {
		ctx.JSON(code, gin.H{"message": err.Error()})
		return
	}
	user.Avatar = filePath


	exits := models.CheckUserExists(user.Email)
	if exits {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email already exits. Please choose another email."})
		return
	} 

	err = user.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An errorr occurred, Couldn't create user."})
		return
	}

	token ,err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() + ". (User created, you can login)"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully.", "token": token})
}

func login(ctx *gin.Context) {
	var data models.LoginRequest
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please enter a valid email and password."})
		return
	}

	user, err := models.GetUser(bson.M{"email": data.Email})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password."})
		return
	}

	if !utils.CheckPassword(data.Password, user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password."})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't generate auth token."})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged in successfully.", "token": token})
}


