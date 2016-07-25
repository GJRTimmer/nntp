# NNTP #
Golang NNTP Library with multi connection support

### Usage Multi Connection ###
The multiple connection pool handles all the connections and the devides the work along all the available connections.
For this the connection pool requires a channel to send the work to, the 'Job Queue'.

On this channel work for a NNTP connection can be send.
The work is of type ```Request```, and the connection pool provides a return channel where is will deliver all the completed work. The completed work is of type ```Response```.

```go
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
```

#### Create Connection Pool ####
1. [Create *ServerInfo details](README.md#serverinfo-structure)
2. Create a JobQueue Channel
3. Provide max number of connection
4. Creat the connection pool

```go

import "github.com/GJRTimmer/nntp"

// Step 1
s := &ServerInfo {
  ...
}

// Step 2
jobQueue := make(chan *nntp.Request, maxConnections + 10) // additional space on queue

// Step 3
nrConn := 30

// Step 4
pool := nntp.NewConnectionPool(s, jobQueue, nrConn)
```

#### Starting the Connection Pool ###
After the connection pool has been created, it must be started.
Starting the pool will automatically create the number of connections provided and will start connecting.
```go
pool.Start()
```

#### Listening for Responses ####
After the pool has been started, you must start listening for responses of work which has been completed.
This can be achived by fetching the response channel on which all the work from all the connections will be deliverd. This channel can be fetched by calling the ```Collect()``` function on the connections pool.
```go
responses := pool.Collect()

// This example does not provide a full implementation which uses a sync group.

// Listen for reponses.
go func() {
    for {
        select {
            case r := <- responses:
                // Do something with r
        }
    }
}()
```

#### Start Initial Work ###
When the reponse listener is up, work can be send to the job queue to be processed.

NOTE: If you have a lot of work, I suggest the following approach: Create initial work on the queue, enough work for all the connections, example 30 connections, job queue of 40 or 50, and for the inital work, just fill up the entire jobqueue. Then every time you receive a response, you first create a new work item of all the work which has to be done and place it on the job queue. This way the jobqueue will stay filled up and all the connections will have enough to do. (TODO: Link to implementation which uses this)

```go
// From github.com/GJRTimmer/nzb package
n, _ nzb.Parse(fh)
chunks := n.GenerateChunkList()
initWork := chunks(nrConn + 10)
for _, w := range initWork {
    // Create `CHECK` Request
    req := nntp.NewRequest(w.Segment.ID, w.Groups, nntp.CHECK)
    jobQueue <- req
}
```
