package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    log "github.com/go-kit/kit/log"
    "math/rand"
    "net/http"
    "os"
    "strconv"
    "time"
)

type Post struct {
    ID       string    `json:"id"`
    Title    string    `json:"title,omitempty"`
    Body     string    `json:"body,omitempty"`
    Author   Author    `json:"author"`
    CreateAt time.Time `json:"create_at"`
}
type Author struct {
    Name  string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
}

var posts []Post
var logger log.Logger


func main() {

    setInitData()
    r := mux.NewRouter()
    r.HandleFunc("/", GetAllPosts).Methods("GET")
    r.HandleFunc("/post", CreatePost).Methods("POST")
    r.HandleFunc("/post", QueryPost).Methods("GET")
    r.HandleFunc("/post/{id}", UpdatePost).Methods("PUT")
    r.HandleFunc("/post/{id}", DeletePost).Methods("DELETE")
    http.Handle("/", r)
    if err := http.ListenAndServe(":3000", r); err != nil {
        logger.Log("status", "fatal", "err", err)
        os.Exit(1)
    }
}
func GetAllPosts(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    var post Post
    _ = json.NewDecoder(r.Body).Decode(&post)
    post.ID = strconv.Itoa(rand.Intn(1000000))
    post.CreateAt = time.Now().UTC()
    posts = append(posts, post)
    json.NewEncoder(w).Encode(&post)
}

func QueryPost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    keys, ok := r.URL.Query()["id"]
    if !ok || len(keys[0]) < 1 {
        logger.Log("Url Param 'id' is missing")
        return
    }

    for _, item := range posts {
        if item.ID == keys[0] {
            json.NewEncoder(w).Encode(item)
            break
        }
    }
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", "application/json")
    params := mux.Vars(r)
    for index, item := range posts {
        if item.ID == params["id"] {
            posts = append(posts[:index], posts[index+1:]...)
            var post Post
            _ = json.NewDecoder(r.Body).Decode(&post)
            post.ID = params["id"]
            posts = append(posts, post)
            json.NewEncoder(w).Encode(post)
            return
        }
    }

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range posts {
        if item.ID == params["id"] {
            posts = append(posts[:index], posts[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(posts)
}

func setInitData() {
    posts = append(posts, Post{
        ID:       "1",
        Title:    "First Post",
        Body:     "Body 1",
        CreateAt: time.Now().UTC(),
        Author: Author{
            Name:  "John Doe",
            Email: "john@mail.com",
        },
    })

    posts = append(posts, Post{
        ID:       "2",
        Title:    "Second Post",
        Body:     "Body 2",
        CreateAt: time.Now().UTC(),
        Author: Author{
            Name:  "John Doe",
            Email: "john@mail.com",
        },
    })
}
