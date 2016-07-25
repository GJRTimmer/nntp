package nntp

import (
    "sync"
)

// Implementation for connection pool of 'poolConn(s)'

// ConnectionPool to NNTP server
type ConnectionPool struct {
    Info            *ServerInfo

    poolChannel     chan chan *Request
    maxConnections  int
    reqQueue        chan *Request // Request queue
    pool            map[int]PoolConn
    quit            chan bool
}

// NewConnectionPool create new connection pool to NNTP Server
func NewConnectionPool(i *ServerInfo, reqQueue chan *Request, maxConnections int) *ConnectionPool {
    poolChan := make(chan chan *Request, maxConnections)

    return &ConnectionPool {
        Info: i,
        poolChannel: poolChan,
        maxConnections: maxConnections,
        reqQueue: reqQueue,
        pool: make(map[int]PoolConn),
        quit: make(chan bool),
    }
}

// Start connection pool
func (cp *ConnectionPool) Start() {

    // Start Connections
    for i:= 0; i < cp.maxConnections; i++ {
        cp.pool[i] = cp.newPoolConn(i, cp.Info, cp.poolChannel)
        cp.pool[i].Start()
    }

    go cp.dispatch()
}

// Stop connection pool
func (cp *ConnectionPool) Stop() {

    // Stop Connections
    for _, conn := range cp.pool {
        conn.Stop()
    }

    // Stop Connection pool dispatcher
    cp.quit <- true

    // Clean pool
    cp.pool = make(map[int]PoolConn)
}

func (cp *ConnectionPool) dispatch() {
    for {
        select {
        case req := <- cp.reqQueue:
            // Dispatch request to available connection
            go func() {
                // Fetch connection reqQueue from poolChannel
                connReqQueue := <- cp.poolChannel
                // Send request to connection
                connReqQueue <- req
            }()
        case <- cp.quit:
            // Received termination request
            return
        }
    }
}

// Collect *Response(s) from ConnectionPool
func (cp *ConnectionPool) Collect() chan *Response {

    wg := new(sync.WaitGroup)
    out := make(chan *Response)

    wg.Add(cp.maxConnections)
    for _, conn := range cp.pool {
        go collect(conn.ResponseChannel(), out, wg)
    }

    // Start goroutine to close out channel
    // when all collectors have been stopped
    // when alll connections are closed
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}

// Function used to run within a go routine
func collect(c <-chan *Response, out chan<- *Response, wg *sync.WaitGroup) {

    // Collect Response(s)
    for r := range c {
        out <- r
    }

    // Channel has been closed
    // Signal Done
    wg.Done()
}

// EOF
