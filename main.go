package main

import (
	"context"
	"fmt"
	gomuxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wmt-challenge/handlers"
	"wmt-challenge/util"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind Address for Server")

func main() {
	env.Parse()

	l := log.New(os.Stdout, "wmt-challenge-api", log.LstdFlags)
	serverMux := mux.NewRouter()

	muxHandlers := gomuxhandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	muxMethods := gomuxhandlers.AllowedMethods([]string{"GET", "HEAD", "OPTIONS"})
	muxOrigins := gomuxhandlers.AllowedOrigins([]string{"*"})
	cors := gomuxhandlers.CORS(muxHandlers, muxOrigins, muxMethods)

	productsHandler := handlers.NewProductsHandler(l)
	getRouter := serverMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/api/products/byId/{id}", productsHandler.GetProductById)
	getRouter.HandleFunc("/api/products/byToken/{token}", productsHandler.GetProductByToken)
	gomuxhandlers.AllowedOrigins([]string{"*"})

	// Create the server
	server := http.Server{
		Addr:         *bindAddress,
		Handler:      cors(serverMux),
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  180 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 9090")

		err := server.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// Trapping SIGTERM or interrupts to gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until signal is received
	sig := <-c
	log.Println("Got signal: ", sig)

	// Shutdown the server gracefully
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}

func TestPalindrome() {
	var s string
	for {
		print("Enter something: ")
		fmt.Scanf("%s", &s)
		isPalindrome := util.IsPalindrome(s)
		if isPalindrome {
			print("Yes, its a palindrome")
		} else {
			print("No, it is not a palindrome")
		}
	}
}
