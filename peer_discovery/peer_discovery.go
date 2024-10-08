package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"log"
	"net"
	"time"
)

func main() {
	var privateKey *ecdsa.PrivateKey

	println("Started !!")
	n, err := node.New(&node.Config{}) //n stores the node instance
	if err != nil {
		log.Printf("Errorr occured while creating node instance☠️: %v", err)
	}
	// fmt.Println(n)

	//Start the node
	if err := n.Start(); err != nil {
		log.Printf("Error occured while starting node: %v", err)
	}

	defer n.Close() //ensures node will stops gracefully when main function ends up

	//=================creating a fileDB using enodeDB==============//
	db, err := enode.OpenDB("node.db")
	if err != nil {
		log.Printf("Error occured while settingup DB:%v", err)
	}

	defer db.Close()
	//==================creating localnode===================//
	localNode := enode.NewLocalNode(db, privateKey)

	//====================================================//
	//setup discovery parameters (bootnodes which helps to find other peer nodes)
	bootnodes := []string{
		"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303", // bootnode-aws-ap-southeast-1-001
		"enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",   // bootnode-aws-us-east-1-001
		"enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303", // bootnode-hetzner-hel
		"enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303", // bootnode-hetzner-fsn
	}

	//convert bootnodes into enode.Node structure
	nodes := make([]*enode.Node, len(bootnodes))

	for i, bn := range bootnodes {
		nd, err := enode.Parse(enode.ValidSchemes, bn)
		if err != nil {
			log.Printf("Error occured while parsing the bootstrap node: %v", err)
		}
		//append the parsed pointer to the nodes slice
		nodes[i] = nd
		// println(nd)
	}

	//setting up udp discovery session
	udpAddr := &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 30303,
	}
	//create a udp connection
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatalf("Error occured whie creating the udp connection:%v", err)
	}

	//configuring the discovery protocol and creating v4 instance
	cfg := discover.Config{
		PrivateKey: privateKey,
		Bootnodes:  nodes,
	}

	udp, err := discover.ListenV4(udpConn, localNode, cfg)
	if err != nil {
		log.Fatalf("Error occured while creating v4 instance:%v", err)
	}

	//ping the bootnodes to start discovering the peers
	go func() {
		for {
			for _, bn := range nodes {
				log.Printf("Pining the nodes:%s", bn.String())
				if err := udp.Ping(bn); err != nil {
					log.Printf("Failed to ping this node%s:%v", bn, err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	//logs new peers
	go func() {
		knownNodes := make(map[string]struct{})

		for {
			time.Sleep(10 * time.Second)
			buckets := udp.TableBuckets()
			for _, bucket := range buckets {
				for _, bucketNode := range bucket {
					node := bucketNode.Node
					nodeId := node.ID().String()
					if _, exists := knownNodes[nodeId]; !exists {
						knownNodes[nodeId] = struct{}{}
						log.Printf("Discovered new peer:%s\n", node.String())
					}

				}
			}
		}
	}()
	select {}

}

/*
var MainnetBootnodes = []string{
	Ethereum Foundation Go Bootnodes

	"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303", // bootnode-aws-ap-southeast-1-001
	"enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",   // bootnode-aws-us-east-1-001
	"enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303", // bootnode-hetzner-hel
	"enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@x:30303", // bootnode-hetzner-fsn
}
*/

//telnet 18.138.108.67 30303
