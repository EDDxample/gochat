# Go Chat

Go implementation of my old TCP socket chat

Also me trying to learn this go project layout: https://github.com/golang-standards/project-layout

### TODO:
- Client implementation
- Handle Custom Usernames
- Merge apps into one binary with client/server modes
- Actual login system (?) sqlite (??)

### How to use

```sh
# server
go run .\cmd\server -h <host> -p <port>

# client
go run .\cmd\client -h <host> -p <port>
``` 

### Resources
Previous version: https://github.com/EDDxample/Socket_Client-Server_Chat
