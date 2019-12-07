package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Device struct
type Device struct {
	Name string
	ID   string
}

// Phone looks better capitalized
var Phone Device

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func disrupt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		js, err := json.Marshal(Device{
			Name: Phone.Name,
			ID:   Phone.ID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			fmt.Printf("Error %q", err)
		}
		return
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.StatusText(200)
		fmt.Println(Phone)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", disrupt)
	http.HandleFunc("/favicon.ico", faviconHandler)
	fmt.Printf("Starting Disrupter\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
