package utils

import "net"

// GetFreePort 动态获取一个空闲的端口
func GetFreePort() (int, error) {
	// 其中，"0" 代表让系统自动分配一个空闲端口号。
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	// 如果地址解析成功，调用 ListenTCP 函数，用于创建并监听一个 TCP 网络连接。
	// 其中，"tcp" 代表协议类型，addr 则是上一步解析出的 TCP 地址结构。
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
