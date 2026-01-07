package main

import (
	"fmt"
	"log"
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
}
