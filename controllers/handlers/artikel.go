package handlers

import (
	"kumparantes/models"
	"kumparantes/models/cache"
	"net/http"
	"strconv"

	"github.com/thedevsaddam/govalidator"

	"github.com/labstack/echo"
)

func AddArticles(c echo.Context) error {
	artikel := models.Article{}

	defer c.Request().Body.Close()

	rules := govalidator.MapData{
		"author": []string{"required"},
		"title":  []string{"required"},
		"body":   []string{"required"},
	}

	vld := ValidateRequest(c, rules, &artikel)
	if vld != nil {
		return c.JSON(http.StatusUnprocessableEntity, vld)
	}

	result, err := models.CreateArticle(&artikel)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to create new article")
	}

	cache.SetArticle(artikel)

	return c.JSON(http.StatusCreated, result)
}

func ValidateRequest(c echo.Context, rules govalidator.MapData, data interface{}) map[string]interface{} {
	var err map[string]interface{}

	opts := govalidator.Options{
		Request: c.Request(),
		Data:    data,
		Rules:   rules,
	}

	v := govalidator.New(opts)

	e := v.ValidateJSON()

	if len(e) > 0 {
		err = map[string]interface{}{"validation_errors": e}
	}

	return err
}

func GetArticles(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	author := c.QueryParam("author")
	query := c.QueryParam("query")

	var response interface{}
	defer c.Request().Body.Close()

	response, err = cache.GetArticles(page)
	if err != nil || author != "" || query != "" {
		if err != nil {
			cache.ReloadArticles()
		}
		//no article found in cache then return manually from db
		response, err = models.GetAllArticles(page, &models.ArticleFilter{author, query})
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, response)
}
