# Crawl
# Crawl-Golang
												
Nộp bài trên github.com để chế độ private (hoặc internal) sau đó mời account @haicon232 (hoặc email haicon2321993@gmail.com) vào để check												
Đọc bài sau để xem VD của http request và response: https://medium.com/@masnun/making-http-requests-in-golang-dd123379efe7												
Mô hình sẽ như sau:												
1. Request đến https://malshare.com/daily --> Nếu malshare bị lỗi (thỉnh thoảng website chết thì crawl từ virusshare https://virusshare.com/hashes/)
2. Đọc Response (Response sẽ chính là text HTML - bản chất là string)
3. Duyệt string để lấy các thông tin cần thiết (đường link đến các ngày) --> Có thể sử dụng regex để việc lấy các thông tin cần thiết đơn giản.
4. Sau khi lấy được link đến từng ngày thì lại Request đến link đó (có dạng https://malshare.com/daily/yyyy-MM-dd/malshare_fileList.yyyy-MM-dd.all.txt)
Nhận Response là các md5, sha1, sha256, ....
5. Lưu các thông tin cần có lại thành 1 map: ngày --> md5, ngày --> sha1, ....
6. Duyệt map vừa lưu, tạo các thư mục tương ứng theo ngày, tháng, năm và viết file.

*Note: Các bước 5-6 có thể thay đổi theo logic của người code. Đây chỉ là VD để hiểu rõ hơn"										
Vào link sau: https://malshare.com/daily/												
Sẽ thấy các thư mục theo ngày, trong mỗi thư mục đó, tìm đến file có dạng .all.txt												
Trong file đó, nội dung sẽ như sau:												
md5 sha1 sha256 ...*												
								
### Nhiệm vụ:												
1. Crawl tất cả các md5, sha1, sha256 của trang malshare. 	
	Lưu các thông tin crawl được vào mongo DB. Tự định nghĩa các trường cần thiết để lưu lại								
									
2.	Tiếp tục sử dụng Task 3								
	Viết RESTFUL API cho các data có sẵn trong Database.					
	Tự định nghĩa các đầu API để xem về khả năng phân tích data. (Có data rồi thì cần viết API gì?)	