# Snowcast

This project emulates an internet radio station where different clients can join and leave as they feel. Each client can tune into a specific station and leave as they please. 

### Server Design
There are a few data structures that help the server function. The first data structure is a map which holds connection structs. A connection struct is the tcp and udp connections of a client as well as other fields that help describe the current state of a client. When a client joins the server, the server first establishes a tcp connection and sets it in a map. Each client has a go routine designated to itself which manages any incoming requests. Another part of the server is the radio struct which manages the different stations. Stations manage which song is currently playing. When a client decides to listen to a station, the client will "subscribe" to the station. Essentially, each station contains its own map of client connections and a subscriber struct. A subscriber struct consists of a udp connection as well as a changeSong channel and an endstation channel. When a client is streaming a song, the station publishes the data to all subscribed clients by writing to each udp connection the current data pulled from the song. The EndStation channel is alerted when a station leaves and the changeSong channel is alerted when a new song is going to play. To protect shared access, RWMutexes were used to protect maps and channels were used for message passing amongst different data structures. 

### Client Design
The client has a much simpler design as it first makes a handshake and then waits for replys from server. With each entry in the CLI, the client can send a command to the server.

## Building Files

`make server` --> builds the server

`make client` --> builds the client

`make listener` --> builds the listener

`make build` --> builds the snowcast_control, snowcast_server, and snowcast_listener 

`make clean` --> removes old build files

`make test` --> runs current tests

## Extra Credit
Extra commands were implemented on both the client and server side. In order to run the extra credit, the `-e` flag must be present before the other command line arguments. When adding multiple songs to the same station, they must be comma separated, no space in between

#### To Start a Server with multiple songs per station:
`./snowcast_server -e [port] [song1],[song2],[song3],... [song4] [song5] `

Example:
`./snowcast_server -e 8888 ./mp3/tinyfile,./mp3/mediumfile,./mp3/VanillaIce-IceIceBaby.mp3 ./mp3/tinyfile ./mp3/mediumfile`

### Server Commands
`print/p` --> prints a list of the stations and all the clients listening to each station

`help/h` --> prints the help menu

`addStation/a [song1] [song2]...` --> adds a new station to server with [songs...] as music

`removeStation/r [stationNumber]` --> removes station [stationNumber] from radio

### Client Commands

`getsongs [station]` --> gets all the songs that are playing on the station

`playlist [station] [num songs]` --> gets the next num songs that will be played on the station


