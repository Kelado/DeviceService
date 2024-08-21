package main

import (
	"fmt"
	"net/http"

	"github.com/Kelado/DeviceService/controllers"
	"github.com/Kelado/DeviceService/repositories"
	"github.com/go-chi/chi/v5"
)

const (
	Addr = ":8000"
)

func main() {
	deviceRepo := repositories.NewSQLiteDeviceRepo(nil)
	controller := controllers.NewDeviceController(deviceRepo)
	router := chi.NewRouter()

	controller.InitRouter(router)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	fmt.Println("Service listening on", Addr)
	if err := http.ListenAndServe(Addr, router); err != nil {
		fmt.Println(err)
	}
}
