package server

import (
	"bufio"
	"context"
	"io"
	"time"
)

// Handler is a callback function to process messages from io.Reader.
type Handler func(Message) error

// Message is a message struct to be received/published.
type Message struct {
	Data []byte
}

// StreamReader is the server listener that reads messages from the
// byte stream provided by io.Reader.
type StreamReader struct {
	reader  io.ReadCloser
	handler Handler
}

// NewStreamReader returns a new instance of StreamReader.
func NewStreamReader(reader io.ReadCloser, handler Handler) *StreamReader {
	return &StreamReader{
		reader:  reader,
		handler: handler,
	}
}

// ListenAndServe continuously reads bytes from io.Reader
// and invokes handlers once the complete message is received.
func (l *StreamReader) ListenAndServe() error {
	for {
		scanner := bufio.NewScanner(l.reader)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			if l.handler == nil {
				continue
			}

			if err := l.handler(Message{Data: scanner.Bytes()}); err != nil {
				return err
			}
		}

		time.Sleep(time.Second)
	}
}

// Shutdown closes the io.Reader.
func (l *StreamReader) Shutdown(context.Context) error {
	return l.reader.Close()
}
