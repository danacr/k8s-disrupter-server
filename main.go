package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Device struct
type Device struct {
	Name string
	ID   string
}

type Payload struct {
	Target  Target  `json:"target"`
	Command Command `json:"command"`
}
type Target struct {
	Type  string `json:"type"`
	Hosts string `json:"hosts"`
	Exact int    `json:"exact"`
}
type Command struct {
	Type        string   `json:"type"`
	CommandType string   `json:"commandType"`
	Args        []string `json:"args"`
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

func disrupt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "K8s Disrupter Server")
		return
	case "POST":
		phone := Device{}
		err := json.NewDecoder(r.Body).Decode(&phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = rebootnode(os.Getenv("GREMLIN_TEAM_ID"), os.Getenv("GREMLIN_API_KEY"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.StatusText(200)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	var err error
	if os.Getenv("GREMLIN_TEAM_ID") == "" {
		err = errors.New("GREMLIN_TEAM_ID env var required")
		log.Fatal(err)
	}
	if os.Getenv("GREMLIN_API_KEY") == "" {
		err = errors.New("GREMLIN_API_KEY env var required")
		log.Fatal(err)
	}
	http.HandleFunc("/", disrupt)
	http.HandleFunc("/favicon.ico", faviconHandler)
	fmt.Printf("Starting Disrupter\n")
	if err = http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func rebootnode(teamid, apikey string) error {
	data := Payload{
		Target: Target{
			Type:  "Random",
			Hosts: "all",
			Exact: 1,
		},
		Command: Command{
			Type:        "shutdown",
			CommandType: "Shutdown",
			Args:        []string{"-r", "-d", "0"},
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.gremlin.com/v1/attacks/new?teamId="+teamid, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+apikey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	log.Panicln("Triggered Node Reboot")
	return nil
}
