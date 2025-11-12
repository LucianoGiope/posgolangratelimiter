package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/app"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/getenv"
	"github.com/LucianoGiope/posgolangdesafioFinalRateLimiter/config/redisdb"
	"github.com/go-chi/chi/v5"
)

func main() {
	if err := getenv.NewConfig("../../.env", ".env"); err != nil {
		log.Fatal("Não foi possível ler as configurações: ", err)
	}

	rdb := redisdb.NewRedisClient(getenv.AppConfig.RedisAddr)
	defer rdb.Close()

	router := chi.NewRouter()

	if err := app.StartDependencies(router, rdb); err != nil {
		log.Fatal("Falha ao iniciar o sistema: ", err)
	}

	fmt.Printf("O servidor está rodano na porta: %s\n", getenv.AppConfig.ServerPortDefault)

	err := http.ListenAndServe(getenv.AppConfig.ServerPortDefault, router)
	if err != nil {
		log.Fatal(err)
	}
}
