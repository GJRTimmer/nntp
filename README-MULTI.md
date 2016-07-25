# NNTP #
Golang NNTP Library with multi connection support

### Usage Multi Connection ###
The multiple connection pool handles all the connections and the devides the work along all the available connections.
For this the connection pool requires a channel to send the work to, the 'Job Queue'.

On this channel work for a NNTP connection can be send.
The work is of type ```Request```, and the connection pool provides a return channel where is will deliver all the completed work. The completed work is of type ```Response```.


#### Create Connection Pool ####
1. [Create *ServerInfo details](README.md#serverinfo-structure)
2. Create a JobQueue Channel
3. Provide max number of connection
