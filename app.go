package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sy_backend/api"
	"sy_backend/config"
	"sy_backend/db"
	"syscall"
	"time"
)

func init() {
	if err := config.LoadConf(); err != nil {
		log.Fatal(err)
	}
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
}

func cleanup() {
	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	mux := &http.ServeMux{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("2006/01/02 15:04:05")
		w.Write([]byte(now))
	})

	mux.HandleFunc("/menu/list/{page}/{limit}", api.GetMenuList)
	mux.HandleFunc("/menu/random", api.GetRandomMenu)

	mux.HandleFunc("/menu", api.PostMenu)
	mux.HandleFunc("/menu/{id}", api.DeleteMenu)

	mux.HandleFunc("/menu/image/{id}", api.MenuImage)
	mux.HandleFunc("/menu/{id}/{field}", api.PatchMenu)

	srv := &http.Server{
		Addr:         config.Conf.ApiPort,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}
	quit := make(chan struct{})

	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGINT)
		<-interrupt

		cleanup()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
		close(quit)
	}()
	log.Printf("http server started on \033[32m%s\033[0m\n", config.Conf.ApiPort)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-quit
	println("\nserver closed")
}
