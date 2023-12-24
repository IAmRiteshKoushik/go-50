package main

import(
  "fmt"
  "log"
  "net/http"
)

// Project Overview:
// 1. /      -> index.html
// 2. /hello -> hello function (returns JSON)
// 3. /form  -> form function  -> form.html

func formHandler(w http.ResponseWriter, r *http.Request){
  if err := r.ParseForm(); err != nil {
    fmt.Fprintf(w, "ParseForm() err : %v", err)
  }
  fmt.Fprintf(w, "POST request successful\n")
  name := r.FormValue("name")
  address := r.FormValue("address")
  fmt.Fprintf(w, "Name: %v\n", name)
  fmt.Fprintf(w, "Address: %v\n", address)

}


func helloHandler(w http.ResponseWriter, r *http.Request){
  // Guard clauses
  if r.URL.Path != "/hello" {
    http.Error(w, "404 not found", http.StatusNotFound)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method is not supported", http.StatusNotFound)
    return
  }

  fmt.Fprintf(w, "Hello!")
}

func main(){
  // Earlier the file path below was ./static 
  // But this is problematic because while this runs when - go run main.go
  // It fails to run as a build executable. That is becuase in the previous
  // method our run method relies on the working directory to locate the 
  // code whereas our build executable relies on the go.mod file to find and 
  // redirect to the appropriate directory. This can be achieved by using the 
  // go mod init <path> instead of the relative path and then accessing 
  // the folders where static files are kept.
  fileServer := http.FileServer(http.Dir("github.com/IAmRiteshKoushik/go-50/1-go-server/static"))

  // http.Handle accepts functions which satisfy an interface - http.Handler
  http.Handle("/", fileServer)

  // http.HandleFunc is a simple type that satisfied http.Handler
  // It automatically implements the http.ServeHTTP method from the package
  http.HandleFunc("/form", formHandler)
  http.HandleFunc("/hello", helloHandler)

  // Basically ::
  // http.Handle     -> An object responds to my http.request
  // http.HandleFunc -> A function responds to my http.request

  fmt.Println("Starting server at port 8000")

  // http.ListenAndServe takes in a port number and a handler. Here we 
  // are not providing a handler because we have different tasks for 
  // different paths which have already been handled previously.
  // Adding a handler here would be like adding a parent handler for everything
  if err := http.ListenAndServe(":8000", nil); err != nil {
    log.Fatal(err) // can use : fmt.Println("%v", err.Error())

    // Select between log and fmt using these facts :
    // 1. The log functions print to stderr by default and can be directed 
    //    to an arbitrary writer. The fmt.Printf function prints to stdout
    // 2. The log functions can print timestamp, source code location and 
    //    other information
    // 3. The log functions and fmt.Printf are both thread safe, but 
    //    concurrent writes by fmt.Printf above an OS dependent size can 
    //    be interleaved.

    // A logger can be used simulatenously from multiple go-routines; it 
    // guarentees to serialize access to the Writer.
  }
}
