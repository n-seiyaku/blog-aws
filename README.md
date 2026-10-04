# Blog AWS

Dự án blog gồm backend Go/Gin sử dụng AWS DynamoDB và frontend React/TypeScript chạy bằng Vite.

## Chức năng hiện tại

- Backend: đăng ký, đăng nhập, làm mới token, đăng xuất và kiểm tra trạng thái server.
- Xác thực: mật khẩu được băm bằng bcrypt, access token dùng JWT, refresh token được lưu dưới dạng hash trong session.
- Frontend: trang đăng nhập `/login` và trang blog `/blog` với dữ liệu mẫu.

Frontend chưa kết nối API backend. Đăng nhập giả lập hiện dùng email `123` và mật khẩu `123`; trạng thái đăng nhập chỉ được giữ trong bộ nhớ.

## Cấu trúc

```text
blog-aws/
├── backend/
│   ├── cmd/api/                 # Điểm khởi chạy HTTP server
│   └── internal/
│       ├── adapter/http/handler/ # Xử lý request và response
│       ├── domain/              # Entity và interface repository
│       ├── infrastructure/      # DynamoDB và token
│       └── usecase/             # Nghiệp vụ
└── frontend/
    └── src/
        ├── app/                 # Router và layout
        ├── auth/                # Context xác thực
        └── features/            # Trang đăng nhập và blog
```

## Yêu cầu

- Go 1.27.1 hoặc mới hơn, theo `backend/go.mod`.
- Node.js 20.19+ thuộc nhánh 20, hoặc 22.12+; pnpm để dùng lockfile hiện có.
- AWS credentials và region được cấu hình cho SDK khi chạy backend với DynamoDB.
- Hai bảng DynamoDB cho user và session, được tạo trước khi sử dụng API auth.

Các lệnh bên dưới dùng PowerShell. Chạy backend và frontend trong hai terminal riêng từ thư mục gốc dự án.

## Cấu hình Backend

Tạo file `backend/.env` với các giá trị phù hợp với môi trường của bạn:

```dotenv
JWT_SECRET=replace-with-a-long-random-secret
DYNAMODB_USERS_TABLE=blog-users
DYNAMODB_SESSIONS_TABLE=blog-sessions
AWS_REGION=ap-northeast-1
```

`JWT_SECRET`, `DYNAMODB_USERS_TABLE` và `DYNAMODB_SESSIONS_TABLE` là bắt buộc. AWS SDK sử dụng chuỗi cấu hình mặc định để lấy credentials, ví dụ từ AWS profile hoặc IAM role. Với AWS CLI đã cài đặt, có thể cấu hình profile mặc định bằng `aws configure`; dùng `AWS_PROFILE` nếu cần chọn profile khác.

File `backend/.env` đã được bỏ qua trong Git. Không đưa secret hoặc AWS credentials vào repository.

### Bảng DynamoDB

Tên bảng phải khớp với giá trị trong `.env`. Cấu hình khóa và index như sau:

| Bảng | Partition key | Global secondary index | Partition key của index |
| --- | --- | --- | --- |
| Users | `id` (String) | `email-index` | `email` (String) |
| Sessions | `id` (String) | `refreshTokenHash-index` | `refreshTokenHash` (String) |

Hai bảng và các index không cần sort key. Chọn projection `ALL` cho các index vì repository đọc dữ liệu user/session trực tiếp từ kết quả query. Có thể bật TTL trên thuộc tính `ttl` (Number) của bảng Sessions để tự động dọn session hết hạn.

Credentials của backend cần quyền `GetItem`, `PutItem`, `Query`, `UpdateItem` và `DeleteItem` tương ứng với bảng và index được sử dụng. Server không tự tạo bảng.

### Chạy Server

```powershell
cd backend
go mod download
go run ./cmd/api
```

Backend lắng nghe tại `http://localhost:8080`. Chạy từ thư mục `backend` để ứng dụng đọc đúng file `.env`.

Kiểm tra server từ terminal khác:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/health
```

Response mong đợi: `{"status":"ok"}`. Endpoint này chỉ kiểm tra HTTP server, không kiểm tra kết nối DynamoDB.

## API

Các endpoint auth nhận JSON với header `Content-Type: application/json`.

| Method | Endpoint | Request body | Thành công |
| --- | --- | --- | --- |
| GET | `/health` | Không có | `200`, trạng thái server |
| POST | `/api/auth/register` | `email`, `password` (ít nhất 8 ký tự) | `201`, thông tin user |
| POST | `/api/auth/login` | `email`, `password` | `200`, `accessToken` và `user` |
| POST | `/api/auth/refresh` | `refreshToken` | `200`, `accessToken` và `refreshToken` mới |
| POST | `/api/auth/logout` | `refreshToken` | `204`, không có body |

Ví dụ đăng ký:

```powershell
$body = @{ email = 'demo@example.com'; password = 'example-password' } | ConvertTo-Json
Invoke-RestMethod -Uri http://localhost:8080/api/auth/register -Method Post -ContentType 'application/json' -Body $body
```

Đăng nhập với cùng tài khoản:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/api/auth/login -Method Post -ContentType 'application/json' -Body $body
```

Access token có thời hạn 15 phút; refresh session có thời hạn 30 ngày. Hiện usecase đăng nhập tạo refresh token nhưng HTTP handler chưa trả `refreshToken` trong response, nên client chưa thể hoàn tất luồng refresh/logout từ kết quả login. API bài viết chưa được đăng ký trong HTTP router.

## Chạy Frontend

```powershell
cd frontend
pnpm install --frozen-lockfile
pnpm dev
```

Mở URL được Vite in trong terminal, thường là `http://localhost:5173`, rồi truy cập `/login` hoặc `/blog`.

Các lệnh kiểm tra và build, chạy trong thư mục `frontend`:

```powershell
pnpm lint
pnpm build
pnpm preview
```

`pnpm preview` phục vụ bản build sau khi `pnpm build` thành công. Frontend hiện chưa có script chạy test tự động.

## Test Backend

Chạy bộ unit test auth và token từ thư mục gốc dự án:

```powershell
cd backend
go test -v ./internal/usecase/auth ./internal/infrastructure/token
```

Các test kiểm tra đăng ký, băm mật khẩu, đăng nhập, tạo session, xoay refresh token, thu hồi session, đăng xuất và tạo/xác minh token. Test sử dụng repository giả lập, không cần `.env`, AWS credentials hoặc DynamoDB. Phần post không nằm trong bộ test này.

Xem độ phủ test:

```powershell
go test -cover ./internal/usecase/auth ./internal/infrastructure/token
```

Nếu gặp lỗi `Access is denied` với Go build cache trên Windows:

```powershell
$env:GOCACHE = Join-Path $env:TEMP 'codex-go-cache-blogaws'
go test -v ./internal/usecase/auth ./internal/infrastructure/token
```

Kết quả `PASS` hoặc `ok` cho biết các test đã chạy thành công. Unit test hiện chưa kiểm tra HTTP handler hoặc tích hợp DynamoDB thực tế.
