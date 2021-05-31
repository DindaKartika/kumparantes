package cache

import (
	"encoding/json"
	"errors"
	"kumparantes/models"
	"kumparantes/models/orm"
	"math"
	"sort"
	"strconv"
)

func SetArticle(article models.Article) {
	var key1 string
	key1 = "article:" + strconv.Itoa(article.ID)

	// pack.Products = nil
	Set(key1, article)
}

func ReloadArticles() {
	var articles []models.Article
	orm.FindAll(&articles)

	// Set partner cache.
	for _, article := range articles {
		Clear("article:" + strconv.Itoa(article.ID))
		SetArticle(article)
	}
}

func GetArticles(page int) (result interface{}, err error) {
	var articles []models.Article
	articleCache := GetAll("article:*")

	if len(articleCache) > 0 {
		//prepare pagination
		rp := 25
		total_rec := len(articleCache)
		var startNum, endNum int
		if page == 0 || page == 1 {
			page = 1
			startNum = 0
		} else {
			startNum = rp * (page - 1)
		}
		endNum = rp * page
		lastPage := int(math.Ceil(float64(total_rec) / float64(rp)))

		//unmarshal into struct and group by range of data
		for _, data := range articleCache {
			var article models.Article
			err = json.Unmarshal([]byte(data), &article)
			if err != nil {
				return result, errors.New("No article found")
			}

			if article.ID > startNum && article.ID <= endNum {
				articles = append(articles, article)
			}
		}

		//sort and paginate
		if err == nil {
			sort.SliceStable(articles, func(i, j int) bool {
				return articles[i].ID > articles[j].ID
			})
			result := orm.PaginationResponse{
				Total:       total_rec,
				PerPage:     rp,
				CurrentPage: page,
				LastPage:    lastPage,
				Data:        articles,
			}

			return result, nil
		}
	}

	return result, errors.New("No article found")
}
