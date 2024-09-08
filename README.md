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
- [x] RSA or Elliptic-curve based digital signature for verifications
- [ ] GUI usage instead of CLI

## Stack
- Go
- MongoDB

## Usage
  
## Considerations
- HTTP based room messaging can be replaced by a custom protocol such that peers can communicate over a small layer on top of TCP directly.
- For production usage, as centralized lookup server can be developed for peers to connect each other over the web.
- Old-school password based authentication requires an additional layer of encryption. To improve this point, authentication must be done with user signatures where peers validate each others signatures and hashes. This means every user should generate their own private keys and use them for communication with other peers.
In this flow:

1. A user passes a public key and a signature to a peer.

2. The peer can validate the signature with the public key and generate a authentication token for the user. To pass this token, the peer can encrypt it with the users public key.

3. The user can decrypt the encrypted authentication key with his private key and start sending requests to the peer.

4. The user can send room and message requests and the peer encrypt the room secret keys to the user which is similar to the authentication key encryption. 

5. The user can decrypt the room secret key (symmetric key) and use this key to encrypt and decrypt the messages in the room.

6. Master peer can choose the renew the room master key and it can send a synchronization requests to the all users in a room. This would be easier to handle with UDP or WebSocket protocols.
