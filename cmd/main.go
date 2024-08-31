package main

import (
	"DiaSync/config"
	"DiaSync/server"
	"DiaSync/utils"
	"fmt"
)

func main() {
	cfg := config.Init()

	fmt.Print(cfg)

	utils.Init(cfg)

	storage := server.InitStorage(cfg.Db)
	router := server.InitRouter(storage)
	httpServer := server.InitHttpServer(cfg, router)

	go storage.Clear()

	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}
