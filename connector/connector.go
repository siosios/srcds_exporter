package connector

import (
	"fmt"
	"time"

	steam "github.com/galexrt/go-steam"
	cache "github.com/patrickmn/go-cache"
)

// Connector struct contains the connections
type Connector struct {
	connections map[string]*Connection
}

// NewConnector creates a new Connector object
func NewConnector() *Connector {
	return &Connector{
		connections: make(map[string]*Connection),
	}
}

// GetConnections holds all connections and reconnects/reopens them if necessary
func (cn *Connector) GetConnections() (map[string]*Connection, error) {
	return cn.connections, nil
}

// NewConnection Add a new connection and initiates first contact connection
func (cn *Connector) NewConnection(name string, opts *ConnectionOptions) error {
	if _, ok := cn.connections[opts.Addr]; ok {
		return nil
	}
	con, err := steam.Connect(opts.Addr,
		&steam.ConnectOptions{
			RCONPassword: opts.RconPassword,
			Timeout:      opts.ConnectTimeout,
		})
	if err != nil {
		return err
	}
	var (
		conTimeoutParsed   time.Duration
		cacheTimeoutParsed time.Duration
	)
	conTimeoutParsed, err = time.ParseDuration(opts.ConnectTimeout)
	if err != nil {
		return err
	}
	cacheTimeoutParsed, err = time.ParseDuration(opts.CacheTimeout)
	if err != nil {
		return err
	}
	fmt.Print(cacheTimeoutParsed)
	// TODO make cache time configurable?
	cn.connections[opts.Addr] = &Connection{
		Name:  name,
		con:   con,
		cache: *cache.New(cacheTimeoutParsed, 11*time.Second),
		opts: map[string]string{
			"Address":      opts.Addr,
			"RCONPassword": opts.RconPassword,
			"Timeout":      opts.ConnectTimeout,
		},
		created: time.Now().Add(conTimeoutParsed),
	}
	return nil
}

// CloseAll closes all open connections
func (cn *Connector) CloseAll() {
	for _, con := range cn.connections {
		con.Close()
	}
}
