[![GoDoc](https://godoc.org/github.com/voloshink/dggchat?status.svg)](https://godoc.org/github.com/voloshink/dggchat) [![Go Report Card](https://goreportcard.com/badge/github.com/voloshink/dggchat)](https://goreportcard.com/report/github.com/voloshink/dggchat)

# dggchat
[Destinygg](https://www.destiny.gg) chat go bindings. You can acquire a login key [here](https://www.destiny.gg/profile/authentication).

# Simple example

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/voloshink/dggchat"
)

func main() {
	// Create a new client
	dgg, err := dggchat.New("loginkey")
	if err != nil {
		log.Fatalln(err)
	}

	// Open a connection
	err = dgg.Open()
	if err != nil {
		log.Fatalln(err)
	}

	// Cleanly close the connection
	defer dgg.Close()

	dgg.AddMessageHandler(onMessage)
	dgg.AddErrorHandler(onError)

	// Wait for ctr-C to shut down
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT)
	<-sc
}

func onMessage(m dggchat.Message, s *dggchat.Session) {
	log.Printf("New message from %s: %s\n", m.Sender.Nick, m.Message)

	if m.Message == "!test" {
		s.SendPrivateMessage(m.Sender.Nick, "testing")
	}
}

func onError(e string, s *dggchat.Session) {
	log.Printf("error %s\n", e)
}
```
