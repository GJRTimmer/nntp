package nntp

import (
    "fmt"
)

// NNTP Operations

// Close the connection
func (c *conn) Close() error {
    return c.conn.Close()
}

//SwitchGroup will change the current group, using the supplied
//group name. It returns an error, if any
// Exported so this can be used at users digression
func (c *conn) SwitchGroup(group string) error {

    if group == c.group {
		return nil
	}

    if c.conn == nil {
        return fmt.Errorf("[%s] no connection to server", c)
    }

	cmd, err := c.conn.Cmd("GROUP %s", group)
	if err != nil {
		return err
	}

	c.conn.StartResponse(cmd)
    defer c.conn.EndResponse(cmd)

    // Read response code
	_, _, err = c.conn.ReadCodeLine(211)
    if err != nil {
        return err
    }

	// No error; set group
	c.group = group

	return err
}

// ArticleExists check if article with provided id exists
// on the NNTP server
// Exported so this can be used at users digression
func (c *conn) ArticleExists(id string) (bool, error) {

    if c.conn == nil {
        return false, fmt.Errorf("[%s] no connection to server", c)
    }

    cmd, err := c.conn.Cmd("STAT <%s>", id)
    if err != nil {
        return false, err
    }

    c.conn.StartResponse(cmd)
    defer c.conn.EndResponse(cmd)

    // Read response code
    code, _, err := c.conn.ReadCodeLine(223)
    if err != nil {
        return false, err
    }

    switch code {
    case 223:
        // Article exists
        return true, nil
    case 412:
        // Nnnewsgroup selected
        return false, fmt.Errorf("[%s] no newsgroup selected", c)
    case 420:
        // Current article number is invalid
        return false, fmt.Errorf("[%s] current article id <%s> is invalid", c, id)
    case 430:
        // Article NOT FOUND
        return false, nil
    default:
        // Unexpected code ARTICLE NOT FOUND
        return false, fmt.Errorf("[%s] unexcepted return code: %d", c, code)
    }
}

// FetchArticle fetch article with id from NNTP server
func (c *conn) FetchArticle(id string) ([]byte, error) {

    if c.conn == nil {
        return nil, fmt.Errorf("[%s] no connection to server", c)
    }

    cmd, err := c.conn.Cmd("BODY <%s>", id)
    if err != nil {
        return nil, err
    }

    c.conn.StartResponse(cmd)
    defer c.conn.EndResponse(cmd)

    code, _, err := c.conn.ReadCodeLine(222)
    if err != nil {
        return nil, err
    }

    switch code {
    case 222:
        // Body follows
        return c.conn.ReadDotBytes()
    case 412:
        // No newsgroup selected
        return nil, fmt.Errorf("[%s] no newsgroup selected", c)
    case 420:
        // Current article number is invalid
        return nil, fmt.Errorf("[%s] current article id <%s> is invalid", c, id)
    case 423:
        // no article with that message-id
        return nil, fmt.Errorf("[%s] no article with id <%s>", c, id)
    default:
        return nil, fmt.Errorf("[%s] unexpected return code: %d", c, code)
    }

}

// +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
/// Private

func (c *conn) checkArticle(req *Request, resp *Response) *Response {

    for _, g := range req.Article.Groups {
        c.SwitchGroup(g)
        b, err := c.ArticleExists(req.Article.ID)
        if err != nil {
            resp.Error = err
        } else {
            resp.Article.Exists = b
            break
        }
    }

    // Clear request memory
    req = nil

    return resp
}

func (c *conn) fetchArticle(req *Request, resp *Response) *Response {

    for _, g := range req.Article.Groups {
        c.SwitchGroup(g)
        b, err := c.FetchArticle(req.Article.ID)
        if err != nil {
            resp.Error = err
        } else {
            resp.Article.Content = b
            resp.Article.Exists = true
            break
        }
    }

    // Clear request memory
    req = nil

    return resp
}

// EOF
