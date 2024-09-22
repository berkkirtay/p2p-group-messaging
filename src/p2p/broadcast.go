// Copyright (c) 2024 Berk Kirtay

package p2p

import (
	"fmt"
	"main/commands"
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
	LOCAL_BROADCAST_ADDRESS string = "224.0.0.1:9999"
	NETWORK_NAME            string = "udp4"
	MAX_BROADCAST_AMOUNT    int    = 2
	ADDRESS_FORMAT          string = HTTP + "%s" + PORT + API
)

var unique_udp_request_identifier int = 0
var currentHost string = ""

func HandlePeerConnection() {
	initializeOutboundIP()
	generateNextConnectionIdentifier()
	go listenForPeerBroadcast()
	startPeerBroadcast()
}

func startPeerBroadcast() {
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
		if commands.IsPeerInitialized() {
			return
		}
		fmt.Printf("Broadcasting for an active peer... %d\n", i)
		_, err = broadcast.Write([]byte(strconv.Itoa(unique_udp_request_identifier)))
		if err != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
	}

	// Become a main peer in case no other peer exists:
	fmt.Println("No active peer found, making yourself an active peer.")
	commands.InitializeAMasterPeer(
		currentHost,
		fmt.Sprintf(ADDRESS_FORMAT, currentHost))
}

func listenForPeerBroadcast() {
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
			var remote_udp_request_identifier, _ = strconv.Atoi(string(buf[:n]))
			if remote_udp_request_identifier != unique_udp_request_identifier {
				var targetAddress string = fmt.Sprintf(
					ADDRESS_FORMAT,
					strings.Split(addr.String(),
						":")[0])
				commands.RegisterPeer(
					targetAddress,
					currentHost,
					fmt.Sprintf(ADDRESS_FORMAT, currentHost))
			}
		}
	}
}

func initializeOutboundIP() {
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
	unique_udp_request_identifier = int(rand.Int31())
}
