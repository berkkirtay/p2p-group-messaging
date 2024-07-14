# group-messaging-in-go
I have been developing this program to utilize a messaging interface for my own local networks. I used P2P approach on top of HTTP. This lets us treat this program as both a peer and a centralized server. How we use it merely depends on the use case. It is possible to use a CLI interface to send a message, as well as to use an HTTP tool such as cURL. It must be noted that clients are responsible to encrypt their own data in P2P programs and that is what I did in this implementation. Surely we can use HTTPS instead of HTTP. This is a must if this program is to be used in a production environment. I plan to add this along with Diffie-hellman key exchange between peers.

Right now this is only a fun side project. Please contact me for any consideration and improvement point. As I mentioned earlier, I will continue developing this project as long as I have time.

## Features 
- [x] P2P local node lookup using UDP Multicast 
- [x] HMAC based user authentication
- [x] Chat room authorization
- [x] HTTP based async messagging with DB synchronization among peers (requires improvement)
- [ ] Diffie-hellman key exchange for p2p key agreement (requires improvement)
- [x] CBC AES encryption for messages
- [ ] RSA or Elliptic-curve based digital signature for verifications
- [ ] GUI usage instead of CLI

## Stack
- Go
- MongoDB

## Usage
  
## Considerations
- HTTP based room messaging can be replaced by a custom protocol such that peers can communicate over a small layer on top of TCP directly.
- For production usage, as centralized lookup server can be developed for peers to connect each other over the web.
