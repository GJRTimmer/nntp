package nntp

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

// EOF
