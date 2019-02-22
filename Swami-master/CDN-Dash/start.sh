#!/bin/bash
set -x

IP_ADDR=`ifconfig en0 | grep "inet " | awk -F '[: ]+' '{ print $2 }'`
NODE_NUM=$1

case "$1" in
    "1")
        echo `go run node.go :8001 ${IP_ADDR}:8002 ${IP_ADDR}:8003 ${IP_ADDR}:8004 ${IP_ADDR}:8000 ${IP_ADDR} >> output1.log`
        ;;
    "2")
	echo `go run node.go :8002 ${IP_ADDR}:8001 ${IP_ADDR}:8003 ${IP_ADDR}:8004 ${IP_ADDR}:8000 ${IP_ADDR} >> output2.log`
    	;;
    "3")
	echo `go run node.go :8003 ${IP_ADDR}:8001 ${IP_ADDR}:8002 ${IP_ADDR}:8004 ${IP_ADDR}:8000 ${IP_ADDR} >> output3.log`
        ;;
    "4")
	echo `go run node.go :8004 ${IP_ADDR}:8001 ${IP_ADDR}:8002 ${IP_ADDR}:8003 ${IP_ADDR}:8000 ${IP_ADDR} >> output4.log`
	;;
esac

