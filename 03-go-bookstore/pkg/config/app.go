package config

import(
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

// Initializing a global variable of the gorm.DB type which is going to hold
// a pointer to the database connection after it has been initialized
var db *gorm.DB

func Connect() {
    // Alter the connection string later
    d, err := gorm.Open("mysql", "user:password/tableName?charset=utf8&parseTime=True&loc=Local")
    if err != nil {
        panic(err)
    }
    db = d
}

func GetDB() *gorm.DB {
    // Getter function for the database
    return db
}


