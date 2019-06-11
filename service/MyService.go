package service

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/auth"

	firebase "firebase.google.com/go"

	"cloud.google.com/go/firestore"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type FirebaseService struct {
	firestore *firestore.Client
	firebase  *firebase.App
	app       *auth.Client
}

func NewFirebaseService(projectID, keyFile string) *FirebaseService {

	ctx := context.Background()
	log.Println("Connecting to: ", projectID)
	config := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(ctx, config)

	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	//	client, err := firebase.NewClient(ctx, projectID)
	if keyFile != "" {
		log.Println("Using keyfile  ", keyFile)
		opt := option.WithCredentialsFile(keyFile)
		app, err = firebase.NewApp(context.Background(), config, opt)
		if err != nil {
			log.Fatalf("error initializing app: %v\n", err)
		}
		//client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile(keyFile))
	}
	client, err := app.Auth(context.Background())
	// Get a Firestore client.

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)

	}

	// Close client when done.
	firestore, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}

	return &FirebaseService{firebase: app, firestore: firestore, app: client}
}

/**
* Retrieve profile by name
 */

func (fs FirebaseService) GetProfile(name string) bool {

	ctx := context.Background()

	iter := fs.firestore.Collection("profile").Where("name", "==", name).Documents(ctx)
	found := false
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		found = true
		fmt.Println(doc.Data())
	}
	return found
}
