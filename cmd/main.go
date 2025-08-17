package main

import (
	"fmt"
	"net/http"

	"github.com/GrigoDev/linker/configs"
	"github.com/GrigoDev/linker/internal/auth"
	"github.com/GrigoDev/linker/internal/link"
	"github.com/GrigoDev/linker/pkg/db"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	router := http.NewServeMux()

	// Hanlder
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
