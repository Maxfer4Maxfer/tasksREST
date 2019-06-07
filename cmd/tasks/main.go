package main

import (
	"tasksREST/pkg/config"
	"tasksREST/pkg/server"
	"tasksREST/pkg/service"
	storage "tasksREST/pkg/storage/gorm"
)

func main() {

	cfg := config.GetConfig()

	storage, sCloser := storage.New(cfg.DSN)
	defer sCloser()

	service := service.New(storage)
	server := server.New(cfg.HTTPAddr, service)

	server.Run()
}
