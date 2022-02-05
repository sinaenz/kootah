package main

import (
	"alibaba/shortener/handler"
	"alibaba/shortener/store/redis"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	// dependency injections
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	hr := handler.HttpHandler{
		Store:  redis.NewRedisStore(),
		Logger: logger,
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/save", hr.Save).Methods("POST")
	router.HandleFunc("/get-original", hr.GetOriginal).Methods("POST")
	router.HandleFunc("/get-info", hr.GetInfo).Methods("POST")
	router.HandleFunc("/sh/{short}", hr.Redirect)

	logger.Info("http server started!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
