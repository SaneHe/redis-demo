package helper

import "net"

func Connect(network, address string) (net.Conn, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	
	return conn, nil
}
