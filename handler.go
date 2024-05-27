package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// ScanQrCode scans and prints the QR code for login.
func ScanQrCode(client *whatsmeow.Client) {
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		fmt.Println("Already logged in.")
		err := client.Connect()
		if err != nil {
			panic(err)
		}
	}
}

// scanQrCodeHandler handles the QR code scanning process
func scanQrCodeHandler(w http.ResponseWriter, r *http.Request) {
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			http.Error(w, "Failed to connect client", http.StatusInternalServerError)
			return
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				qrResp := qrResponse{Code: evt.Code}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(qrResp)
				return
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		fmt.Println("Already logged in.")
		err := client.Connect()
		if err != nil {
			http.Error(w, "Failed to connect client", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Already logged in"))
	}
}

type qrResponse struct {
	Code string `json:"code"`
}

// getMessagesHandler returns all stored messages
func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages := messageStorage.GetMessages()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// ReceiveAllMessages processes incoming messages and stores them
func ReceiveAllMessages(v *events.Message) {
	if !v.Info.IsFromMe {
		msg := Message{
			Sender:  v.Info.Sender.String(),
			Content: v.Message.GetConversation(),
		}
		messageStorage.AddMessage(msg)
	}
}
