package main

import (
	"github.com/saarzur123/task-management/backend/service"
	"github.com/saarzur123/task-management/backend/utils"
	"log"
	"net/http"
)

func main() {
	dbInstance, err := service.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	defer dbInstance.Close()

	router := utils.SetupRoutes(&service.TaskManager{DB: dbInstance})

	log.Fatal(http.ListenAndServe(":8080", router))
}
