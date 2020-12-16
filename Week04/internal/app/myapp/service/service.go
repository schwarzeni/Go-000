package service

import (
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/dao"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/dto"
)

type Service struct {
	dao *dao.Dao
}

func NewService(dao *dao.Dao) *Service {
	return &Service{dao: dao}
}

func (svc Service) GetArticles() (articles []*dto.Article, err error) {
	rawArticles, err := svc.dao.GetArticles()
	if err != nil {
		return nil, err
	}
	for _, rawArticle := range rawArticles {
		articles = append(articles, &dto.Article{
			ID:      rawArticle.ArticleID,
			Title:   rawArticle.Title,
			Content: rawArticle.Content,
		})
	}
	return
}
