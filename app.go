package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	Router      *mux.Router
	Port        int
	WaitTimeout time.Duration
}

func NewApp(port int, waitTimeout time.Duration) *App {
	return &App{
		Port:        port,
		WaitTimeout: waitTimeout,
	}
}

func (a *App) HandleMain(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling connection: %+v", r)
	w.Write([]byte("Hello World!"))
}

func (a *App) HandleSendWOL(w http.ResponseWriter, r *http.Request) {
	log.Debugf("Handling connection: %+v", r)
	params := mux.Vars(r)
	hwAddr, err := net.ParseMAC(params["mac"])
	if err != nil {
		log.Errorf("Invalid request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%v\"}\n", err)))
		return
	}

	macAddress := hwAddr.String()
	log.Debugf("Sending WakeOnLAN with MAC: %v", macAddress)
	if err := SendWakeOnLAN(macAddress); err != nil {
		log.Errorf("Error sending WakeOnLAN packet: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("{\"message\": \"%v\"}\n", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"OK\"}\n"))
}

func (a *App) InitializeRouter() {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/", a.HandleMain).Methods("GET")
	a.Router.HandleFunc("/wakeonlan/{mac}", a.HandleSendWOL).Methods("GET")
}

func (a *App) Initialize() {
	log.Debugf("Initializing application")
}

func (a *App) Run() error {
	// Gorilla mux graceful shutdown
	// https://github.com/gorilla/mux#graceful-shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", a.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("%v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(ch, os.Interrupt)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), a.WaitTimeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Warningf("Shutting down")
	return nil
}
