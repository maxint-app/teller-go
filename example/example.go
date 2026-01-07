package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/maxint-app/teller-go"
)

func main() {
	client, err := teller.NewClient(
		"./example/certs/certificate.pem",
		"./example/certs/private_key.pem",
		nil, // optional access token
	)
	if err != nil {
		log.Fatal(err)
	}

	institutions, err := client.Institutions.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, institution := range institutions {
		fmt.Printf("Institution: %s (ID: %s)\n", institution.Name, institution.ID)
		fmt.Println("Products:", strings.Join(institution.Products, ", "))
	}

	RunWebhookServer()
}

// This example spins up a minimal HTTP server that verifies and parses
// Teller webhook events using ConstructWebhook.
// You can use ngrok or a similar tool to expose this server to the internet
// for testing with real Teller webhooks.
func RunWebhookServer() {
	// Replace with your active signing secrets from Teller dashboard.
	signingSecrets := []string{"your_signing_secret_1", "your_signing_secret_2"}

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("read body error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sigHeader := r.Header.Get("Teller-Signature")

		event, err := teller.ConstructWebhook(body, sigHeader, signingSecrets)
		if err != nil {
			log.Printf("signature verify/parse error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("webhook verified: type=%s id=%s", event.Type, event.ID)
		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
