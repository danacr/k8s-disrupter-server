package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
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
		// err = rebootnode(os.Getenv("GREMLIN_TEAM_ID"), os.Getenv("GREMLIN_API_KEY"))
		instancelist, err := getinstances()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var goodbye string
		for _, instanceWithNamedPorts := range instancelist {
			if instanceWithNamedPorts.Status == "RUNNING" {
				goodbye = instanceWithNamedPorts.Instance
			}
		}
		goodbye = goodbye[strings.LastIndex(goodbye, "/")+1:]

		result, err := deleteinstance(goodbye)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(result)

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
	if err = serviceAccount(); err != nil {
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

// serviceAccount shows how to use a service account to authenticate.
func serviceAccount() error {
	// Download service account key per https://cloud.google.com/docs/authentication/production.
	// Set environment variable GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json
	// This environment variable will be automatically picked up by the client.
	client, err := pubsub.NewClient(context.Background(), "your-project-id")
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	// Use the authenticated client.
	_ = client

	return nil
}

func getinstances() ([]*compute.InstanceWithNamedPorts, error) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	computeService, err := compute.New(c)
	if err != nil {
		return nil, err
	}

	// Project ID for this request.
	project := "kill-the-cluster"

	// The name of the zone where the instance group is located.
	zone := "europe-west6-a" //

	// The name of the instance group from which you want to generate a list of included instances.
	instanceGroup := "gke-chaotic-cluster-preemptible-pool-ae9c7dfb-grp"

	rb := &compute.InstanceGroupsListInstancesRequest{}
	var instancelist []*compute.InstanceWithNamedPorts
	req := computeService.InstanceGroups.ListInstances(project, zone, instanceGroup, rb)
	if err := req.Pages(ctx, func(page *compute.InstanceGroupsListInstances) error {
		instancelist = append(instancelist, page.Items...)
		return nil
	}); err != nil {
		return nil, err
	}
	return instancelist, nil
}

func deleteinstance(instance string) (string, error) {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		return "", err
	}

	// Project ID for this request.
	project := "kill-the-cluster"

	// The name of the zone where the instance group is located.
	zone := "europe-west6-a" //

	resp, err := computeService.Instances.Delete(project, zone, instance).Context(ctx).Do()
	if err != nil {
		return "", err
	}

	// TODO: Change code below to process the `resp` object:
	result := fmt.Sprintf("%#v\n", resp)
	return result, nil
}
