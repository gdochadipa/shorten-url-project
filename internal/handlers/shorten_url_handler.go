package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ochadipa/url-shorterner-project/internal/model"
	"github.com/ochadipa/url-shorterner-project/internal/repositories"
	"github.com/ochadipa/url-shorterner-project/internal/service"
	"go.uber.org/zap"
)

func ShorternUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/favicon.ico" {
        // You can return 204 No Content (idiomatic) or 404
        w.WriteHeader(http.StatusNoContent)
        return
    }

	urlModel := &model.Url{}
	urlRepo := repositories.NewUrlRepo(urlModel)
	urlService := service.NewService(urlRepo)
	ctx := context.Background()

	urlModel, err := urlService.GetUrl(ctx, id)
	if err != nil {
		http.Error(w, "failed get url", http.StatusBadRequest)
		return
	}

	// http.RedirectHandler(urlModel.URL, http.StatusFound)
	http.Redirect(w,r,urlModel.URL,http.StatusFound)
	// fmt.Fprintf(w, `{"status":"success", "message":"get ID: %s"}`, urlModel.URL)
}

func CreteSorternUrlHandler(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	w.Header().Set("Content-Type", "application/json")

	urlModel := &model.Url{}
	urlRepo := repositories.NewUrlRepo(urlModel)
	urlService := service.NewService(urlRepo)
	ctx := context.Background()
	var uri string
	switch contentType {
	case "application/json":
		// Handle JSON payload
		var data map[string]any
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Could not decode JSON", http.StatusBadRequest)
			return
		}

		url, ok := data["url"]
		if !ok {
			http.Error(w, "url not found", http.StatusNotAcceptable)
		}
		uri = fmt.Sprint(url)

	case "application/x-www-form-urlencoded":
		// Handle Form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Could not parse form", http.StatusBadRequest)
			return
		}
		zap.L().Info("Received JSON  ", zap.Any("data", r.PostForm))

		url, ok := r.PostForm["url"]
		if !ok {
			http.Error(w, "url not found", http.StatusNotAcceptable)
		}
		uri = url[0]

	default:
		// Unsupported media type
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	urlModel, err := urlService.StoreUrl(ctx, uri)

	if err != nil {
		http.Error(w, "Could not process url", http.StatusBadRequest)
		return
	}

	parseToResponseUrlModel(w, urlModel)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success", "message":"Deleted item with ID: %s"}`, id)
}

func parseToResponseUrlModel(w http.ResponseWriter, url *model.Url) {
	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"status":  "success",
		"message": "URL shortened successfully",
		"data":    *url,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResponse)
}
