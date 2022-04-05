package main

import (
	"github.com/dentych/encrypt-service/http"
	"github.com/dentych/encrypt-service/storage"
	"log"
)

func main() {
	storage.Setup("./files")

	httpService := http.CreateHttp()
	err := httpService.Setup()
	if err != nil {
		log.Fatalf("failed to setup http service: %s", err)
	}

	err = httpService.Start()
	if err != nil {
		log.Fatalf("failed to run http service: %s", err)
	}
}
