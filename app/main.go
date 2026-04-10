package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		os.Exit(1)
	}
	defer l.Close()

	// Danh sách tất cả khách hàng đang kết nối
	var clients []net.Conn

	fmt.Println("Event Loop started...")

	for {
		// BƯỚC 1: KIỂM TRA KHÁCH MỚI (NON-BLOCKING ACCEPT)
		// Ép l.Accept() chỉ chờ trong 1ms, nếu không có ai thì chạy tiếp luôn
		l.(*net.TCPListener).SetDeadline(time.Now().Add(time.Millisecond))
		conn, err := l.Accept()
		if err == nil {
			fmt.Println("New client joined!")
			clients = append(clients, conn)
		}

		// BƯỚC 2: DUYỆT QUA TẤT CẢ CÁC KHÁCH ĐANG CÓ
		for i := 0; i < len(clients); i++ {
			client := clients[i]

			// Ép việc Đọc dữ liệu chỉ chờ trong 1ms (Non-blocking Read)
			client.SetReadDeadline(time.Now().Add(time.Millisecond))
			
			buf := make([]byte, 1024)
			_, err := client.Read(buf)

			if err == nil {
				// CÓ DỮ LIỆU: Xử lý ngay lập tức
				client.Write([]byte("+PONG\r\n"))
			} else {
				// KIỂM TRA XEM LÀ LỖI HAY CHỈ LÀ CHƯA CÓ TIN NHẮN
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// Chỉ là khách chưa gửi gì, bỏ qua và sang khách tiếp theo
					continue
				} else {
					// Lỗi thật hoặc khách đã ngắt kết nối
					fmt.Println("Client left.")
					client.Close()
					// Xóa khách khỏi danh sách (Event Loop)
					clients = append(clients[:i], clients[i+1:]...)
					i-- // Lùi chỉ số lại vì mảng vừa bị co lại
				}
			}
		}

		// Nghỉ 1 chút để CPU không bị quá nóng (vì vòng lặp for chạy quá nhanh)
		time.Sleep(10 * time.Millisecond)
	}
}