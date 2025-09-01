package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/GrigoDev/linker/configs"
	"github.com/GrigoDev/linker/internal/auth"
	"github.com/GrigoDev/linker/internal/link"
	"github.com/GrigoDev/linker/internal/user"
	"github.com/GrigoDev/linker/pkg/db"
	"github.com/GrigoDev/linker/pkg/middleware"
)

func main() {
	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Done task")
	case <-ctxWithTimeout.Done():
		fmt.Println("Timeout")
	}
}

func main2() {
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
