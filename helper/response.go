package helper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	// 其中 13 10 代表 \r\n
	// 返回示例: [43 80 79 78 71 13 10]
	statusReply = '+'
	
	// 返回示例: [45 69 82 82 32 117 110 107 110 111 119 110 32 99 111 109 109 97 110 100 32 96 103 101 116 49 96 44 32 119 105 116 104 32 97 114 103 115 32 98 101 103 105 110 110 105 110 103 32 119 105 116 104 58 32 96 102 111 111 50 50 96 44 32 13 10]
	errorReply = '-'
	
	// 返回示例: [58 50 13 10]
	integerReply = ':'
	
	// 返回示例: [36 52 13 10 115 97 110 101 13 10]
	bulkReply = '$'
	
	// 返回示例: [42 50 13 10 36 50 13 10 49 54 13 10 36 52 13 10 115 97 110 101 13 10]
	multiBulkReply = '*'
)

func GenerateResponse(response []byte) ([]byte, error) {
	replayStrategy := response[0]
	
	// 获取第一位标识符的位置
	position := getFirstDelimiter(response, '\r')
	// 打印响应内容 切片
	fmt.Println(response)
	
	switch replayStrategy {
	case bulkReply:
		fmt.Printf("bulkReplyResponse mark: %#v, content: %#v\n", string(replayStrategy), string(response))
		return bulkReplyResponse(response, position)
	case errorReply:
		fmt.Printf("errorReplyResponse mark: %#v, content: %#v\n", string(replayStrategy), string(response))
		return errorReplyResponse(response, position)
	case statusReply:
		fmt.Printf("statusReplyResponse mark: %#v, content: %#v\n", string(replayStrategy), string(response))
		return statusReplyResponse(response, position)
	case integerReply:
		fmt.Printf("integerReplyResponse mark: %#v, content: %#v\n", string(replayStrategy), string(response))
		return integerReplyResponse(response, position)
	case multiBulkReply:
		fmt.Printf("multiBulkReplyResponse mark: %#v, content: %#v\n", string(replayStrategy), string(response))
		return multiBulkReplyResponse(response, position)
	}
	return response, nil
}

/**
$<返回值 的字节数量> CR LF
<返回值 的数据> CR LF
*/
func bulkReplyResponse(response []byte, position int) ([]byte, error) {
	if dataLen, err := strconv.Atoi(string(response[1:position])); err != nil {
		return nil, err
	} else if dataLen < 0 {
		return nil, nil
	}
	
	return response[position+2 : len(response)-2], nil
}

/**
-<错误消息 的数据> CR LF
*/
func errorReplyResponse(response []byte, position int) ([]byte, error) {
	return response[1 : position-2], nil
}

/**
*<返回值 的数据> CR LF
 */
func statusReplyResponse(response []byte, position int) ([]byte, error) {
	return response[1:position], nil
}

/**
:<返回值 的数据> CR LF
*/
func integerReplyResponse(response []byte, position int) ([]byte, error) {
	return response[1:position], nil
}

/**
*<返回值 的数量> CR LF
$<返回值 1 的字节数量> CR LF
<返回值 1 的数据> CR LF

$<返回值 N 的字节数量> CR LF
<返回值 N 的数据> CR LF
*/
func multiBulkReplyResponse(response []byte, position int) ([]byte, error) {
	dataLen, err := strconv.Atoi(string(response[1:position]))
	if err != nil {
		return nil, err
	} else if dataLen < 0 {
		return nil, nil
	}
	
	result, responseSlice := make([]string, dataLen), strings.Split(string(response), commandSuffix)
	for index := range result {
		result[index] = responseSlice[index*2+2]
	}
	
	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func getFirstDelimiter(data []byte, delimiter byte) int {
	for i := 0; i < len(data); i++ {
		if data[i] == delimiter {
			return i
		}
	}
	
	return -1
}
