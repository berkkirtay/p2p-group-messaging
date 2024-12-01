# group-messaging-in-go
I have been developing this program to utilize a messaging interface for my own local networks. I used P2P approach on top of HTTP. This lets us treat this program as both a peer and a centralized server. How we use it merely depends on the use case. It is possible to use a CLI interface to send a message, as well as to use an HTTP tool such as cURL. It must be noted that clients are responsible to encrypt their own data as well as agreeing on key exchange algorithms in P2P programs and that is what I did in this implementation. 

Right now this is only a fun side project. Please contact me for any consideration and improvement point. As I mentioned earlier, I will continue developing this project as long as I have time.

## Simple flow of the program
1. A requestor peer passes a public key and a signature to serving peer.

2. The serving peer validates the signature with the public key of the requestor peer by RSA algorithm and initializes a new session and authentication key for the peer. Then both peers agree on a new key by using ECDH and the serving peer passes the encrypted authentication key to the requestor.
   
4. In every new request, peers agree on a new key by using newly generated public keys and use this key to encrypt the generated authentication key.

5. Key agreement in every transaction can be done both RSA and ECDH algorithms. In RSA however, serving peer decides the key and passes to the requestor.

6. The requestor peer can now send room and message requests. And can use the exchanged key to encrypt the messages.

7. Every room can have a master key to encrypt the messages and this key can be distributed to the member peers of the room securely.

8. Master peer can choose the renew the room master key and it can send a synchronization requests to the all users in a room. This would be easier to handle with WebSocket protocols (Not done yet).

## Features 
- [x] P2P local node lookup using UDP Multicast 
- [x] ECDH based user and peer authentication
- [x] Chat room authorization
- [x] CBC AES encryption for messages with the ECDH exchanged keys
- [x] RSA based digital signature usage for verifications
- [x] HTTP based async messagging between peers by polling.
- [ ] Websocket based async messaging between peers.
- [ ] Votalite memory usage for data.

## Stack
- Go
- Go Gin framework
- PKCS libraries
- MongoDB

## Usage
  
## Considerations
- HTTP based room messaging can be replaced by a custom protocol such that peers can communicate over a small layer on top of TCP directly.
- For production usage, as centralized lookup server can be developed for peers to connect each other over the web.
- Peers agree on a new exchanged key in every new transaction. But this key is only used for encrypting text/data field in the transmissions. The whole transmitted data can be encrypted as well.
