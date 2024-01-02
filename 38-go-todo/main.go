package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostName       string = "localhost:27017" // PORT for MongoDB
	dbName         string = "demo_todo"
	collectionName string = "todo" 
    port           string = ":9000"
)

type(

  // We have two similar structs because, one struct is working with BSON data
  // and another one is working with JSON data.
  todoModel struct{
    ID        bson.ObjectId `bson:"_id,omitempty"`
    Title     string        `bson:"title"`
    Completed bool          `bson:"completed"`
    CreatedAt time.Time     `bson:"createdAt"`
  }

  todo struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
  }
)

func init(){
  rnd = renderer.New()
  sess, err := mgo.Dial(hostName)
  checkErr(err)
  sess.SetMode(mgo.Monotonic, true)  // Modes in MongoDB - String, Monotonic and Eventual
  db = sess.DB(dbName)
}

func checkErr(err error){
  if err != nil {
    log.Fatal(err)
  }
}

func homeHandler(w http.ResponseWriter, r * http.Request){
  err := rnd.Template(w, http.StatusOK, []string{"static/home.tpl"}, nil)
  checkErr(err)
}

func todoHandlers() http.Handler {
  rg := chi.NewRouter();
  rg.Group(func(r chi.Router){
    r.Get("/", fetchTodos)
    r.Post("/", createTodo)
    r.Put("/{id}", updateTodo)
    r.Delete("/{id}", deleteTodo)
  })
  return rg
}

func fetchTodos(w http.ResponseWriter, r *http.Request){
  todos := []todoModel{}

  if err := db.C(collectionName).Find(bson.M{}).All(&todos); err != nil {
    rnd.JSON(w, http.StatusProcessing, renderer.M{
      "message":"Failed to fetch data",
      "error":err,
    })
  }

  todoList := []todo{}

  // We are ignoring the iterator and only taking in the value. Think of this
  // as an enumerated list of objects and we are ingoring the enumeration.
  // Also, we need to conver the BSON type to JSON hence, we are making 
  // this loop and adding things one by one
  for _, t := range todos{
    todoList = append(todoList, todo{
      ID: t.ID.Hex(),
      Title: t.Title,
      Completed: t.Completed,
      CreatedAt: t.CreatedAt,
    })
  }
  rnd.JSON(w, http.StatusOK, renderer.M{
    "data" : todoList,
  })
}

func createTodo(w http.ResponseWriter, r *http.Request){
  var t todo

  if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
    rnd.JSON(w, http.StatusProcessing, err)
    return
  }

  // simple validation
  if t.Title == "" {
    rnd.JSON(w, http.StatusBadRequest, renderer.M{
      "message" : "The title field is required",
    })
    return
  }

  // if input is okay, create a todo
  tm := todoModel{
    ID: bson.NewObjectId(),
    Title: t.Title,
    Completed: false,
    CreatedAt: time.Now(),
  }
  if err := db.C(collectionName).Insert(&tm); err != nil {
    rnd.JSON(w, http.StatusProcessing, renderer.M{
      "message" : "Failed to save todo",
      "error" : err,
    })
    return
  }

  rnd.JSON(w, http.StatusCreated, renderer.M{
    "message" : "Todo created successfully",
    "todo_id" : tm.ID.Hex(),
  })
}

func updateTodo(w http.ResponseWriter, r *http.Request){
  id :=  strings.TrimSpace(chi.URLParam(r, "id"))

  if !bson.IsObjectIdHex(id){
    rnd.JSON(w, http.StatusBadRequest, renderer.M{
      "message" : "This id is invalid",
    })
    return
  }

  var t todo

  if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
    rnd.JSON(w, http.StatusProcessing, err)
    return
  }

  // simple validation
  if t.Title == ""{
    rnd.JSON(w, http.StatusBadRequest, renderer.M{
      "message" : "The title field is required",
    })
    return
  }

  // If input is okay, update a todo
  if err := db.C(collectionName).
    Update(
    bson.M{"_id":bson.ObjectIdHex(id)},
    bson.M{"title": t.Title, "completed" : t.Completed},
    ); err != nil {
    rnd.JSON(w, http.StatusProcessing, renderer.M{
      "message" : "Failed to update todo",
      "error" : err,
    })
    return
  }

  rnd.JSON(w, http.StatusOK, renderer.M{
    "message" : "Todo updated successfully",
  })
}

func deleteTodo(w http.ResponseWriter, r *http.Request){
  id := strings.TrimSpace(chi.URLParam(r, "id"))

  if !bson.IsObjectIdHex(id) {
    rnd.JSON(w, http.StatusBadRequest, renderer.M{
      "message" : "The id is invalid",
    })
    return
  }

  if err := db.C(collectionName).RemoveId(bson.ObjectIdHex(id)); err != nil {
    rnd.JSON(w, http.StatusProcessing, renderer.M{
      "message" : "Failed to delete todo",
      "error" : err,
    })
  }

  rnd.JSON(w, http.StatusOK, renderer.M{
    "message" : "Todo successfully deleted",
  })
}

func main(){

  // Stopping the server gracefully
  stopChan := make(chan os.Signal)
  signal.Notify(stopChan, os.Interrupt)

  r := chi.NewRouter();
  r.Use(middleware.Logger)
  r.Get("/", homeHandler)
  r.Mount("/todo", todoHandlers())

  srv := &http.Server{
    Addr: port,
    Handler: r,
    ReadTimeout: 60 * time.Second,
    WriteTimeout: 60 * time.Second,
    IdleTimeout: 60 * time.Second,
  }

  go func(){ // goroutines
    log.Println("Listening on port", port)
    if err := srv.ListenAndServe(); err != nil {
      log.Printf("Listen:%s\n", err)
    }
  }()

  <-stopChan // Yet to figure out

  log.Println("Shutting down server...")
  ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
  srv.Shutdown(ctx)
  defer cancel()
  log.Println("Server gracefully stopped")
}
