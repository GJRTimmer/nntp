# NNTP #
Golang NNTP Library with multi connection support

## Install ##
```bash
go get github.com/GJRTimmer/nntp
```

## Package Info ##
Conn Interface is primarly used for a single connection to a NNTP server.
```go
type Conn interface {
    Connect() bool
    Close() error

    SwitchGroup(group string) error
    ArticleExists(id string) (bool, error)
    FetchArticle(id string) ([]byte, error)
}
```

PoolConn is a connection used within a connection pool for multi connection to a NNTP server.
PoolConn inherit Conn.
```go
type PoolConn interface {
    Conn
    Start()
    Stop()
    ResponseChannel() chan *Response
}
```

## Usage ##
For both single and multi user connection, details of the NNTP server must be provided.

### ServerInfo Structure ###
- Host ```string```
- Port ```uint16``` Port number 0-65536
- TLS ```bool``` Boolean to use TLS for connection
- Auth ```ServerAuth``` If Auth == nil no authentication is preformed by the library
  - Username ```string```
  - Password ```string```

```go
// ServerInfo for nntp connection
type ServerInfo struct {
    Host        string
    Port        uint16
    TLS         bool
    Auth        *ServerAuth
}

// ServerAuth authentication info for ServerInfo
type ServerAuth struct {
    Username    string
    Password    string
}
```

### Example Server Info ###
```go
    sInfo := &nntp.ServerInfo {
        Host: "<HOST>",
        Port: <PORT>, // if port variable != uint16, don't forget to cast uint16(port)
        TLS: true,
        Auth: &nntp.ServerAuth {
            Username: "<USERNAME>",
            Password: "<PASSWORD>",
        },
    }
```
  

### [Usage Single Connection](README-SINGLE.md) ###

### [Usage Multi Connection](README-MULTI.md) ###

