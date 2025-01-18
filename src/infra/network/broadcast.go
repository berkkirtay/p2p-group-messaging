// Copyright (c) 2024 Berk Kirtay

package network

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/ipv4"
)

/*
 * UDP multicast peer broadcast implementation based on RFC 1301 and golang.org/x/net/ipv4.
 *
 * Provides fast peer lookup in any local network. Basically, when a listener peer
 * receives a broadcast signal from another peer, the listener peer will send a peer
 * HTTP request to the broadcaster peer. If the request is successful, peers store
 * each others addresses with INBOUND and OUTBOUND information. At last both peers
 * can proceed with handshaking upon the initial broadcasting connection.
 */

const (
	CUSTOM_PEER_LOOKUP_API  string = "http://127.0.0.1:8080/api/peer"
	LOCAL_BROADCAST_ADDRESS string = "224.0.0.1:9999"
	NETWORK_NAME            string = "udp4"
	MAX_BROADCAST_AMOUNT    int    = 3
	PORT                    string = ":8080"
	API                     string = "/api"
	HTTP                           = "http://"
	ADDRESS_FORMAT          string = HTTP + "%s" + "%s" + API
)

var unique_udp_request_identifier string = ""
var currentHost string = ""
var currentHostPort = PORT

func InitializeBroadcast(identifier string) {
	if identifier != "" {
		currentHostPort = identifier
	}
	initializeHostOutboundAddress()
	generateNextConnectionIdentifier()
}

func StartPeerBroadcast() {
	remote, err := net.ResolveUDPAddr(NETWORK_NAME, LOCAL_BROADCAST_ADDRESS)
	if err != nil {
		panic(err)
	}
	broadcast, err := net.DialUDP(NETWORK_NAME, nil, remote)
	if err != nil {
		panic(err)
	}
	defer broadcast.Close()

	// Continuously check if the current host has established a connection with any peer:
	for i := 0; i < MAX_BROADCAST_AMOUNT; i++ {
		fmt.Printf("Broadcasting for active peers... %d\n", i)
		_, err = broadcast.Write([]byte(unique_udp_request_identifier))
		if err != nil {
			panic(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

/*
 * Listens for the next peer connection and returns @targetAddress
 */
func ListenForPeerBroadcast() (string, string) {
	ipv4Addr, _ := net.ResolveUDPAddr(NETWORK_NAME, LOCAL_BROADCAST_ADDRESS)
	conn, err := net.ListenUDP(NETWORK_NAME, ipv4Addr)
	if err != nil {
		panic(err)
	}

	pc := ipv4.NewPacketConn(conn)
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	err = pc.JoinGroup(&interfaces[0], ipv4Addr)
	if err != nil {
		panic(err)
	}
	loop, err := pc.MulticastLoopback()
	if err == nil {
		if !loop {
			pc.SetMulticastLoopback(true)
		}
	} else {
		panic(err)
	}
	buf := make([]byte, 1024*4)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err == nil {
			// Send your peer info to the sender and conclude broadcasting:
			var remote_udp_request_identifier = string(buf[:n])
			if remote_udp_request_identifier != unique_udp_request_identifier {
				return remote_udp_request_identifier, fmt.Sprintf(
					ADDRESS_FORMAT,
					strings.Split(addr.String(),
						":")[0],
					":"+strings.Split(remote_udp_request_identifier, ":")[1])
			}
		}
	}
}

func GetHostAddress() (string, string) {
	return unique_udp_request_identifier, fmt.Sprintf(ADDRESS_FORMAT, currentHost, currentHostPort)
}

func initializeHostOutboundAddress() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		ipnet := address.(*net.IPNet)
		if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			currentHost = ipnet.IP.String()
			return
		}
	}
}

/*
 * Assings @unique_udp_request_identifier with a random generated
 * value to avoid peer collisions on local networks.
 */
func generateNextConnectionIdentifier() {
	unique_udp_request_identifier = strconv.Itoa(int(rand.Int31())) + currentHostPort
}
