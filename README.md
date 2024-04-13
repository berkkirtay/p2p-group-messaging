# group-messaging-in-go
## Features 
- P2P node lookup (requires improvement)
- HMAC based user authentication
- Chat room authorization
- HTTP based async messagging with DB synchronization among peers (requires improvement)
- Diffie-hellman key exchange for p2p key agreement (requires improvement)
- CBC AES encryption for messages
- RSA based digital signature for verifications (TODO)

## Considerations
- HTTP based room messaging can be replaced by a custom protocol such that peers can communicate on a small layer on top of TCP directly.