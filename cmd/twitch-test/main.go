package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/eriktate/go-twitch"
	"github.com/pressly/chi"
)

var client twitch.Client
var clientAccess twitch.Access

func main() {
	clientKey := os.Getenv("TWITCH_CLIENT")
	secret := os.Getenv("TWITCH_SECRET")

	log.Printf("%s\n%s", clientKey, secret)

	client = twitch.NewClient(clientKey, secret, "http://localhost:8080/authorized")
	r := chi.NewRouter()

	r.Get("/", client.Authorize("openid", "user_read"))
	r.Get("/authorized", client.HandleAuthorization(handleAccess))
	r.Get("/user", handleGetUser)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleAccess(access twitch.Access, err error) {
	if err != nil {
		log.Printf("Failed to get access: %s", err)
	}

	clientAccess = access

	log.Printf("Got access: %+v", access)
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	user, err := client.WithAccess(clientAccess).GetUser()
	if err != nil {
		log.Printf("Failed to get user: %s", err)
	}

	data, err := json.Marshal(&user)
	if err != nil {
		log.Printf("Failed to marshal response: %s", err)
	}

	w.Write(data)
}
