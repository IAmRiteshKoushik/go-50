package main

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "path/filepath"

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

func New(dir string, options *Options) (*Driver, error){
    dir = filepath.Clean(dir)
    opts := Options{}
    if options != nil {
        opts = *options
    }
    if opts.Logger == nil {
        opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
    }
    driver := Driver{
        dir: dir,
        mutexes: make(map[string]*sync.Mutex),
        log: opts.Logger,
    }

    // If database exists
    if _, err := os.Stat(dir); err == nil {
        opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
        return &driver, nil
    }
    // If database does not exist
    opts.Logger.Debug("Creating the database at '%s' ...\n", dir)
    return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(collection, resource string, v interface{}) error {
    // the collection is the directory
    // the resource is the file name
    if collection == ""{
        return fmt.Errorf("Missing collection - no place to save record!")
    }    
    if resource == "" {
        return fmt.Errorf("Missing resource - unable to save record (no name)!")
    }

    mutex := d.getOrCreateMutex(collection)
    mutex.Lock()
    defer mutex.Unlock()

    dir := filepath.Join(d.dir, collection)
    finalPath := filepath.Join(dir, resource + ".json")
    tmpPath := finalPath + ".tmp"

    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    b, err := json.MarshalIndent(v, "", "\t")
    if err != nil {
        return err
    }

    b = append(b, byte('\n'))
    if err := os.WriteFile(tmpPath, b, 0644); err != nil {
        return err
    }

    // After all activities are over, renaming the temporary path to final path
    return os.Rename(tmpPath, finalPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error {
    if collection == "" {
        return fmt.Errorf("Missing collection - no place to save record!")
    }
    if resource == "" {
        return fmt.Errorf("Missing resource - unable to save record (no name)!")
    }
    record := filepath.Join(d.dir, collection, resource)
    if _, err := stat(record); err != nil {
        return err
    }
    b, err := os.ReadFile(record + ".json")
    if err != nil {
        return err
    }
    return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string)([]string, error) {
    if collection == "" {
        return nil, fmt.Errorf("Missing collection - unable to read!")
    }
    dir := filepath.Join(d.dir, collection)
    if _, err := stat(dir); err != nil {
        return nil, err
    }
    files, _ := os.ReadDir(dir)
    var records []string 
    for _, file := range files{
        b, err := os.ReadFile(filepath.Join(dir, file.Name()))
        if err != nil {
            return nil, err
        }
        records = append(records, string(b))
    }
    return records, nil
}

func (d *Driver) Delete(collection, resource string) error {
    // The resource may be optional and can be an empty string
    // 1. If resource is passed then only that record is deleted
    // 2. If resource is empty string, then entire dir is deleted
    path := filepath.Join(collection, resource) 
    mutex := d.getOrCreateMutex(collection)
    mutex.Lock()
    defer mutex.Unlock()

    dir := filepath.Join(d.dir, path)

    switch fi, err := stat(dir); {
    case fi == nil, err != nil: 
        return fmt.Errorf("Unable to find files or directory named %v\n", path)
    case fi.Mode().IsDir():
        return os.RemoveAll(dir)
    case fi.Mode().IsRegular():
        return os.RemoveAll(dir + ".json")
    }
    return nil
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {

    d.mutex.Lock()
    defer d.mutex.Unlock()
    m, ok := d.mutexes[collection]
    if !ok {
        m =  &sync.Mutex{}
        d.mutexes[collection] = m
    }
    return m
}

func stat(path string)(fi os.FileInfo, err error){
    if _, err := os.Stat(path); os.IsNotExist(err){
        fi, err = os.Stat(path + ".json")
    }
    return fi, err
}

func main(){
    dir := "./"
    db, err := New(dir, nil)
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

    // Delete One Document
    // if err := db.Delete("users", "John"); err != nil {
    //     fmt.Println("Error", err)
    // }

    // Delete Entire Collection
    // if err := db.Delete("user", ""); err != nil {
    //     fmt.Println("Error", err)
    // }
}
