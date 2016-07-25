# NNTP #
Golang NNTP Library with multi connection support

### Usage Single Connection ###
To create a new single connection the function ```NewConn(id int, i *ServerInfo)``` can be used. The id is just to identify the connection number.

#### Create ####
```go
// Info is pointer to ServerInfo
conn := nntp.NewConn(0, info)
```

#### Available Operations on Single Connection ####
```go
// Connect to NNTP Server
// Returns true is connection succeeded, false otherwise 
Connect() (bool, error)

// Close connection
Close() error 

// Switch to newsgroup
SwitchGroup(group string) error

// Check if article with ID exists on NNTP server
// Uses the NNTP STAT command
ArticleExists(id string) (bool, error) 

// Fetch the article with id from NNTP server
FetchArticle(id string) ([]byte, error) 
```

#### Example Single Connection ####
```go
import (
    "fmt"
    "os"
    
    "github.com/GJRTimmer/nntp"
)

func main() {
    
    // Setup serverInfo
    sInfo := &nntp.ServerInfo {
        Host: "<HOST>",
        Port: <PORT>, // if port comes from a flag, don't forget to cast uint16(port)
        TLS: true,
        Auth: &nntp.ServerAuth {
            Username: "<USERNAME>",
            Password: "<PASSWORD>",
        },
    }
    
    conn := nntp.NewConn(0, sInfo)
    connected, err := conn.Connect()
    if !connected {
        fmt.Printf("Failed to connect: %s\n", err)
    }
    defer conn.Close()
    
    // Switch to newsgroup
    conn.SwitchGroup("alt.binaries.boneless")
    
    // Check if article exists
    // Article ID is fake
    ok, err := conn.ArticleExists("part1of181.***********@powerpost2000AA.local")
    if err != nil {
        fmt.Println(err)
    }
    
    if ok {
        // Article is present
        // Article ID is fake
        // Fetch the article
        content, err := conn.FetchArticle("part1of181.***********@@powerpost2000AA.local")
        if err != nil {
            // User TODO: Do something with content
            // Release memory
            content = nil
        }
    }
    
    os.Exit(0)
}

```
