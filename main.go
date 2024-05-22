package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var client *whatsmeow.Client

func ScanQrCode() {
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err := client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
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

func ReceiveAllMessages(v *events.Message) {
	if !v.Info.IsFromMe {
		sender := v.Info.Sender.String()
		switch {
		case v.Message.GetConversation() != "":
			fmt.Printf("Text message from %s: %s\n", sender, v.Message.GetConversation())
		case v.Message.GetImageMessage() != nil:
			fmt.Printf("Image message from %s: %s\n", sender, v.Message.GetImageMessage().GetCaption())
		case v.Message.GetVideoMessage() != nil:
			fmt.Printf("Video message from %s: %s\n", sender, v.Message.GetVideoMessage().GetCaption())
		case v.Message.GetAudioMessage() != nil:
			fmt.Printf("Audio message from %s\n", sender)
		case v.Message.GetDocumentMessage() != nil:
			fmt.Printf("Document message from %s: %s\n", sender, v.Message.GetDocumentMessage().GetTitle())
		default:
			fmt.Printf("Other type of message from %s\n", sender)
		}
	}
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		ReceiveAllMessages(v)
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:wasopingi.db?_foreign_keys=on", dbLog)
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

	ScanQrCode()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
