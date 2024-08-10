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
	"sy_backend/middleware"
	"syscall"
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
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("resource/menu"))
	//mux.Handle("/resource/menu/", http.StripPrefix("/resource/menu", fs))
	mux.HandleFunc(
		"/resource/menu/",
		middleware.Cors(func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/resource/menu", fs).ServeHTTP(w, r)
		}),
	)

	// 获取菜列表
	mux.HandleFunc("/menu/list", middleware.Cors(api.GetMenuList))

	mux.HandleFunc("/menu", api.PostMenu)        // 增菜
	mux.HandleFunc("/menu/{id}", api.DeleteMenu) // 删菜

	// mux.HandleFunc("/menu/image/{id}", api.MenuImage)   // 示例图片的增删
	mux.HandleFunc("/menu/{id}/{field}", api.PatchMenu) // 更新菜字段

	srv := &http.Server{
		Addr:    config.Conf.ApiPort,
		Handler: mux,
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
