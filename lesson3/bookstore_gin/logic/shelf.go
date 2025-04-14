package logic

import (
	"com.bookstore/demo/dao/mysql"
	"com.bookstore/demo/models"
	"go.uber.org/zap"
)

func GetShelf(id int64) (*models.Shelf, error) {
	// 根据书架id获取书架信息
	data, err := mysql.GetShelf(id)
	if err != nil {
		zap.L().Error("mysql.GetShelf failed", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return data, nil
}
