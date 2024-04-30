package models

import(
    "github.com/jinzhu/gorm"
    "github.com/IAmRiteshKoushik/go-50/03-go-bookstore/pkg/config"
)

var db *gorm.DB

type Book struct{
    gorm.Model
    Name string `gorm:""json:"name"`
    Author string `json:"author"`
    Publication string `json:"publication"`
}

// Initializing the ORM
func init()  {
    config.Connect()
    db = config.GetDB()
    db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
    db.NewRecord(b)
    // Passing the values inside the book as b -> pointer
    db.Create(&b)
    return b
}

func GetAllBooks() []Book {
    var Books []Book
    db.Find(&Books)
    return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB){
    var getBook Book
    db := db.Where("ID=?", Id).Find(&getBook)
    return &getBook, db
}

func DeleteBook(Id int64) Book {
    var book Book
    db.Where("ID=?", Id).Delete(book)
    return book
}

// There is not update function for this project
// Approach : Find the book, delete it and then 
// add in the new data
