package mysql

import "com.bookstore/demo/models"

func GetShelf(id int64) (*models.Shelf, error) {
	sqlStr := `select * from shelves where id=?`
	data := new(models.Shelf)
	err := db.Get(data, sqlStr, id)
	return data, err
}
