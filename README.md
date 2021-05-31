# KumparanTes
This simple web service is created for the completion of Kumparan Technical Assesment.

## Databases
All datas in this web service saved in two places:
- MySQL; used as the main database in this service
- Redis; used as cache 

## APIs
#### API Get All Articles
This endpoint used to get all articles saved in the database.
##### Endpoint
```sh
[GET] /articles
```
##### Filters
This API provided 3 filters in the query param :
- page
- author; will show articles with author field containing word from the filter
- query; will show articles with title or body field containing word from the filter

#### API Add New Article
This endpoint used to create new article and save it to the database and cache.
##### Endpoint
```sh
[POST] /articles
```
##### Request body
```sh
{
	"author": "Author Name",
	"title": "Article's Title",
	"body": "Article's body"
}
```

## Testing
In this service, there's testing provided. To do the testing, you can use command below:
```sh
go test ./... -coverprofile=coverage.out
```
To add up, there's manual testing using postman provided that can be accessed [here](https://docs.google.com/document/d/1rEKkEePFkf7BJjrdre9NTqL-dSihEN30xNkhl8--JFE/edit?usp=sharing)