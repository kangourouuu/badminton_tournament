# MASTERPLAN V2.0: BADMINTON TOURNAMENT MANAGER

## 1. Hệ Thống Tính Năng (Detailed Feature List)

### A. Phân Hệ Input (Google Ecosystem)

- **Form Đăng Ký:**
  - Trường bắt buộc: Họ và Tên, Email (dùng làm ID định danh), Phân loại Pool (Dropdown: "Mesoneer" hoặc "Lab").
- **Real-time Sync (Apps Script):**
  - Cơ chế: Trigger `onFormSubmit`.
  - Logic: Chuyển đổi dữ liệu Form thành JSON Payload -> Gọi POST Webhook tới Backend.
  - Data Integrity: Backend phải xử lý Upsert (Nếu email đã tồn tại -> Cập nhật Pool/Tên; Nếu chưa -> Tạo mới). Tránh việc 1 người điền form 2 lần bị tính là 2 người.

### B. Phân Hệ Quản Trị (Admin Dashboard)

- **Authentication:**
  - Cơ chế: Password-based (Env Var).
  - Session: JWT lưu Cookie/LocalStorage.
- **Pool Management (Mesoneer vs Lab):**
  - Dashboard chia 2 tabs rõ rệt. Dữ liệu của Pool nào nằm riêng Pool đó.
- **Team Generator (The Wheel):**
  - Logic: Lấy danh sách Participant trong Pool -> Shuffle ngẫu nhiên -> Cắt đôi danh sách ghép thành Team.
  - Manual Override: Admin có thể xóa Team nếu ghép sai.
- **Bracket Manager (GSL Automator):**
  - Thao tác: Kéo thả (hoặc chọn) 4 Team vào 1 Group.
  - Automation: Hệ thống tự sinh 5 trận đấu (M1->M5) với trạng thái "Scheduled".
- **Match Operator:**
  - Nhập điểm: Set 1, Set 2, Set 3.
  - Upload Video: Paste link YouTube vào ô input -> Backend lưu ID video.
  - Submit Result: Hệ thống tự tính thắng/thua -> Tự điền tên đội thắng vào trận kế tiếp.

### C. Phân Hệ Hiển Thị (Public View)

- **Live Bracket:**
  - Hiển thị sơ đồ cây GSL (CSS Grid).
  - Auto-refresh: Polling 30s/lần để cập nhật tỉ số mới nhất mà không cần F5.
- **Media Experience:**
  - Nút "Watch Replay" (Icon Play phẳng) hiện ra khi trận đấu có link video.
  - Bấm vào hiện Modal xem ngay trên web.

## 2. Technical Stack & Architecture

- **Backend:** Golang 1.22 (Gin, Bun ORM, Logrus).
- **Database:** PostgreSQL (Neon.tech) - Thiết kế Schema chặt chẽ cho GSL.
- **Frontend:** Vue 3 (Script Setup, TypeScript, Vite).
- **Styling:** TailwindCSS + Font "Outfit" (Strict Tech-Flat Theme).
- **Deployment:**
  - Backend: Render (Docker).
  - Frontend: Vercel.
  - CI/CD: Github Actions.
