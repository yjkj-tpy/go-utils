package net_utils

import (
	"fmt"
	"net"
	"os"
)

type UdpReceiver struct {
	addr     *net.UDPAddr
	dataChan chan []byte
}

func NewUdpReceiver(port, size int) *UdpReceiver {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"), // 监听所有可用的接口
	}

	return &UdpReceiver{
		addr:     &addr,
		dataChan: make(chan []byte, size),
	}
}

func (receiver *UdpReceiver) GetDataChan() chan []byte {
	return receiver.dataChan
}

func (receiver *UdpReceiver) Start() {
	conn, err := net.ListenUDP("udp", receiver.addr)
	if err != nil {
		fmt.Println("无法监听 UDP:", err)
		os.Exit(1)
	}
	defer func(conn *net.UDPConn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭UDP时出错:", err)
		}
	}(conn) // 确保在退出时关闭连接

	fmt.Println("UDP 服务器在", receiver.addr.String(), "监听中...")

	buffer := make([]byte, 10240)

	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("读取数据时出错:", err)
			continue
		}

		if n > 0 {
			data := make([]byte, n)
			copy(data, buffer[:n])
			receiver.dataChan <- data
		}
	}
}
