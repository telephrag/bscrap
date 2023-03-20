package main

import (
	"bscrap/bscrap_srv"
	"bscrap/config"
	"bscrap/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func getLogFile(path string) *os.File {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		_, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			log.Fatal(err)
		}
		return f
	}
	return f
}

func main() {

	logFile := getLogFile("/data/log.log")
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)

	mi, err := db.ConnectMongo(config.DBUri)
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println("application startup")

	bScrapEnv := &bscrap_srv.Env{Mi: mi}
	go func() {
		err := http.ListenAndServe(
			":8080",
			bscrap_srv.Run(bScrapEnv),
		)
		if err != nil {
			log.Panic(err)
		}
	}()

	interupt := make(chan os.Signal, 1)
	signal.Notify(interupt, syscall.SIGTERM, syscall.SIGINT)
	<-interupt

	log.Print("application shutdown\n\n")
}
