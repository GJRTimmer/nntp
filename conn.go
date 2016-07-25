package nntp

import (
    "crypto/tls"
    "fmt"
    "net/textproto"
)

// Interface implementation check
var _ Conn = (*conn)(nil)

// Conn represent(s) a single connection to a NNTP Server
type Conn interface {
    Connect() (bool, error)
    Close() error

    SwitchGroup(group string) error
    ArticleExists(id string) (bool, error)
    FetchArticle(id string) ([]byte, error)
}

// Conn represents a NNTP Connection
type conn struct {
    Info        *ServerInfo

    id          int
    conn        *textproto.Conn
    group       string
}

// NewConn create a new NNTP connection based on provided ServerInfo
func NewConn(id int, i *ServerInfo) Conn {
    return &conn {
        id: id,
        Info: i,
    }
}

func (c *conn) String() string {
    return fmt.Sprintf("%s:%d", c.Info.Host, c.id)
}

// Connect to NNTP server
func (c *conn) Connect() (bool, error) {
    // (Re)connect
    // Try to connect to newsgroup server
    // if unable to connect start auto-reconnect timer
    if c.Info.TLS {
        // Secured connection is requested
        err := c.dialTLS()
        if err != nil {
            return false, err
        }

        return true, nil
    }

    err := c.dial()
    if err != nil {
        return false, err
    }

    return true, nil
}

func (c *conn) dial() error {

	var err error
    c.conn, err = textproto.Dial("tcp", fmt.Sprintf("%s:%d", c.Info.Host, c.Info.Port))
	if err != nil {
		return err
	}

    _, _, err = c.conn.ReadCodeLine(20)
	if err != nil {
		c.conn.Close()
		return err
	}

    if c.Info.Auth != nil {
        err = c.authenticate()
    	if err != nil {
    		return err
    	}
    }

	return nil
}

func (c *conn) dialTLS() error {

    tlsConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", c.Info.Host, c.Info.Port), nil)
	if err != nil {
		return err
	}

    c.conn = textproto.NewConn(tlsConn)
    _, _, err = c.conn.ReadCodeLine(20)
	if err != nil {
		c.conn.Close()
		return err
	}

    if c.Info.Auth != nil {
        err = c.authenticate()
    	if err != nil {
    		return err
    	}
    }

	return nil
}

//Authenticate will authenticate with the NNTP server, using the supplied
//username and password. It returns an error, if any
func (c *conn) authenticate() error {

    u, err := c.conn.Cmd("AUTHINFO USER %s", c.Info.Auth.Username)
	if err != nil {
		return err
	}

    c.conn.StartResponse(u)
	code, _, err := c.conn.ReadCodeLine(381)
    c.conn.EndResponse(u)

    switch code {
	case 481, 482, 502:
		//failed, out of sequence or command not available
		return err
	case 281:
		//accepted without password
		return nil
	case 381:
		//need password
		break
	default:
		return err
	}

    p, err := c.conn.Cmd("AUTHINFO PASS %s", c.Info.Auth.Password)
	if err != nil {
		return err
	}

	c.conn.StartResponse(p)
	code, _, err = c.conn.ReadCodeLine(281)
    c.conn.EndResponse(p)

	return err
}

// EOF
