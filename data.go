package main

import (
	"context"
	"errors"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	defaultShelfSize = 100
)

//使用GORM

func NewDB(dsn string) (*gorm.DB, error) {
	//mysql

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	//sqlite

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移 schema
	db.AutoMigrate(&Shelf{}, &Book{})
	return db, nil
}

//定义model

// Shelf 书架
type Shelf struct {
	ID       int64 `gorm:"primiaryKey"` //主键
	Theme    string
	Size     int64
	CreateAt time.Time
	UpdateAt time.Time
}

// Book 图书

type Book struct {
	ID       int64 `gorm:"primaryKey"`
	Author   string
	Titile   string
	ShelfID  int64
	CreateAt time.Time
	UpdateAt time.Time
}

// 数据库操作
type bookstore struct {
	db *gorm.DB
}

//实现增删改查

// 创建书架
func (b *bookstore) CreateShelf(ctx context.Context, data Shelf) (*Shelf, error) {
	//如果主题字符小于等于0，则不符合要求
	if len(data.Theme) <= 0 {
		return nil, errors.New("invalid theme")
	}
	//书架上图书规模默认为100
	size := data.Size
	if size <= 0 {
		size = defaultShelfSize
	}
	v := Shelf{
		Theme:    data.Theme,
		Size:     size,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	// 调用orm方法
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}

// GetShelf 根据id获取书架
func (b *bookstore) GetShelf(ctx context.Context, id int64) (*Shelf, error) {
	v := Shelf{}
	err := b.db.WithContext(ctx).First(&v, id).Error
	return &v, err
}

// ListShelves 获取书架列表
func (b *bookstore) ListShelves(ctx context.Context) ([]*Shelf, error) {
	var vl []*Shelf
	err := b.db.WithContext(ctx).Find(&vl).Error
	return vl, err
}

// DeleteShelf 删除书架
func (b *bookstore) DeleteShelf(ctx context.Context, id int64) error {
	return b.db.WithContext(ctx).Delete(&Shelf{}, id).Error
}

// ListBooks 返回某个书架上的图书列表
func (b *bookstore) ListBooks(ctx context.Context, id int64) ([]*Book, error) {
	var vl []*Book
	err := b.db.WithContext(ctx).Where("ShelfID = ?", id).Find(&vl).Error
	return vl, err
}

// CreateBook 创建书架上的图书
func (b *bookstore) CreateBook(ctx context.Context, data Book) (*Book, error) {
	//如果书籍没有title则不能创建
	if len(data.Titile) <= 0 {
		return nil, errors.New("invalid Titile")
	}
	v := Book{
		Author:   data.Author,
		Titile:   data.Titile,
		ShelfID:  data.ShelfID,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	//调用orm方法
	err := b.db.WithContext(ctx).Create(&v).Error
	return &v, err
}
