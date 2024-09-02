package main

import (
	"DiaSync/config"
	"DiaSync/server"
	"DiaSync/utils"
)

func main() {
	cfg := config.Init()

	utils.Init(cfg)

	storage := server.InitStorage(cfg.Db)
	router := server.InitRouter(storage)
	httpServer := server.InitHttpServer(cfg, router)

	go storage.Clear()

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
