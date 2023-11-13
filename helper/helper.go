package helper

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	commandDelimiter   = " "    // 命令分隔符号
	commandPrefix      = "*"    // 命令前缀
	commandParamPrefix = "$"    // 命令参数前缀
	commandSuffix      = "\r\n" // 命令后缀
)

/**
*<参数数量> CR LF
$<参数 1 的字节数量> CR LF
<参数 1 的数据> CR LF
...
$<参数 N 的字节数量> CR LF
<参数 N 的数据> CR LF
*/

func GenerateRequest(command string) []byte {
	strSlice := strings.Split(command, commandDelimiter)
	commandSlice := make([]string, len(strSlice)*2+1)
	
	commandSlice[0] = commandPrefix + strconv.Itoa(len(strSlice))
	
	for index, str := range strSlice {
		commandSlice[index*2+1] = commandParamPrefix + strconv.Itoa(len(str))
		commandSlice[index*2+2] = str
	}
	
	return []byte(strings.Join(commandSlice, commandSuffix) + commandSuffix)
}

func SendRequest(conn net.Conn, command string) ([]byte, error) {
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	
	done := make(chan struct{})
	var (
		result []byte
		err    error
	)
	go func() {
		result, err = sendRequest(conn, command, done)
	}()
	
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-done:
			return result, err
		}
	}
}

func sendRequest(conn net.Conn, command string, done chan struct{}) ([]byte, error) {
	_, err := conn.Write(GenerateRequest(command))
	if err != nil {
		return nil, err
	}
	
	resp := make([]byte, 0)
	buffer := make([]byte, 1024)
	
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		
		resp = append(resp, buffer[:n]...)
		
		if n < 1024 {
			break
		}
	}
	
	defer func() {
		done <- struct{}{}
	}()
	
	return GenerateResponse(resp)
}

func SendRequestV2(conn net.Conn, command string) ([]byte, error) {
	redisCommand := GenerateRequest(command)
	n, err := conn.Write(redisCommand)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Command len: %d, Command: %#v\n", n, string(redisCommand))
	
	resp := make([]byte, 0)
	buffer := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		
		resp = append(resp, buffer[:n]...)
		
		if n < 1024 {
			break
		}
	}
	
	return GenerateResponse(resp)
}
