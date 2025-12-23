package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ochadipa/url-shorterner-project/internal/db"
	"github.com/ochadipa/url-shorterner-project/internal/handlers"
	"go.uber.org/zap"
)

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	db.NewDB()
	db.NewRedisClient()
}
func main() {
	defer db.StopDB()
	r := mux.NewRouter()

	r.HandleFunc("/{id}", handlers.ShorternUrlHandler).Methods("GET")
	r.HandleFunc("/url", handlers.CreteSorternUrlHandler).Methods("POST")
	r.HandleFunc("/{id}", handlers.DeleteHandler).Methods("DELETE")

	zap.L().Info("Server starting", zap.String("port", ":8000"))
	err := http.ListenAndServe(":8000",r)
	zap.L().Fatal("failed, run serve", zap.Error(err))

}
