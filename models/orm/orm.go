package orm

import (
	"fmt"
	"kumparantes/databases"
	"math"
	"reflect"

	"github.com/jinzhu/gorm"
)

type (
	DBFunc func(tx *gorm.DB) error // func type which accept *gorm.DB and return error

	PaginationResponse struct {
		Total       int         `json:"total"`
		PerPage     int         `json:"per_page"`
		CurrentPage int         `json:"current_page"`
		LastPage    int         `json:"last_page"`
		Data        interface{} `json:"data"`
	}
)

var (
	total_rec int
)

// Create
// Helper function to insert gorm model to database by using 'WithinTransaction'
func Create(v interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if !databases.App.DBConfig.NewRecord(v) {
			return err
		}
		if err = tx.Create(v).Error; err != nil {
			tx.Rollback() // rollback
			return err
		}
		return err
	})
}

// Save
// Helper function to save gorm model to database by using 'WithinTransaction'
// func Save(v interface{}) error {
// 	return WithinTransaction(func(tx *gorm.DB) (err error) {
// 		// check new object
// 		if databases.App.DBConfig.NewRecord(v) {
// 			return err
// 		}
// 		if err = tx.Save(v).Error; err != nil {
// 			tx.Rollback() // rollback
// 			return err
// 		}
// 		return err
// 	})
// }

// FindOneByID
// Helper function to find a record by using 'WithinTransaction'
func FindOneByID(v interface{}, id int) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Last(v, id).Error; err != nil {
			tx.Rollback() // rollback db transaction
			return err
		}
		return err
	})
}

// FindAll
// Helper function to find records by using 'WithinTransaction'
func FindAll(v interface{}) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Find(v).Error; err != nil {
			tx.Rollback() // rollback db transaction
			return err
		}
		return err
	})
}

// FindOneByQuery
// Helper function to find a record by using 'WithinTransaction'
func FindOneByQuery(v interface{}, params map[string]interface{}) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Where(params).Last(v).Error; err != nil {
			tx.Rollback() // rollback db transaction
			return err
		}
		return err
	})
}

// // FindByQuery
// // Helper function to find records by using 'WithinTransaction'
// func FindByQuery(v interface{}, params map[string]interface{}) (err error) {
// 	return WithinTransaction(func(tx *gorm.DB) error {
// 		if err = tx.Where(params).Find(v).Error; err != nil {
// 			tx.Rollback() // rollback db transaction
// 			return err
// 		}
// 		return err
// 	})
// }

// FindAllWithPage
// Helper function to find all records in pagination by using 'WithinTransaction'
// v interface{}	Gorm model struct
// page int	Page number
// rp int	Record per page to be showed
// filters int	Gorm model struct for filters
func FindAllWithPage(v interface{}, page int, filters interface{}) (resp PaginationResponse, err error) {
	var (
		offset   int
		lastPage int = 1
	)

	rp := 25

	// tx := databases.App.DBConfig.Begin()
	tx := databases.App.DBConfig

	// loop through filters
	refOf := reflect.ValueOf(filters).Elem()
	typeOf := refOf.Type()
	for i := 0; i < refOf.NumField(); i++ {
		f := refOf.Field(i)
		// ignore if empty
		// just make sure ModelFilterable its all in string type
		if f.Interface() != "" {
			valueString := fmt.Sprintf("%s", f.Interface())
			if typeOf.Field(i).Name == "Query" {
				tx = tx.Where("title LIKE ? or body LIKE ?", "%"+valueString+"%", "%"+valueString+"%")
			} else {
				tx = tx.Where(fmt.Sprintf("%s LIKE ?", typeOf.Field(i).Name), "%"+valueString+"%")
			}
		}
	}

	// copy of tx
	ctx := tx

	// get total record include filters
	ctx.Find(v).Count(&total_rec)

	if page == 0 {
		page = 1
	}

	offset = (page * rp) - rp

	lastPage = int(math.Ceil(float64(total_rec) / float64(rp)))

	tx.Limit(rp).Offset(offset).Order("created_at desc").Find(v)

	resp = PaginationResponse{
		Total:       total_rec,
		PerPage:     rp,
		CurrentPage: page,
		LastPage:    lastPage,
		Data:        &v,
	}
	if err != nil {
		// tx.Rollback() // rollback db transaction
		return resp, err
	}

	// tx.Commit()

	return resp, err
}

// WithinTransaction
// accept DBFunc as parameter
// call DBFunc function within transaction begin, and commit and return error from DBFunc
func WithinTransaction(fn DBFunc) (err error) {
	tx := databases.App.DBConfig.Begin() // start db transaction
	defer tx.Commit()
	err = fn(tx)
	// close db transaction
	return err
}
