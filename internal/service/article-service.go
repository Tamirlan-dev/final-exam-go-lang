package service

import (
	"fmt"

	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/dto"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/models"
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/repository"
	"github.com/mashingan/smapping"
)

type ArticleService interface {
	Insert(a dto.ArticleCreateDTO) models.Article
	Update(a dto.ArticleUpdateDTO) models.Article
	Delete(a models.Article)
	All() []models.Article
	FindByID(articleID uint64) models.Article
	IsAllowedToEdit(userID string, articleID uint64) bool
}

type articleService struct {
	articleRepository repository.ArticleRepository
}

func NewArticleService(articleRepo repository.ArticleRepository) ArticleService {
	return &articleService{
		articleRepository: articleRepo,
	}
}

func (service *articleService) Insert(a dto.ArticleCreateDTO) models.Article {
	article := models.Article{}
	err := smapping.FillStruct(&article, smapping.MapFields(&a))
	if err != nil {
		panic(err)
	}
	res := service.articleRepository.InsertArticle(article)
	return res
}

func (service *articleService) Update(a dto.ArticleUpdateDTO) models.Article {
	article := models.Article{}
	err := smapping.FillStruct(&article, smapping.MapFields(&a))
	if err != nil {
		panic(err)
	}
	res := service.articleRepository.UpdateArticle(article)
	return res
}

func (service *articleService) Delete(a models.Article) {
	service.articleRepository.DeleteArticle(a)
}

func (service *articleService) All() []models.Article {
	return service.articleRepository.AllArticles()
}

func (service *articleService) FindByID(articleID uint64) models.Article {
	return service.articleRepository.FindArticleByID(articleID)
}

func (service *articleService) IsAllowedToEdit(userID string, articleID uint64) bool {
	b := service.articleRepository.FindArticleByID(articleID)
	id := fmt.Sprintf("%v", b.UserID)
	return userID == id
}
