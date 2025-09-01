package main

import (
	"fmt"
	"net/http"

	"github.com/GrigoDev/linker/configs"
	"github.com/GrigoDev/linker/internal/auth"
	"github.com/GrigoDev/linker/internal/link"
	"github.com/GrigoDev/linker/internal/user"
	"github.com/GrigoDev/linker/pkg/db"
	"github.com/GrigoDev/linker/pkg/middleware"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	LinkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Hanlder
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: LinkRepository,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
