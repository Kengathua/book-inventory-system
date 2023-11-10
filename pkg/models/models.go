package models

import (
	"fmt"
	"time"

	"github.com/Kengathua/book-inventory-system/pkg/common"
	"github.com/Kengathua/book-inventory-system/pkg/common/variables"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Expired bool   `json:"expired"`
}

type User struct {
	common.BaseModel `gorm:"embedded"`
	FullName         string    `gorm:"column:full_name;not null" json:"full_name"`
	Email            string    `gorm:"unique;column:email;not null" json:"email"`
	Password         string    `gorm:"column:password;not null" json:"password"`
	UserType         string    `gorm:"column:user_type;not null" json:"user_type"`
	LastLogin        time.Time `gorm:"column:last_login" json:"last_login"`
}

func (u *User) GeneratePasswordHarsh() error {
	argon2ID := NewArgon2ID()
	passwordHash, err := argon2ID.Hash(u.Password)
	u.Password = passwordHash

	return err
}

func (user *User) GetUserID_UUID() (uuid.UUID, error) {
	return uuid.Parse(*user.ID)
}

func (u *User) CheckPasswordHarsh(password string) bool {
	argon2ID := NewArgon2ID()
	valid, err := argon2ID.Verify(password, u.Password)
	if err != nil {
		return false
	}

	return valid
}

func (user *User) CreateAuthor(tx *gorm.DB) error {
	userID, err := user.GetUserID_UUID()
	if err != nil {
		return err
	}
	author := Author{
		BaseModel: common.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName: user.FullName,
		Email:    user.Email,
		UserID:   *user.ID,
	}

	err = tx.Create(&author).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CreateLibrarian(tx *gorm.DB) error {
	userID, err := user.GetUserID_UUID()
	if err != nil {
		return err
	}
	librarian := Librarian{
		BaseModel: common.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName: user.FullName,
		Email:    user.Email,
		UserID:   *user.ID,
	}

	err = tx.Create(&librarian).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CreateStoreKeeper(tx *gorm.DB) error {
	userID, err := user.GetUserID_UUID()
	if err != nil {
		return err
	}
	storeKeeper := StoreKeeper{
		BaseModel: common.BaseModel{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		FullName: user.FullName,
		Email:    user.Email,
		UserID:   *user.ID,
	}

	err = tx.Create(&storeKeeper).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *User) CreateRespectiveUserType(tx *gorm.DB) error {
	switch user.UserType {
	case variables.AUTHOR:
		return user.CreateAuthor(tx)

	case variables.LIBRARIAN:
		return user.CreateLibrarian(tx)

	case variables.STOREKEEPER:
		return user.CreateStoreKeeper(tx)

	default:
		return fmt.Errorf("invalid user type")
	}
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	return user.CreateRespectiveUserType(tx)
}

type Author struct {
	common.BaseModel `gorm:"embedded"`
	FullName         string `gorm:"column:full_name" json:"full_name"`
	Email            string `gorm:"unique;column:email" json:"email"`
	UserID           string `gorm:"foreignKey:ID;column:user_id" json:"user_id"`
	User             User   `gorm:"foreignkey:UserID;references:id" json:"user"`
}

type StoreKeeper struct {
	common.BaseModel `gorm:"embedded"`
	FullName         string `gorm:"column:full_name" json:"full_name"`
	Email            string `gorm:"unique;column:email" json:"email"`
	UserID           string `gorm:"foreignKey:ID;column:user_id" json:"user_id"`
	User             User   `gorm:"foreignkey:UserID;references:id" json:"user"`
}

type Librarian struct {
	common.BaseModel `gorm:"embedded"`
	FullName         string `gorm:"column:full_name" json:"full_name"`
	Email            string `gorm:"unique;column:email" json:"email"`
	UserID           string `gorm:"foreignKey:ID;column:user_id" json:"user_id"`
	User             User   `gorm:"foreignkey:UserID;references:id" json:"user"`
}

type Book struct {
	common.BaseModel    `gorm:"embedded"`
	Title               string `gorm:"unique;column:title" json:"title"`
	AuthorID            string `gorm:"foreignKey:ID;column:author_id" json:"author_id"`
	Author              Author `gorm:"foreignkey:AuthorID;references:id" json:"author"`
	ReleaseYear         int    `gorm:"column:release_year" json:"release_year"`
	Genre               string `gorm:"column:genre" json:"genre"`
	ISBNNumber          string `gorm:"column:isbn_number" json:"isbn_number"`
	Status              string `gorm:"default:PENDING;column:status" json:"status"`
	Quantity            int64  `gorm:"column:quantity" json:"quantity"`
	IsLibrarianVerified bool   `gorm:"column:is_librarian_verified" json:"is_librarian_verified"`
	IsKeeperVerified    bool   `gorm:"column:is_keeper_verified" json:"is_keeper_verified"`
}

func (book *Book) ValidateUniqueBook(tx *gorm.DB) error {
	count := int64(0)
	err := tx.Table("books").Where("books.title = ?", book.Title).Not("books.id = ?", *book.ID).Count(&count).Error
	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("a book with a similar title already exists")
	}

	return nil
}

func (book *Book) ValidateISBNNumber(tx *gorm.DB) error {
	count := int64(0)
	if book.ISBNNumber == "" {
		return fmt.Errorf("please supply a valid ISBN number")
	}

	err := tx.Table("books").Where("books.isbn_number = ?", book.ISBNNumber).Not("books.id = ?", *book.ID).Count(&count).Error
	if err != nil {
		return err
	}

	if count != 0 {
		return fmt.Errorf("a book with a similar ISBN number already exists")
	}

	return nil
}

func (book *Book) BeforeCreate(tx *gorm.DB) error {
	err := book.BaseModel.BeforeCreate(tx)
	if err != nil {
		return err
	}

	err = book.ValidateUniqueBook(tx)
	if err != nil {
		return err
	}

	err = book.ValidateISBNNumber(tx)
	if err != nil {
		return err
	}

	return nil
}
