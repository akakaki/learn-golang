package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"

  "github.com/gorilla/mux"
)

type Item struct {
  ID	string	`json:"id"`
  Name	string	`json:"name"`
  Object *Director  `json:"object"`
}

type Director struct {
  Type  string  `json:"type"`
}

var items []Item

func setItems(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  // params := mux.Vars(r)["id"]
  id := mux.Vars(r)["id"]

  switch r.Method {
    case "GET":
      if len(id) == 0 {
        json.NewEncoder(w).Encode(items)
      } else {
        for index, item := range items {
          if item.ID == id {
            json.NewEncoder(w).Encode(item)
            break
          } else if len(items) == index + 1 {
            var item Item
            json.NewEncoder(w).Encode(item)
          }
        }
      }
    case "POST":
      var item Item
      _ = json.NewDecoder(r.Body).Decode(&item)
      items = append(items, item)
      json.NewEncoder(w).Encode(item)
    case "PUT":
      for index, item := range items {
        if item.ID == id {
          items = append(items[:index], items[index+1:]...)
          var item Item
          _ = json.NewDecoder(r.Body).Decode(&item)
          items = append(items, item)
          json.NewEncoder(w).Encode(item)
        }
      }
    case "DELETE":
      for index, item := range items {
        if item.ID == id {
          items = append(items[:index], items[index+1:]...)
          break
        }
      }
      json.NewEncoder(w).Encode(items)
  }
  // fmt.Printf("Get Methods\n")
}

func main() {
  items = append(items, Item{ ID: "1", Name: "Item One", Object: &Director{ Type: "monster" } })
  items = append(items, Item{ ID: "2", Name: "Item Two", Object: &Director{ Type: "item" } })
  
  r := mux.NewRouter()

  r.HandleFunc("/items", setItems).Methods("GET")
  r.HandleFunc("/items/", setItems).Methods("GET")
  // curl http://localhost:8000/items

  r.HandleFunc("/items/{id}", setItems).Methods("GET")
  // curl http://localhost:8000/items/1

  r.HandleFunc("/items", setItems).Methods("POST")
  // curl -X POST -H "Content-Type: application/json" -d '{"id":"3","name":"Item Three","object":{"type": "effect"}}' http://localhost:8000/items

  r.HandleFunc("/items/{id}", setItems).Methods("PUT")
  // curl -X PUT -H "Content-Type: application/json" -d '{"id":"22","name":"Item Three","object":{"type": "effect"}}' http://localhost:8000/items/2

  r.HandleFunc("/items/{id}", setItems).Methods("DELETE")
  // curl -X DELETE http://localhost:8000/items/1

  fmt.Printf("Starting server at port 8000\n")
  log.Fatal(http.ListenAndServe(":8000", r))
  // http.ListenAndServe(":8000", r)
}

// go run restful-api.go
