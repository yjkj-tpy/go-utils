package net_utils

import (
	"fmt"
	"net"
	"os"
	"time"
)

type TcpReceiver struct {
	addr       *net.TCPAddr
	heartbeat  string
	bufferSize int
	dataChan   chan []byte
}

func NewTcpReceiver(port, size int, heartbeat string) *TcpReceiver {
	addr := net.TCPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"), // 监听所有可用的接口
	}

	return &TcpReceiver{
		addr:       &addr,
		heartbeat:  heartbeat,
		bufferSize: size,
		dataChan:   make(chan []byte, size),
	}
}

func (receiver *TcpReceiver) Start() {
	listener, err := net.ListenTCP("tcp", receiver.addr)
	if err != nil {
		fmt.Println("无法监听 TCP:", err)
		os.Exit(1)
	}
	defer listener.Close() // 确保在退出时关闭监听器

	fmt.Println("TCP 服务器在", receiver.addr, "监听中...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接时出错:", err)
			continue
		}

		// 启动一个 goroutine 处理连接
		go receiver.handleConnection(conn)
	}
}

func (receiver *TcpReceiver) GetDataChan() chan []byte {
	return receiver.dataChan
}

func (receiver *TcpReceiver) handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("关闭TCP时出错:", err)
		}
	}(conn) // 确保在处理结束时关闭连接

	buffer := make([]byte, receiver.bufferSize) // 创建一个缓冲区

	// 启动心跳检测
	go receiver.startHeartbeat(conn)

	for {
		n, err := conn.Read(buffer) // 逐字节读取数据
		if err != nil {
			fmt.Println("读取数据时出错:", err)
			return
		}

		receiver.dataChan <- buffer[:n]
	}
}

func (receiver *TcpReceiver) startHeartbeat(conn net.Conn) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := fmt.Fprintln(conn, receiver.heartbeat)
			if err != nil {
				fmt.Println("发送心跳时出错:", err)
				return
			}
		}
	}
}
