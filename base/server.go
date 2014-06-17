package base

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type terminateSignal struct{}

func (t terminateSignal) String() string {
	return "Normal program termination."
}

func (t terminateSignal) Signal() {}

func ServeMain(ltype, laddr string, server func(net.Listener)) (err error) {
	// Create the socket to listen on:
	var l net.Listener
	l, err = net.Listen(ltype, laddr)
	if err != nil {
		return
	}

	// NOTE(jsd): Unix sockets must be removed before being reused.

	// Handle common process-killing signals so we can gracefully shut down:
	// TODO(jsd): Go does not catch Windows' process kill signals (yet?)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// Start a server:
		server(l)

		// Signal completion:
		sigc <- terminateSignal{}
		signal.Stop(sigc)
	}()

	// Wait for a termination signal (normal or otherwise):
	sig := <-sigc

	// Stop listening:
	l.Close()

	// Delete the unix socket, if applicable:
	if ltype == "unix" {
		os.Remove(laddr)
	}

	// Our own terminateSignal is not an error condition:
	if _, ok := sig.(*terminateSignal); ok {
		return nil
	}

	return errors.New(sig.String())
}
