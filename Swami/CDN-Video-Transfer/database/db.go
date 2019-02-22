package main

import (
	"context"
	"log"

	"firebase.google.com/go"
	"google.golang.org/api/option"
)

func main() {
	sa := option.WithCredentialsFile("./swami-database-firebase-adminsdk-6e01e-6d14ba6b69.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	m := make(map[string]string)
	m["Latency"] = "10"

	log.Print(m)
	result, err := client.Collection("sampleData").Doc("inspiration").Set(context.Background(), m)
	if err != nil {
		log.Fatalln(err)
	}
	log.Print(result)
	defer client.Close()

}
