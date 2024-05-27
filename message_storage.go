package main

import (
	"sync"
)

// Message represents a WhatsApp message
type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

// MessageStorage stores messages in memory
type MessageStorage struct {
	mu       sync.Mutex
	messages []Message
}

// AddMessage adds a new message to the storage
func (s *MessageStorage) AddMessage(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, msg)
}

// GetMessages returns all stored messages
func (s *MessageStorage) GetMessages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.messages
}

// messageStorage is an instance of MessageStorage to store messages in memory
var messageStorage = &MessageStorage{messages: make([]Message, 0)}
