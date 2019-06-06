package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jacoelho/fortuna"
	"github.com/jacoelho/fortuna/entropy/random"
	"github.com/jacoelho/fortuna/internal/api"
	"github.com/jacoelho/fortuna/internal/oracle"
)

func main() {
	f, err := fortuna.New()
	if err != nil {
		log.Fatal(err)
	}

	s := &random.Random{}
	m := fortuna.Manager{Pools: 32}

	m.Accumulate(f, s)

	seedFile := fortuna.NewSeedFile(f, "seed.dat")
	b, err := seedFile.Read()
	if err != nil {
		f.Seed(time.Now().UTC().UnixNano())
		log.Println(err)
	}

	f.SeedBytes(b)
	defer seedFile.Write()

	go func() {
		tickler := time.NewTicker(time.Minute)

		for {
			select {
			case <-tickler.C:
				if err := seedFile.Write(); err != nil {
					log.Println("seed file", err)
				}
			}
		}
	}()

	oracleService := oracle.New(rand.New(rand.Source(f)))

	router := http.NewServeMux()
	router.HandleFunc("/numbers", api.IntegersHandler(oracleService.Numbers))
	router.HandleFunc("/sequence", api.IntegersHandler(oracleService.Sequence))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen: %v\n", err)
	}

	<-done
}
