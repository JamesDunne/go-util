package base

import (
	"errors"
	"net"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

type terminateSignal struct{}

func (t terminateSignal) String() string {
	return "Normal program termination."
}

func (t terminateSignal) Signal() {}

type Dialable struct {
	Network, Address string
}

func ParseDialable(s string) (d *Dialable, err error) {
	var u *url.URL
	u, err = url.Parse(s)
	if err != nil {
		return nil, err
	}

	var ltype, laddr string
	ltype = u.Scheme
	if ltype == "unix" {
		if u.Host != "" {
			return nil, errors.New("Dialable unix URI must have blank host, e.g. unix:///path/to/socket")
		}
		laddr = u.Path
	} else {
		laddr = u.Host
	}

	return &Dialable{Network: ltype, Address: laddr}, nil
}

type Listenable struct {
	Network, Address string
}

func ParseListenable(s string) (l *Listenable, err error) {
	var u *url.URL
	u, err = url.Parse(s)
	if err != nil {
		return nil, err
	}

	var ltype, laddr string
	ltype = u.Scheme
	if ltype == "unix" {
		if u.Host != "" {
			return nil, errors.New("Listenable unix URI must have blank host, e.g. unix:///path/to/socket")
		}
		laddr = u.Path
	} else {
		laddr = u.Host
	}

	return &Listenable{Network: ltype, Address: laddr}, nil
}

// Main method to start up a server.
func ServeMain(la *Listenable, server func(net.Listener) error) (sig os.Signal, err error) {
	// Create the folder for any unix sockets to live in:
	if la.Network == "unix" {
		// TODO(jsd): 0770 permissions on the folder?
		// TODO(jsd): Hide mkdir error?
		os.MkdirAll(filepath.Dir(la.Address), os.FileMode(0770)|os.ModeDir)
	}

	// Create the socket to listen on:
	var l net.Listener
	l, err = net.Listen(la.Network, la.Address)
	if err != nil {
		return
	}

	sig, err = Daemonize(func() error { return server(l) })

	// Stop listening:
	l.Close()

	// Delete the unix socket, if applicable:
	if la.Network == "unix" {
		os.Remove(la.Address)
	}

	return
}

type Daemon func() error

func Daemonize(start Daemon) (sig os.Signal, err error) {
	// Handle common process-killing signals so we can gracefully shut down:
	// TODO(jsd): Go does not catch Windows' process kill signals (yet?)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// Start a server; `err` will be returned to the caller:
		err = start()

		// Signal completion:
		sigc <- terminateSignal{}
		signal.Stop(sigc)
	}()

	// Wait for a termination signal (normal or otherwise):
	sig = <-sigc

	return
}
