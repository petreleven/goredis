package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	SimpleString   = '+'
	SimpleError    = '-'
	SimpleInt      = ':'
	BulkString     = '$'
	Array          = '*'
	CarriageReturn = '\r'
	LineFeed       = '\n'
	CRLF           = "\r\n"
	NULL           = "$-1\r\n"
	OK             = string(SimpleString) + "OK" + CRLF
)

var RESPONSES = map[string]func([]string) string{
	"PING": func([]string) string {
		// log.Println("called ping")
		return "+PONG\r\n"
	},
	"ECHO": func(slice []string) string {
		if len(slice) == 0 || len(slice) == 1 {
			return string(SimpleError) + "ERROR: FEW ARGUMENTS" + CRLF
		} else if len(slice) > 2 {
			return string(SimpleError) + "ERROR: TOO MANY ARGS" + CRLF
		}
		return string(SimpleString) + slice[1] + CRLF
	},
	"SET": func(slice []string) string {
		invalidArgsErr := string(SimpleError) + "ERROR IN NUMBER OF ARGUMENTS:(hint SET hello world or SET hello world px 100)" + CRLF
		argc := len(slice)

		switch argc {
		case 3:
			(*setDB)[slice[1]] = slice[2]
			return OK
		case 5:
			if strings.ToLower(slice[3]) != "px" {
				return invalidArgsErr
			}
			deadline, err := strconv.Atoi(slice[4])
			if err != nil {
				return "ERROR EXPIRY SHOULD BE INT"
			}
			(*setEXDB)[slice[1]] = ValExp{slice[2], int64(deadline)}
			return OK
		default:
			return invalidArgsErr
		}
	},
	"GET": func(slice []string) string {
		argc := len(slice)
		if argc != 2 {
			return string(SimpleError) + "ERROR IN NUMBER OF ARGUMENTS:(hint GET hello)" + CRLF
		}
		// CHECK CONFIG
		if slice[1] == "dir" {
			reply := fmt.Sprintf("%c%d%s%c%d%s%s%s%c%d%s%s%s",
				Array, 2, CRLF,
				BulkString, 3, CRLF, "dir", CRLF,
				BulkString, len(cfg.Dir), CRLF, cfg.Dir, CRLF)
			return reply
		}
		// CHECK SETDB
		valueSetDB := (*setDB)[slice[1]]
		if valueSetDB != "" {
			return string(SimpleString) + valueSetDB + CRLF
		}
		// CHECK SETEXDB
		valueSetEXDB := (*setEXDB)[slice[1]]
		if valueSetEXDB.Val != "" {
			if valueSetEXDB.Exp > 0 {
				return string(SimpleString) + valueSetEXDB.Val + CRLF
			} else {
				(*setEXDB)[slice[1]] = ValExp{}
			}
		}
		return NULL
	},
}

func Parse(conn *net.Conn, data []byte) string {
	commandsAndArgs := getCommandsAndArgs(data)
	command := strings.ToUpper(commandsAndArgs[0])
	response := RESPONSES[command](commandsAndArgs)

	return response
}

func getCommandsAndArgs(data []byte) []string {
	//*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	var commands []string
	stringData := strings.Split(string(data), "\n")
	for index, val := range stringData {
		stringData[index] = strings.Trim(val, "\r")
	}
	var dtype byte
	if len(stringData) > 0 {
		dtype = stringData[0][0]
	}
	switch dtype {
	case SimpleString:
		if len(stringData) > 1 {
			commands = stringData[1:]
		}
	case Array:
		l := stringData[0][1:]
		lengthOfArr, err := strconv.Atoi(l)
		if err != nil {
			break
		}
		start := 1
		for range lengthOfArr {
			// skip current as its type of data
			start++
			if start < len(stringData) {
				commands = append(commands, stringData[start])
				start++
			}
		}

	}
	return commands
}

func ReduceKeyTTL() {
	for {
		time.Sleep(1 * time.Second)
		for key, v := range *setEXDB {
			(*setEXDB)[key] = ValExp{v.Val, v.Exp - 1}
		}
	}
}
