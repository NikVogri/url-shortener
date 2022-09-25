package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NikVogri/url-shortener/database"
	"github.com/NikVogri/url-shortener/server"
)

const (
	PORT int = 3000
)

type AddUrlBody struct {
	Url      string `json:"url"`
	Duration int    `json:"duration"`
}

type AddUrlResponse struct {
	Url string `json:"url"`
}

var db = new(database.Db)

func main() {
	// connect to DB
	db = database.Connect(os.Getenv("DB_CONN_STR"))
	defer db.Close()

	// init server
	app := server.Create()

	// register handlers
	app.HandleFunc("/add", handleAddUrl)
	app.HandleFunc("/", handleRedirect)

	// assign port, if not provided in env. fallback to default port
	port := PORT

	if os.Getenv("PORT") != "" {
		v, _ := strconv.Atoi(os.Getenv("PORT"))
		port = v
	}

	// start listening to requests
	server.Listen(port, app)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")[1]
	ri, e := db.FindRecordById(id)

	if e != nil {

		if strings.Contains(e.Error(), "no rows in result set") {
			http.Error(w, "This route does not exist", http.StatusNotFound)
		} else {
			fmt.Println(e.Error())
			http.Error(w, "Something went wrong on our end", http.StatusInternalServerError)
		}

		return
	}

	expireTimestamp := ri.CreatedTimestamp + ri.Duration
	currTimestamp := int(time.Now().UnixMilli())

	if expireTimestamp < currTimestamp {
		http.Error(w, "This route has expired", http.StatusBadRequest)
		return
	}

	er := db.IncrementClick(id)

	if er != nil {
		fmt.Println(e.Error())
		http.Error(w, "Something went wrong on our end", http.StatusInternalServerError)
	}

	http.Redirect(w, r, ri.OriginalUrl, http.StatusTemporaryRedirect)
}

func handleAddUrl(w http.ResponseWriter, r *http.Request) {
	var body AddUrlBody
	e := json.NewDecoder(r.Body).Decode(&body)

	if e != nil {

		if e.Error() == "EOF" {
			http.Error(w, "JSON payload expected", http.StatusBadRequest)
		} else {
			fmt.Println(e.Error())
			http.Error(w, "Something went wrong on our end", http.StatusInternalServerError)
		}

		return
	}

	if body.Url == "" {
		http.Error(w, "Missing 'url' property in payload", http.StatusBadRequest)
		return
	}

	if body.Duration == 0 {
		body.Duration = 300000 // 5 minutes
	}

	newRand := rand.New(rand.NewSource(time.Now().UnixMicro()))
	id := strconv.Itoa(newRand.Int())

	er := db.AddRecord(&database.RecordItem{
		Id:               id,
		OriginalUrl:      string(body.Url),
		Clicks:           0,
		CreatedTimestamp: int(time.Now().UnixMilli()),
		Duration:         body.Duration,
	})

	if er != nil {
		fmt.Println(e.Error())
		http.Error(w, "Something went wrong on our end", http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(AddUrlResponse{
		Url: os.Getenv("URL") + id,
	})
}
