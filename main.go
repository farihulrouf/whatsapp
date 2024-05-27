package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	client    *whatsmeow.Client
	container *sqlstore.Container
)

func init() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	var err error
	container, err = sqlstore.New("sqlite3", "file:whatsapp.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}
	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/scan", scanQrCodeHandler).Methods("GET")
	r.HandleFunc("/messages", getMessagesHandler).Methods("GET")

	go func() {
		fmt.Println("Starting server on port 8080")
		if err := http.ListenAndServe(":8080", r); err != nil {
			fmt.Println("Failed to start server:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	container.Close()
	fmt.Println("Server stopped")
}

func scanQrCodeHandler(w http.ResponseWriter, r *http.Request) {
	qrChan, _ := client.GetQRChannel(context.Background())
	err := client.Connect()
	if err != nil {
		http.Error(w, "Failed to connect to WhatsApp: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, w)
			return
		}
	}

	http.Error(w, "Failed to retrieve QR code", http.StatusInternalServerError)
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Fetch and return all messages
	// You would implement this function based on your requirements
	// For demonstration purposes, let's just return a dummy response
	dummyMessages := []string{"Message 1", "Message 2", "Message 3"}
	json.NewEncoder(w).Encode(dummyMessages)
}

func eventHandler(evt interface{}) {
	// Handle incoming WhatsApp events here
	// You would implement this function based on your requirements
}
