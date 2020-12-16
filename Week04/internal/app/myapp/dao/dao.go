package dao

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/errcode"
	"github.com/schwarzeni/Go-000/Week04/internal/app/myapp/model"
	"gorm.io/gorm"
)

type Dao struct {
	engine *gorm.DB
}

func NewDao(engine *gorm.DB) (*Dao, error) {
	return &Dao{engine: engine}, nil
}

func (dao *Dao) GetArticles() (articles []*model.Article, err error) {
	result := dao.engine.Find(&articles)
	if result.Error != nil {
		return nil, errors.Wrapf(errcode.ErrDB, fmt.Sprintf("failed to get articles with err: %v", result.Error))
	}
	return
}
