package main

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/acoshift/configfile"
	_ "github.com/lib/pq"

	"github.com/acoshift/postlist/pkg/app"
)

var config = configfile.NewReader("config")

func main() {
	db, err := sql.Open("postgres", config.String("sqlurl"))
	if err != nil {
		log.Fatal(err)
	}

	{
		// init table
		table, err := ioutil.ReadFile("table.sql")
		if err != nil {
			log.Fatal("can not load table.sql;", err)
		}
		_, err = db.Exec(string(table))
		if err != nil {
			log.Println("exec table.sql error;", err)
		}
	}

	h := app.MakeHandler(db)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown error;", err)
		return
	}
	log.Println("server shutdown")
}
