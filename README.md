# ⚙️ Expense Management - Go Backend API

Hệ thống Backend cung cấp RESTful API chất lượng cao cho ứng dụng quản lý chi tiêu cá nhân. Được xây dựng bằng ngôn ngữ **Go (Golang)** tối ưu hiệu năng, kiến trúc phân lớp sạch sẽ, và được cấu hình sẵn môi trường Docker hóa tiện lợi.

---

## 🛠️ Công Nghệ Sử Dụng

- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin) - Framework nhanh và mạnh mẽ nhất cho Go.
- **ORM (Object-Relational Mapping)**: [GORM](https://gorm.io/) - Hỗ trợ tương tác cơ sở dữ liệu dễ dàng, an toàn, và tự động Migration.
- **Cơ sở dữ liệu**: [PostgreSQL](https://www.postgresql.org/) phiên bản 16-alpine.
- **Bảo mật**: JWT (JSON Web Tokens) cho phân quyền API, băm mật khẩu bằng `bcrypt`, băm mã PIN bằng mã hóa an toàn.
- **Tài liệu API**: [swaggo/swag](https://github.com/swaggo/swag) tự động tạo tài liệu đặc tả OpenAPI và tích hợp giao diện tương tác Swagger UI.
- **Containerization**: Docker & Docker Compose để khởi chạy nhanh toàn bộ môi trường (App & Database).

---

## 📁 Cấu Trúc Thư Mục Dự Án (Clean Layered Architecture)

Backend được tổ chức theo cấu trúc phân lớp logic (Layered Architecture) giúp cô lập mã nguồn nghiệp vụ:

```text
expense_management_backend/
├── cmd/
│   └── main.go                 # Điểm khởi chạy ứng dụng (Entrypoint)
├── docs/                       # Tài liệu Swagger API phục vụ Client
├── config/                     # Đọc và khởi tạo các cấu hình (config.go)
├── internal/                   # Chứa toàn bộ mã nguồn nghiệp vụ nội bộ
│   ├── config/                 # Định nghĩa và nạp biến môi trường (.env)
│   ├── controllers/            # Lớp Handler tiếp nhận HTTP Request & gửi Response
│   │   ├── auth/               # Controller quản lý Đăng nhập, Đăng ký, OAuth Google
│   │   └── user/               # Controller quản lý Hồ sơ, Avatar, API PIN Bảo mật
│   ├── db/                     # Cấu hình kết nối PostgreSQL & GORM
│   ├── docs/                   # Tài liệu Swagger phục vụ hiển thị trên Swagger UI
│   ├── middleware/             # Các bộ lọc JWT Auth, CORS, Logger...
│   ├── models/                 # Định nghĩa các struct và lược đồ cơ sở dữ liệu (GORM)
│   ├── repositories/           # Tương tác trực tiếp với Database qua GORM
│   ├── services/               # Lớp trung gian chứa Logic xử lý nghiệp vụ chính
│   └── router/                 # Định nghĩa và phân nhóm các API endpoints
├── Dockerfile                  # Cấu hình đóng gói ứng dụng Go
├── docker-compose.yml          # Cấu hình chạy cụm Container (Go App + PostgreSQL)
└── go.mod                      # Quản lý Golang Dependencies
```

---

## 📋 Cấu Hình Biến Môi Trường (Environment Variables)

Hãy copy file `.env.example` thành file `.env` tại thư mục gốc và tùy chỉnh các tham số dưới đây:

### 1. Application Config
- `APP_PORT`: Cổng khởi chạy máy chủ Go (mặc định: `8080`).
- `GIN_MODE`: Chế độ chạy của Gin Framework (`debug` hoặc `release`).

### 2. Database Config (PostgreSQL)
- `DB_HOST`: Địa chỉ IP/Domain của PostgreSQL (chọn `db` nếu chạy Docker, `localhost` nếu chạy máy thật).
- `DB_PORT`: Cổng PostgreSQL (mặc định: `5432`).
- `DB_USER`: Tên tài khoản quản trị DB (ví dụ: `expense_user`).
- `DB_PASSWORD`: Mật khẩu tài khoản DB (ví dụ: `expense_secret`).
- `DB_NAME`: Tên database lưu trữ dữ liệu (ví dụ: `expense_management`).
- `DB_SSLMODE`: Cấu hình SSL (`disable` hoặc `require`).
- `DB_TIMEZONE`: Múi giờ lưu dữ liệu (`Asia/Ho_Chi_Minh`).

### 3. SMTP Config (Gửi Mail khôi phục mật khẩu)
- `SMTP_HOST`: Địa chỉ SMTP Server (ví dụ: `smtp.gmail.com`).
- `SMTP_PORT`: Cổng kết nối SMTP (thường là `587` cho TLS hoặc `465` cho SSL).
- `SMTP_USER`: Tài khoản email gửi tin (ví dụ: `your-email@gmail.com`).
- `SMTP_PASSWORD`: Mật khẩu ứng dụng (App Password) của email.
- `SMTP_FROM`: Nhãn hiển thị người gửi (ví dụ: `no-reply@expensemanagement.com`).

---

## 🛡️ Cơ Chế Bảo Mật & Xác Thực

### 1. JWT Middleware Authorization Flow
Hệ thống sử dụng Authorization Header định dạng Bearer Token để bảo vệ các tài nguyên API riêng tư:
* Client gửi yêu cầu với header: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
* Middleware [auth.go](file:///home/quan/Documents/expense_management_backend/internal/middleware/auth.go) sẽ chặn request, giải mã signature của token:
  * Nếu token hợp lệ: Giải nén thông tin `user_id` & `email` lưu vào Context của Gin và tiếp tục gọi Controller xử lý.
  * Nếu token hết hạn/không hợp lệ: Trả về HTTP `401 Unauthorized` ngay lập tức.

### 2. Bảo mật mã PIN & Câu hỏi khôi phục
Để tránh rủi ro rò rỉ dữ liệu nhạy cảm của người dùng:
- **Mã PIN**: Được băm một chiều an toàn bằng thuật toán **`bcrypt`** (độ phức tạp cao) trước khi lưu vào cột `pin_hash` trong DB.
- **Câu trả lời bảo mật (Security Answer)**: Trước khi băm, chuỗi câu trả lời sẽ được **chuẩn hóa** (loại bỏ khoảng trắng thừa ở hai đầu, chuyển tất cả ký tự thành chữ thường không dấu) để tăng tỷ lệ nhập đúng của người dùng, sau đó mới băm bằng `bcrypt` và lưu vào cột `security_answer_hash`.

### 3. Database Auto-Migration
Mỗi lần Go Backend khởi chạy, lớp cơ sở dữ liệu của GORM sẽ kích hoạt tính năng tự động dò tìm cấu trúc DB:
- Tự động phát hiện các model mới và ánh xạ thành các bảng tương ứng trong PostgreSQL.
- Tự động thêm các cột mới (ví dụ: `pin_hash`, `security_question`, `security_answer_hash`) vào bảng `users` mà **hoàn toàn không làm ảnh hưởng hay mất mát dữ liệu hiện có** của người dùng.

---

## 🔌 Đặc Tả Hệ Thống API PIN Bảo Mật (App Lock APIs)

Các API quản lý khóa và PIN nằm dưới prefix `/api/v1/profile` và đều yêu cầu Header Authenticate:

| Phương thức | Đường dẫn API | Mô tả tính năng | Dữ liệu đầu vào (JSON Body) |
|:---:|---|---|---|
| **POST** | `/api/v1/profile/pin` | Kích hoạt mã PIN và Câu hỏi bảo mật | `{ "pin": "1234", "security_question": "Tên thú cưng?", "security_answer": "lulu" }` |
| **POST** | `/api/v1/profile/pin/verify` | Xác thực mã PIN khi mở/quần ứng dụng | `{ "pin": "1234" }` |
| **GET** | `/api/v1/profile/pin/security-question` | Lấy câu hỏi bảo mật để hiển thị khi khôi phục PIN | *Không có Body* |
| **POST** | `/api/v1/profile/pin/reset` | Đặt mã PIN mới bằng câu trả lời câu hỏi bảo mật | `{ "security_answer": "lulu", "new_pin": "5678" }` |
| **DELETE** | `/api/v1/profile/pin` | Vô hiệu hóa PIN và khóa ứng dụng | `{ "pin": "1234" }` |

---

## 🚀 Hướng Dẫn Cài Đặt & Khởi Chạy Nhanh

### Cách 1: Khởi chạy nhanh bằng Docker Compose (Khuyên dùng)
Bạn không cần cài đặt Go hay PostgreSQL trên máy thật, chỉ cần cài sẵn Docker và Docker Desktop.

1. **Khởi chạy hệ thống**:
   ```bash
   docker compose up --build -d
   ```
2. **Kiểm tra trạng thái**:
   ```bash
   docker compose ps
   ```
   Hệ thống sẽ chạy Go App trên cổng `8080` và PostgreSQL trên cổng `5432`.
3. **Dừng hệ thống**:
   ```bash
   docker compose down
   ```

### Cách 2: Chạy trực tiếp trên máy local (Development mode)
Yêu cầu bạn đã cài đặt sẵn Go 1.20+ và một cơ sở dữ liệu PostgreSQL đang chạy.

1. Cấu hình biến môi trường trong file `.env` tại thư mục gốc.
2. Chạy ứng dụng:
   ```bash
   go run cmd/main.go
   ```

---

## 🧪 Hướng dẫn Kiểm thử API (cURL Commands)

Dưới đây là chuỗi lệnh cURL ví dụ để test nhanh luồng API PIN Bảo mật:

### 1. Đăng ký & Kích hoạt PIN:
```bash
curl -X POST http://localhost:8080/api/v1/profile/pin \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"pin": "1234", "security_question": "Trường tiểu học đầu tiên?", "security_answer": "Kim Dong"}'
```

### 2. Xác thực mã PIN:
```bash
curl -X POST http://localhost:8080/api/v1/profile/pin/verify \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"pin": "1234"}'
```

### 3. Lấy câu hỏi bảo mật khi quên PIN:
```bash
curl -X GET http://localhost:8080/api/v1/profile/pin/security-question \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>"
```

### 4. Đặt lại PIN bằng câu trả lời bảo mật:
```bash
curl -X POST http://localhost:8080/api/v1/profile/pin/reset \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"security_answer": "Kim Dong", "new_pin": "9999"}'
```

---

## 📄 Tài Liệu API Swagger UI

Sau khi hệ thống Backend được khởi chạy thành công:
- Truy cập địa chỉ: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) để xem toàn bộ tài liệu API trực quan.
- Bạn có thể thử nghiệm gửi request trực tiếp trên giao diện Swagger UI này để kiểm tra tính đúng đắn của các endpoint.

*Để cập nhật lại tài liệu Swagger sau khi chỉnh sửa code, hãy chạy lệnh:*
```bash
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o internal/docs
go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o docs
```
