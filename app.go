package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sy/api"
	"sy/config"
	"sy/db"
	"syscall"
	"time"
)

func init() {
	if err := config.LoadConf(); err != nil {
		log.Fatal(err)
	}
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	var conn = db.Conn
	conn.Exec("create table menu(id integer primary key, title text)")
}

func cleanup() {
	if err := db.CloseDB(); err != nil {
		log.Fatal(err)
	}
}

func initApi() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().Format("2006/01/02 15:04:05")
		w.Write([]byte(now))
	})

	http.HandleFunc("/menu/list/:page/:limit", api.GetMenuList)
	http.HandleFunc("/menu/:id", api.GetMenuInfo)
}

func main() {
	initApi()

	srv := &http.Server{
		Addr:         config.C.ApiPort,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
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
	log.Printf("http server started on \033[32m%s\033[0m\n", config.C.ApiPort)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-quit
	println("\nserver closed")
}
