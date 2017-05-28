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

	client = twitch.NewClient(clientKey, secret, "http://localhost:8080/authorized")
	r := chi.NewRouter()

	r.Get("/", client.Authorize("openid", "user_read", "user_subscriptions", "user_follows_edit", "user_blocks_edit"))
	r.Get("/authorized", client.HandleAuthorization(handleAccess))
	r.Get("/user", handleGetUser)
	r.Get("/test", handleTest)

	log.Println("Starting server...")
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

func handleTest(w http.ResponseWriter, r *http.Request) {
	scope := []string{
		"openid",
		"user_read",
		"user_subscriptions",
		"user_follows_edit",
		"user_blocks_edit",
	}

	access := twitch.NewAccess("m724mn6yvu28rx61kxdkmedomsuu75", scope)
	users, err := client.GetUsersByName("TehDotDev", "TehDot")
	if err != nil {
		log.Printf("Failed to get users: %s", err)
	}

	userID := users[0].ID
	blockID := users[1].ID

	block, err := client.WithAccess(access).BlockUser(userID, blockID)
	if err != nil {
		log.Printf("Failed to block: %s", err)
	}

	data, err := json.Marshal(&block)
	if err != nil {
		log.Printf("Failed to marshal json: %s")
	}

	w.Write(data)
}
