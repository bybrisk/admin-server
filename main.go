package main

import (
	"log"
	"os"
	"context"
	"time"
	"net/http"
	"os/signal"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"github.com/go-openapi/runtime/middleware"
	deliveryHl "github.com/bybrisk/delivery-api/handlers"
	accountHl "github.com/bybrisk/businessAccount-api/handlers"
	agentsHl "github.com/bybrisk/agents-api/handlers"
	//"github.com/bybrisk/swagger-automation-docs/docs"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":80","Bind address for the server")

func main() {
	env.Parse()

	l := log.New(os.Stdout,"map-api",log.LstdFlags)
	
	//registering all handlers
	deliveryHandler := deliveryHl.NewDelivery(l)
	accountHandler := accountHl.NewAccount(l)
	agentsHandler := agentsHl.NewAgents(l)

	serveMux := mux.NewRouter()

	//Get Router
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/account/{id}",accountHandler.GetAccountDetail)
	getRouter.HandleFunc("/agents/all/{id}",agentsHandler.GetAllAgents)
	getRouter.HandleFunc("/agents/delete/{id}",agentsHandler.DeleteAgents)
	getRouter.HandleFunc("/agents/one/{id}",agentsHandler.GetOneAgent)

	//Post Router
	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/addDelivery", deliveryHandler.AddDelivery)

	postRouter.HandleFunc("/account", accountHandler.AddNewAccount)
	postRouter.HandleFunc("/accountUpdate", accountHandler.UpdateAccount)
	postRouter.HandleFunc("/accountUpdate/password", accountHandler.UpdateAccountPassword)

	postRouter.HandleFunc("/agents/create", agentsHandler.AddNewAgent)
	postRouter.HandleFunc("/agents/update", agentsHandler.UpdateAgent)
	//Put Router
	//putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	//putRouter.HandleFunc("/delivery/{id}", deliveryHandler.Updatedelivery)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	//Create a new server
	s := http.Server{
		Addr: *bindAddress,			//configure the bind address
		Handler: serveMux,	//set the default handler
		ErrorLog: l,	//set the logger for the server
		ReadTimeout: 5 * time.Second,	//max time for read request for the client
		WriteTimeout: 10 * time.Second,	//Max time for write response for client
		IdleTimeout: 120 * time.Second, //Max timme for connection using TCP Keep alive 
	}

	//start the server
	go func () {
		l.Println("Starting bybrisk maps-api server")
		err := s.ListenAndServe()	
		if err!=nil {
			l.Printf("Error starting bybrisk maps-api server: %s\n",err)
			os.Exit(1)
		}
	}()

	//trap sigterm or intrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Println("Recieved terminate, graceful shutdown bybrisk maps-api server", sig)

	tc,_ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}