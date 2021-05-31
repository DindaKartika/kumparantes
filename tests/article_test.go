package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"kumparantes/router"

	"github.com/gavv/httpexpect"
)

func TestGetArticles(t *testing.T) {
	// create http.Handler
	handler := router.New()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	// without query string
	obj := e.GET("/articles").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
	obj.Value("data").Array().Element(0).Object()

	// with p set
	obj = e.GET("/articles").WithQuery("p", "1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")

	// with filter author
	obj = e.GET("/articles").WithQuery("author", "Dinda").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
	obj.Value("data").Array().Element(0).Object().Value("author").String().Contains("Dinda")

	// with filter query
	obj = e.GET("/articles").WithQuery("query", "slightly").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("data")
	obj.Value("data").Array().Element(0).Object().Value("body").String().Contains("slightly")
}

func TestAddArticles(t *testing.T) {
	// create http.Handler
	handler := router.New()

	// run server using httptest
	server := httptest.NewServer(handler)

	defer server.Close()

	// create httpexpect instance
	e := httpexpect.New(t, server.URL)

	payload := make(map[string]interface{})

	// normal add new
	payload = map[string]interface{}{
		"author": "Author 2",
		"title":  "How to build a snowman",
		"body":   "This is how you build a snowman",
	}
	obj := e.POST("/articles").
		WithJSON(payload).
		Expect().
		Status(http.StatusCreated).JSON().Object()
	obj.ContainsKey("author").ValueEqual("author", "Author 2")
	obj.ContainsKey("title").ValueEqual("title", "How to build a snowman")
	obj.ContainsKey("body").ValueEqual("body", "This is how you build a snowman")

	// failed because author empty
	payload = map[string]interface{}{
		"title": "How to build a snowman",
		"body":  "This is how you build a snowman",
	}
	obj = e.POST("/articles").
		WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
	obj.Value("validation_errors").Object().Value("author").Array().Element(0).String().Equal("The author field is required")

	// failed because title empty
	payload = map[string]interface{}{
		"author": "Author 2",
		"body":   "This is how you build a snowman",
	}
	obj = e.POST("/articles").
		WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
	obj.Value("validation_errors").Object().Value("title").Array().Element(0).String().Equal("The title field is required")

	// failed because body empty
	payload = map[string]interface{}{
		"author": "Author 2",
		"title":  "How to build a snowman",
	}
	obj = e.POST("/articles").
		WithJSON(payload).
		Expect().
		Status(http.StatusUnprocessableEntity).JSON().Object()
	obj.Value("validation_errors").Object().Value("body").Array().Element(0).String().Equal("The body field is required")
}
