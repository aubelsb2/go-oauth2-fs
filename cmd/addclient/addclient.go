package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"flag"
	go_oauth2_fs "github.com/officeadminsorted/form0-proxy-auth-rules/oauth/go-oauth2-fs"
	"github.com/officeadminsorted/oauth2/v5/models"
	"github.com/segmentio/ksuid"
	"log"
)

var (
	userId  = flag.String("user-id", "", "The desired user id (not used)")
	domain  = flag.String("domain", "", "The domain ie www.example.com OR a URL http://www.example.com/path")
	dataDir = flag.String("data-dir", "./data", "Directory to store the data")
)

func main() {
	flag.Parse()

	if len(*dataDir) == 0 {
		log.Panicf("Please specify a data dir")
	}

	if len(*domain) == 0 {
		log.Panicf("Please specify a client domain")
	}

	secretB := make([]byte, 32)
	if _, err := rand.Read(secretB); err != nil {
		log.Panic(err)
	}
	secret := base64.StdEncoding.EncodeToString(secretB)

	var client *go_oauth2_fs.Client = &go_oauth2_fs.Client{
		Client: models.Client{
			ID:     ksuid.New().String(),
			Secret: secret,
			Domain: *domain,
			UserID: *userId,
		},
		Extra: nil,
	}

	clientStore := go_oauth2_fs.NewClientStore(*dataDir).(*go_oauth2_fs.ClientStore)
	ctx := context.Background()
	if err := clientStore.AddClient(ctx, client); err != nil {
		log.Panicf("Error: %v", err)
	}

	log.Printf("Created")
	log.Printf("ID: %v", client.ID)
	log.Printf("Domain: %v", client.Domain)
	log.Printf("Secret: %v", client.Secret)
}
