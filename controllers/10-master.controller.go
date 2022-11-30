package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesh/models"
	"golang.org/x/crypto/bcrypt"
)

func MasterCreateRegistrar(c *gin.Context) {
	// VALIDATE MASTER KEY

	var requestBody models.Registrar

	// BIND REQUEST BODY TO MODEL
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "could not bind request body",
			"err": err.Error,
		})
		return
	}

	// CHECK IF REGISTRAR ALREADY EXISTS
	var dbData models.Registrar
	models.DB.Where("regID = ?", requestBody.RegID).First(&dbData)

	if dbData.RegID != "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "entity already exist",
		})
		return
	}

	// HASH THE PASSWORD
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(requestBody.Hash), bcrypt.DefaultCost)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not hash the password",
			"err": err.Error(),
		})
		return
	}

	requestBody.Hash = string(hashedPass)

	// INSERT ACTION IN DATABASE
	if err := models.DB.Create(&requestBody); err.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "could not create user",
			"err": err.Error,
		})
		return
	}

	// RESPOND STATUS
	c.JSON(http.StatusCreated, gin.H{
		"msg": "entity created",
	})

}
