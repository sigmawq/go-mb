package main

import (
	"net/http"
	"os"
	"fmt"	
	"github.com/gorilla/mux"
	"log"
	"io/ioutil"
)

type Channel struct {
	Name string
	Messages []string
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Please specify host and port")
		return
	}
	host := args[0] 
	port := args[1] 
	fmt.Println("Host: ", host)
	fmt.Println("Port: ", port)
	r := mux.NewRouter()
	var channels []Channel
	r.HandleFunc("/add/{channel}", func(w http.ResponseWriter, r *http.Request){
		log.Println("add")
		channel := mux.Vars(r)["channel"]
		if channel == "" {
			w.WriteHeader(http.StatusNotFound)
			return 
		}

		found := false
		for _, existing_channel := range channels {
			if existing_channel.Name == channel {
				found = true
				break
			}
		}

		if !found {
			channels = append(channels, Channel { Name: channel })
			log.Printf("Channel %v created\n", channel)
		}

		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/push/{channel}", func(w http.ResponseWriter, r *http.Request){
		log.Println("push")
		_channel := mux.Vars(r)["channel"]
		if _channel == "" {
			log.Println("Channel name was empty")
			w.WriteHeader(http.StatusNotFound)
			return 
		}

		found := false
		var channel *Channel
		for i, _ := range channels {
			existing_channel := &channels[i]
			if existing_channel.Name == _channel {
				channel = existing_channel
				found = true
				break
			}
		}

		if !found {
			log.Println("Channel was not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		message, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error while reading message from user: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return	
		}
		channel.Messages = append(channel.Messages, string(message))
		log.Printf("Message (%v) length %v appended to channel %v. Messages total: %v", message, len(message), channel.Name, len(channel.Messages))
		w.WriteHeader(http.StatusOK)
	})

	r.HandleFunc("/pop/{channel}", func(w http.ResponseWriter, r *http.Request){
		log.Println("pop")
		_channel := mux.Vars(r)["channel"]
		if _channel == "" {
			w.WriteHeader(http.StatusNotFound)
			return 
		}

		found := false
		var channel *Channel
		for i, _ := range channels {
			existing_channel := &channels[i]
			if existing_channel.Name == _channel {
				channel = existing_channel
				found = true
				break
			}
		}

		if !found {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if len(channel.Messages) == 0 {
			w.WriteHeader(http.StatusNoContent)			
			return
		}

		first := channel.Messages[0]
		channel.Messages = channel.Messages[1:len(channel.Messages)]
		w.Write([]byte(first))
	})

	log.Println("Queue initiailize")
	http.ListenAndServe(host + ":" + port, r)
}