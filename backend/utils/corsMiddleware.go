package utils

import (
	"github.com/gorilla/mux"
	"github.com/saarzur123/task-management/backend/handler"
	"github.com/saarzur123/task-management/backend/service"
	"net/http"
)

func SetupRoutes(taskRepository service.TaskRepository) *mux.Router {
	taskHandler := handler.TaskHandler{DB: taskRepository}

	router := mux.NewRouter()

	router.Use(corsMiddleware)

	router.HandleFunc("/tasks", taskHandler.CreateTask).Methods(http.MethodPost)
	router.HandleFunc("/tasks", taskHandler.GetAllTasks).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id:[0-9]+}", taskHandler.DeleteTask).Methods("DELETE")

	router.HandleFunc("/tasks", corsHandler).Methods("OPTIONS")
	router.HandleFunc("/tasks/{id:[0-9]+}", corsHandler).Methods("OPTIONS")

	return router
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusOK)
		return
	}
}
