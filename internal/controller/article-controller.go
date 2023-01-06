package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/dto"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/helper"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/models"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ArticleController interface {
	All(context *gin.Context)
	FindByID(context *gin.Context)
	Insert(context *gin.Context)
	Update(context *gin.Context)
	Delete(context *gin.Context)
}

type articleController struct {
	articleService service.ArticleService
	jwtService     service.JWTService
}

func NewArticleController(articleService service.ArticleService, jwtServ service.JWTService) ArticleController {
	return &articleController{
		articleService: articleService,
		jwtService:     jwtServ,
	}
}

func (c *articleController) All(context *gin.Context) {
	var articles []models.Article = c.articleService.All()
	res := helper.BuildResponse(true, "OK", articles)
	context.JSON(http.StatusOK, res)
}

func (c *articleController) FindByID(context *gin.Context) {
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("No parameter id was found", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var article models.Article = c.articleService.FindByID(id)
	if (article == models.Article{}) {
		res := helper.BuildErrorResponse("Data not found", "No data is found with given ID", helper.EmptyObj{})
		context.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK", article)
		context.JSON(http.StatusOK, res)
	}
}

func (c *articleController) Insert(context *gin.Context) {
	var articleCreateDTO dto.ArticleCreateDTO
	errDTO := context.ShouldBind(&articleCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			articleCreateDTO.UserID = convertedUserID
		}
		result := c.articleService.Insert(articleCreateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *articleController) Update(context *gin.Context) {
	var articleUpdateDTO dto.ArticleUpdateDTO
	errDTO := context.ShouldBind(&articleUpdateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.articleService.IsAllowedToEdit(userID, articleUpdateDTO.ID) {
		id, errID := strconv.ParseUint(userID, 10, 64)
		if errID == nil {
			articleUpdateDTO.UserID = id
		}
		result := c.articleService.Update(articleUpdateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("You don not have a permission to update", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *articleController) Delete(context *gin.Context) {
	var article models.Article
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed. ID is not correct", "No paramater ID were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	article.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.articleService.IsAllowedToEdit(userID, article.ID) {
		c.articleService.Delete(article)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You don not have a permission to delete", "You are not the owner of the article", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *articleController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
