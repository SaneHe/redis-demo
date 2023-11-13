package main

import (
	"fmt"
	"net"
	"redis-protocol/helper"
)

const (
	address = "192.168.0.105:6379"
	network = "tcp"
)

func main() {
	redisConn, err := helper.Connect(network, address)
	if err != nil {
		panic(err)
	}
	defer redisConn.Close()
	
	bulkReply(redisConn)
	errorReply(redisConn)
	statusReply(redisConn)
	integerReply(redisConn)
	multiBulkReply(redisConn)
}

func bulkReply(redisConn net.Conn) {
	command := "get sane"
	resp, err := helper.SendRequestV2(redisConn, command)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Command: %#v, Raw Response: %#v\n", command, string(resp))
}

func errorReply(redisConn net.Conn) {
	command := "get2 sane"
	resp, err := helper.SendRequestV2(redisConn, command)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Command: %#v, Raw Response: %#v\n", command, string(resp))
}

func statusReply(redisConn net.Conn) {
	// 若后面增加参数，则为 bulkReply
	command := "ping"
	resp, err := helper.SendRequestV2(redisConn, command)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Command: %#v, Raw Response: %#v\n", command, string(resp))
}

func integerReply(redisConn net.Conn) {
	command := "incr sane"
	resp, err := helper.SendRequestV2(redisConn, command)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Command: %#v, Raw Response: %#v\n", command, string(resp))
}

func multiBulkReply(redisConn net.Conn) {
	command := "mget sane foo foo2 foo22"
	resp, err := helper.SendRequestV2(redisConn, command)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Command: %#v, Raw Response: %#v\n", command, string(resp))
}
