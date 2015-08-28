package main

import "time"

// messae represents a single message
type message struct {
	Name    string
	Message string
	When    time.Time
}
