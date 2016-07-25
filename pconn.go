package nntp

import (
    "fmt"
)

// Implementation for a pooled connection
var _ PoolConn = (*poolConn)(nil)

// PoolConn represents a connection to a NNTP server
// within a connection pool
type PoolConn interface {
    Conn
    Start()
    Stop()
    ResponseChannel() chan *Response
}

type poolConn struct {
    *conn

    reqQueue    chan *Request // Queue to get Request(s) from
    pool        chan chan *Request // Connection Pool
    quit        chan bool //quit channel to close the connection

    respChan    chan *Response // Response Channel
}

// NewPoolConn create a new NNTP connection for use in a connection pool
func NewPoolConn(id int, i *ServerInfo, pool chan chan *Request) PoolConn {
    return &poolConn {
        conn: &conn {
            id: id,
            Info: i,
        },

        reqQueue: make(chan *Request), // Make new request channel
        pool: pool, // Assign provided pool as workerpool
        quit: make(chan bool), // create new quit channel
        respChan: make(chan *Response),
    }
}

// Process work
// Exported so this can be used at users digression
func (p *poolConn) process(req *Request) *Response {
    resp := req.generateNewResponse(p)

    if p.conn == nil {
        resp.Error = fmt.Errorf("no connection to server")
        return resp
    }

    switch req.Command {
    case CHECK:
        return p.checkArticle(req, resp)
    case FETCH:
        return p.fetchArticle(req, resp)
    default:
        resp.Error = fmt.Errorf("operation: '%s' not implemented", req.Command)
    }

    return resp
}

func (p *poolConn) Start() {
    go func() {

        p.Connect()

        for {
            // Add connection request queue (reqQueue)
            // to the worker pool
            p.pool <- p.reqQueue

            select {
            case req := <- p.reqQueue:
                // Incoming request to process'
                r := p.process(req)
                // Send response through response channel
                p.respChan <- r

            case <- p.quit:
                // Requested to stop processing
                // Close connection
                p.Close()
                return
            }

        }
    }()
}

func (p *poolConn) Stop() {
    go func() {
        // Stop Main Connection GoRoutine
        p.quit <- true

        // Stop Response Channel Merge (FanIn) GoRoutine
        close(p.respChan)
    }()
}

func (p *poolConn) ResponseChannel() chan *Response {
    return p.respChan
}

// EOF
