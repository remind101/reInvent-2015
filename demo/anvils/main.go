package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var queue *Queue

func init() {
	flag.Usage = usage
	queue = newQueue(os.Getenv("REDIS_URL"))
}

type DropRequest struct {
	Target string
}

func runWeb() {
	port := os.Getenv("PORT")
	fmt.Printf("Starting web server on %s\n", port)

	http.HandleFunc("/drop", func(w http.ResponseWriter, r *http.Request) {
		var dropRequest DropRequest
		if err := json.NewDecoder(r.Body).Decode(&dropRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		queue.Push(dropRequest)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Fprintln(os.Stdout, http.ListenAndServe(":"+port, nil))
}

func runWorker() {
	fmt.Println("Starting worker")
	for range time.Tick(1 * time.Second) {
		dropRequest := queue.Pop()
		fmt.Printf("Dropping anvil on %s...\n", dropRequest.Target)
	}
}

func main() {
	flag.Parse()

	if len(os.Args) <= 1 {
		usage()
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "web":
		runWeb()
	case "worker":
		runWorker()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `Usage of %s:

Commands:

  web  Run the web server
  worker Run the background worker
`, os.Args[0])
	flag.PrintDefaults()
}
