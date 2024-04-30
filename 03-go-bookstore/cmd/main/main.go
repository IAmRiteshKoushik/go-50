package main

import(
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/IAmRiteshKoushik/go-50/03-go-bookstore/pkg/routes"
)

func main(){
    // Getting a new router from gorilla mux
    r := mux.NewRouter()
    // passing the router to the routes function
    routes.RegisterBookStoreRoutes(r)
    // global handler made simple using net/http package
    http.Handle("/", r)
    // Logging if ListenAndServe fails
    log.Fatal(http.ListenAndServe("localhost:9000", r))
}
