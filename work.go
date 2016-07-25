package nntp

import (
    "fmt"
)

// Operation defines what kind
// of operation to perform for provided article
// on the connection
type Operation uint8

const (
    // CHECK Operation
    // Check if article exists
    CHECK Operation = iota

    // FETCH Operation
    // Perform Fetch Operation for article
    FETCH
)

// Request on connection
type Request struct {
    Article     *Article        // Article Information
    Command     Operation       // Command Operation to perform
}

// Response of work on connection
type Response struct {
    Article     *Article
    Commmand    Operation       // Operation which was executed
    Error       error           // Error of operation if any
    Source      string          // source of response HOST:PORT of connection
}

// Article describes the article information
type Article struct {
    ID          string          // Article ID
    Groups      []string        // Groups which holds the article

    // Result
    Exists      bool            // Flag which defines if article exists on NNTP server
    Content     []byte          // Content of fetched article, filled on operation fetch, otherwise NIL
}

func (a *Article) String() string {
    return fmt.Sprintf("<%s>", a.ID)
}

// NewRequest create new request for connection
func NewRequest(id string, groups []string, oper Operation) *Request {
    return &Request {
        Article: &Article {
            ID: id,
            Groups: groups,
        },
        Command: oper,
    }
}

func (req *Request) generateNewResponse(p *poolConn) *Response {
    return &Response {
        Article: req.Article,
        Commmand: req.Command,
        Source: p.String(),
    }
}

// EOF
