package repository

import (
	"github.com/Omar-Temirgali/final-exam-go-lang/internal/domain/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	InsertArticle(a models.Article) models.Article
	UpdateArticle(a models.Article) models.Article
	DeleteArticle(a models.Article)
	AllArticles() []models.Article
	FindArticleByID(articleID uint64) models.Article
}

type articleConnection struct {
	connection *gorm.DB
}

func NewArticleRepository(dbConn *gorm.DB) ArticleRepository {
	return &articleConnection{
		connection: dbConn,
	}
}

func (db *articleConnection) InsertArticle(a models.Article) models.Article {
	db.connection.Save(&a)
	db.connection.Preload("User").Find(&a)
	return a
}

func (db *articleConnection) UpdateArticle(a models.Article) models.Article {
	db.connection.Save(&a)
	db.connection.Preload("User").Find(&a)
	return a
}

func (db *articleConnection) DeleteArticle(a models.Article) {
	db.connection.Delete(&a)
}

func (db *articleConnection) FindArticleByID(articleID uint64) models.Article {
	var article models.Article
	db.connection.Preload("User").Find(&article, articleID)
	return article
}

func (db *articleConnection) AllArticles() []models.Article {
	var articles []models.Article
	db.connection.Preload("User").Find(&articles)
	return articles
}
