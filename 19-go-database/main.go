package main

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "github.com/blend/go-sdk/stringutil"
    "github.com/jcelliott/lumber" // A simple logger for Go
)

const Version = "1.0.0"

type(
    Logger interface {
        Error(string, ...interface{})
        Warn(string, ...interface{})
        Info(string, ...interface{})
        Debug(string, ...interface{})
        Trace(string, ...interface{})
    }
    Driver struct {
        mutex sync.Mutex
        mutexes map[string]*sync.Mutex
        dir string
        log Logger
    }
)

type Address struct {
    City string
    State string
    Country string
    Pincode json.Number
}

type User struct {
    Name string
    Age json.Number 
    Contact string
    Company string
    Address Address
}

type Options struct {
    Logger
}

func New()(){

}

func (d *Driver) Write() error {

}

func (d *Driver) Read() error {

}

func (d *Driver) ReadAll()() {

}

func (d *Driver) Delete() error {

}

func (d *Driver) getOrCreateMutex() *sync.Mutex {

}

func main(){
    dir := "./"
    db, err := new(dir, nil)
    if err != nil {
        fmt.Println("Error", err)
    }

    // Sample data
    employees := []User{
        {"John", "23", "23344333", "Myrl Tech", Address{"bangalore", "Karnataka", "india", "410013"}},
        {"Paul", "25", "23344333", "Google", Address{"sf", "california", "india", "410013"}},
        {"Robert", "27", "23344333", "Microsoft", Address{"bangalore", "Karnataka", "india", "410013"}},
        {"Vince", "29", "23344333", "Facebook", Address{"bangalore", "Karnataka", "india", "410013"}},
        {"Neo", "31", "23344333", "Remote-Teams", Address{"bangalore", "Karnataka", "india", "410013"}},
        {"Albert", "32", "23344333", "Dominate.ai", Address{"bangalore", "Karnataka", "india", "410013"}},
    }

    // Writing data to the database (JSON files)
    for _, value := range employees{
        db.Write("users", value.Name, User{
            Name: value.Name,
            Age: value.Age,
            Contact: value.Contact,
            Company: value.Company,
            Address: value.Address,
        })
    }

    // Reading all users
    records, err := db.ReadAll("users")
    if err != nil {
        fmt.Println("Error", err)
    }
    fmt.Println(records)

    // The problem is that all the records would be in JSON format and 
    // that will be printed. To be able to use / manipulate the data 
    // in Go we need to unmarshal it.
    allusers := []User{}
    for _, f := range records {
        employeeFound := User{} // empty struct to capture the data
        if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
            fmt.Println("Error", err)
        }
        allusers = append(allusers, employeeFound)
    }
    fmt.Println(allusers)

    // if err := db.Delete("user", "john"); err != nil {
    //     fmt.Println("Error", err)
    // }
    // if err := db.Delete("user", "john"); err != nil {
    //     fmt.Println("Error", err)
    // }
}
