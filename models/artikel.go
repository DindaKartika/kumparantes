package models

import (
	"kumparantes/models/orm"
)

type (
	Article struct {
		BaseModel
		Author string `json:"author" sql:"type:text(5000)"`
		Title  string `json:"title" sql:"type:text(5000)"`
		Body   string `json:"body" sql:"type:text(5000)"`
	}

	// just use string type, since it will be use on query at DB layer
	ArticleFilter struct {
		Author string `json:"author"`
		Query  string `json:"query"`
	}
)

func CreateArticle(article *Article) (*Article, error) {
	var err error
	err = orm.Create(&article)
	return article, err
}

func GetArticle(id int) (Article, error) {
	var (
		artikel Article
		err     error
	)
	err = orm.FindOneByID(&artikel, id)
	return artikel, err
}

func GetAllArticles(page int, filters interface{}) (interface{}, error) {
	var (
		artikel []Article
		err     error
	)
	resp, err := orm.FindAllWithPage(&artikel, page, filters)
	return resp, err
}
