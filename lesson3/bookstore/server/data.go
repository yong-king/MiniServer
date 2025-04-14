package main

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shelf struct{
	ID int64 `gorm:"primaryKey"`
	Theme string
	Size int64
	CreateAt time.Time
	UpdataAt time.Time
}

type Book struct{
	ID int64 `gorm:"primaryKey"`
	Author string
	Title string
	ShelfID int64
	CreateAt time.Time
	UpdateAt time.Time
}

type bookstore struct{
	db *gorm.DB
}

func NewDB() (*gorm.DB, error) {
	dsn := "root:youngking98@tcp(127.0.0.1:3306)/bookstore?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect to mysql failed: %v\n", err)
	}

	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

const(
	defaultSzie = 100
)

// CreateShelf 创建书架
func (b *bookstore) CreateShelf(ctx context.Context, data Shelf) (*Shelf, error){
	if len(data.Theme) <= 0 {
		return nil, errors.New("invalid theme")
	}
	size := data.Size
	if size <= 0 {
		size = defaultSzie
	}
	v := Shelf{Theme: data.Theme, Size: size, CreateAt: time.Now(), UpdataAt: time.Now()}
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

// ListShelves 获取书架列表
func (b *bookstore) ListShelves (ctx context.Context) ([]*Shelf, error) {
	var vl []*Shelf
	err := b.db.WithContext(ctx).Find(&vl).Error
	return vl, err
}

// GetShelf 获取书架
func (b *bookstore) GetShelf (ctx context.Context, id int64) (*Shelf, error) {
	v :=  Shelf{}
	err := b.db.WithContext(ctx).First(&v, id).Error
	return &v, err
}

// DelateShelf 删除书架
func (b *bookstore) DelateShelf (ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
}

// CreateBook 创建书籍
func (b *bookstore) CreateBook(ctx context.Context, data Book) (*Book, error) {
	// 检查参数
	if len(data.Author) == 0{
		return nil, errors.New("invalid author")
	}
	if len(data.Title) == 0{
		return nil, errors.New("invalid title")
	}
	if data.ShelfID <= 0{
		return nil, errors.New("invalid shelfid")
	}
	// 处理数据
	book := Book{Author: data.Author, Title: data.Title, ShelfID: data.ShelfID, CreateAt: time.Now(), UpdateAt: time.Now()}
	err := b.db.WithContext(ctx).Create(&book).Error
	return &book, err
}

// ListBooks 根据书架id获取书架书籍列表
func (b *bookstore) ListBooks(ctx context.Context, sid int64, NextId string, PageSize int) ([]*Book, error) {
	// 检查参数
	if sid <= 0 {
		return nil, errors.New("invalid shelf id")
	}
	// 查询书架书籍
	var bl []*Book
	err := b.db.WithContext(ctx).Where("shelf_id = ? and id > ?", sid, NextId).Limit(PageSize).Order("id desc").Find(&bl).Error
	return bl, err
}

// GetBook 获取指定的书籍
func (b *bookstore) GetBook(ctx context.Context, id int64, sid int64) (*Book, error) {
	// 检查参数
	if id <= 0{
		return nil, errors.New("invalid book id")
	}
	if sid <= 0{
		return nil, errors.New("invalid shelf id")
	}
	// 查询书籍
	var book Book
	err := b.db.WithContext(ctx).Where("shelf_id = ? AND id = ?", sid, id).First(&book).Error
	return &book, err
}

// 删除书籍
func (b *bookstore) DelateBook (ctx context.Context, id int64, sid int64) error {
	// 检查参数
	if id <= 0{
		return errors.New("invalid book id")
	}
	if sid <= 0{
		return errors.New("invalid shelf id")
	}
	// 删除操作
	var book Book
	return b.db.WithContext(ctx).Where("shelf_id = ? AND id = ?", sid, id).Delete(&book).Error
}