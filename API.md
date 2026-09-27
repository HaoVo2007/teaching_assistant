## 1. Sản phẩm (để UI đúng việc)

Giáo viên tiểu học (lớp 1–5) soạn **câu trắc nghiệm / đúng-sai**, gom **bộ** (chỉ để ôn, xóa bộ không xóa câu), tạo **lớp**, gắn **một phụ huynh ↔ một học sinh**, **giao bài** cho cả lớp (chọn từng câu + hạn nộp). Phụ huynh đăng nhập, thấy **đúng một con**, làm bài **hộ**, nộp **một lần**, rồi xem điểm.

Quy tắc bắt buộc trên UI:

- Phụ huynh **chưa nộp**: hiện đề, **cấm** hiện đáp án đúng và lời giải.
- Phụ huynh **đã nộp**: hiện đáp án + điểm của con.
- Giáo viên **luôn** thấy đáp án và điểm cả lớp.
- Không nộp lại. Không nộp sau hạn.
- Gắn/đổi phụ huynh bằng **mã học sinh** + chọn tài khoản phụ huynh.
- Sửa danh sách lớp: giữ học sinh bằng **`id`**, thêm bằng **tên**. Không gửi mã `HS…` khi sửa lớp.

---

## 2. Kết nối

| Mục | Giá trị |
|---|---|
| Base URL | Biến môi trường `VITE_API_BASE_URL`, ví dụ `http://localhost:2000/api/v1` |
| CORS FE dev | `http://localhost:5173` |
| Auth | Header `Authorization: Bearer <token>` |
| Token | JWT, mặc định 24 giờ. Logout **không** thu hồi — chỉ xóa token phía client |
| Role | `admin` \| `teacher` \| `parent` |
| ID | Chuỗi hex 24 ký tự (Mongo ObjectID) |
| Thời gian | `created_at`, `updated_at`, `submitted_at` = RFC3339 (`2026-09-16T08:00:00Z`) |
| `due_date` | Luôn `YYYY-MM-DD` (UTC). **Không** đổi sang ngày giờ Việt Nam |

Thiếu/sai token → `401` `UNAUTHORIZED`. Sai vai trò trên route → `403` `FORBIDDEN`.

---

## 3. Envelope — mọi response

Thành công:

```json
{
  "success": true,
  "message": "Login successful",
  "data": {}
}
```

Lỗi:

```json
{
  "success": false,
  "message": "Invalid email or password",
  "error": { "code": "INVALID_CREDENTIALS" }
}
```

Quy tắc parse:

1. HTTP 2xx và `success === true` → dùng `data`.
2. Còn lại → đọc `error.code` + `message`. Đừng parse HTML.
3. **Create / update / delete** thường `data: null` (trừ login và `POST /admin/users/create`).
4. **Danh sách:** `meta` nằm **trong `data`**, không ở root.

```json
{
  "success": true,
  "message": "Questions fetched successfully",
  "data": {
    "questions": [],
    "meta": { "page": 1, "limit": 10, "total": 40, "total_pages": 4 }
  }
}
```

Query phân trang (mọi list):

| Query | Mặc định | Ghi chú |
|---|---|---|
| `page_index` | `1` | ≥ 1 |
| `page_size` | `10` | Server clamp tối đa 100 |

`meta.page` = `page_index` đã chuẩn hóa. `meta.limit` = `page_size` đã clamp. `meta.total_pages` có thể `0` khi `total = 0`.

HTTP client:

- JSON: `Content-Type: application/json`.
- Câu hỏi + lớp: `multipart/form-data` (để browser set boundary; **đừng** tự set Content-Type).
- `credentials` không bắt buộc (token trên header). CORS `AllowCredentials: true`.
- Timeout hợp lý (15–30s). Ảnh lớp tối đa 5MB.

---

## 4. Auth, lưu trữ, routing

Sau `POST /auth/login` lưu:

```ts
{ token: string, user: { id, username, email, role } }
```

Gợi ý: `localStorage` key `ta_auth`. Axios/fetch interceptor gắn Bearer.

`POST /auth/logout`: gọi rồi xóa storage, về `/login`. Token cũ vẫn dùng được đến hết hạn — không sao.

**Không** làm màn Register trên sản phẩm. `POST /auth/register` luôn tạo `teacher` — bỏ qua.

Route gợi ý (đặt guard theo `role`):

| Path | Role | Màn |
|---|---|---|
| `/login` | public | Đăng nhập |
| `/admin/users` | admin | Tạo giáo viên / phụ huynh |
| `/teacher/questions` | teacher | Ngân hàng câu |
| `/teacher/questions/new` | teacher | Tạo câu |
| `/teacher/questions/:id/edit` | teacher | Sửa câu |
| `/teacher/question-sets` | teacher | Bộ câu |
| `/teacher/question-sets/new` | teacher | Tạo bộ |
| `/teacher/question-sets/:id` | teacher | Chi tiết / sửa bộ |
| `/teacher/classes` | teacher | Danh sách lớp |
| `/teacher/classes/new` | teacher | Tạo lớp |
| `/teacher/classes/:id` | teacher | Chi tiết lớp + roster + gắn PH |
| `/teacher/homeworks` | teacher | Danh sách bài giao |
| `/teacher/homeworks/new` | teacher | Giao bài |
| `/teacher/homeworks/:id` | teacher | Chi tiết / sửa bài |
| `/teacher/homeworks/:id/submissions` | teacher | Điểm cả lớp |
| `/parent` | parent | Con + danh sách bài |
| `/parent/homeworks/:id` | parent | Làm bài (ẩn đáp án nếu chưa nộp) |
| `/parent/homeworks/:id/result` | parent | Điểm + đáp án **sau khi nộp** |

Sai role → redirect về đúng home. Chưa login → `/login`. Login rồi vào `/login` → home theo role.

---

## 5. Enum và nhãn UI (bắt buộc dùng đúng value API)

### Role

| Value | Hiện trên UI |
|---|---|
| `admin` | Quản trị |
| `teacher` | Giáo viên |
| `parent` | Phụ huynh |

Tạo user **chỉ** `teacher` \| `parent`. Không hiện option `admin` / `student`.

### Student `status`

| Value | UI |
|---|---|
| `active` | Đang học |
| `inactive` | Đã nghỉ |

GET lớp chỉ trả học sinh còn trên roster (đang học). Em nghỉ không còn trong `students[]`.

### Question / set `type`

| Value | UI |
|---|---|
| `multiple_choice` | Trắc nghiệm |
| `true_false` | Đúng / sai |

Cấm `matching`. Tạo/sửa/filter loại khác → `400` `INVALID_QUESTION_TYPE`.

### `subject`

| Value | UI |
|---|---|
| `vietnamese` | Tiếng Việt |
| `mathematics` | Toán |
| `ethics` | Đạo đức |
| `english` | Tiếng Anh |
| `nature_and_society` | Tự nhiên và xã hội |
| `history_and_geography` | Lịch sử và địa lý |
| `science` | Khoa học |
| `informatics` | Tin học |
| `technology` | Công nghệ |
| `physical_education` | Thể dục |
| `music` | Âm nhạc |
| `art` | Mỹ thuật |
| `experiential_activities` | Hoạt động trải nghiệm |

### `grade`

Value API là **string**: `"1"` `"2"` `"3"` `"4"` `"5"`. UI: Lớp 1 … Lớp 5.

### `difficulty`

| Value | UI | Trọng số điểm |
|---|---|---|
| `easy` | Dễ | 1 |
| `medium` | Vừa | 2 |
| `hard` | Khó | 3 |

Điểm bài: trọng số / tổng trọng số × 100, làm tròn 2 chữ số, tổng `max_score` = 100. Server **chấm lúc GET bài nộp**. POST nộp **không** gửi điểm.

### Hạn nộp

- Gửi / nhận `due_date`: `YYYY-MM-DD`.
- Server lưu 00:00:00 UTC ngày đó.
- Nộp bị từ chối khi `now UTC > due_date 00:00 UTC` → gần như hết hạn **ngay đầu ngày UTC** đó (`400` `DUE_DATE_EXPIRED`).
- UI ghi chú: “Hạn theo ngày UTC, không đổi giờ Việt Nam.”
- Disable nút Nộp nếu ngày UTC hôm nay **sau** `due_date` (so sánh chuỗi ngày UTC `YYYY-MM-DD` với `due_date`; nếu `today > due_date` thì hết hạn. Nếu `today === due_date` vẫn có thể hết hạn vì mốc 00:00 UTC — an toàn hơn: coi `today >= due_date` là hết hạn trên UI, hoặc vẫn cho bấm và hiện lỗi server).

### Mã học sinh

Format: `HS` + `YYMMDD` + `-` + 4 ký tự, ví dụ `HS260916-M7KQ`. **Chỉ dùng cho claim.** Sửa lớp dùng `id`.

---

## 6. TypeScript types (khớp JSON)

Field có `omitempty` trên backend: khi **ẩn đáp án** hoặc rỗng, field **không có trong JSON** (không phải `null`). FE dùng optional `?`.

```ts
type Role = "admin" | "teacher" | "parent";
type StudentStatus = "active" | "inactive";
type QuestionType = "multiple_choice" | "true_false";
type Subject =
  | "vietnamese" | "mathematics" | "ethics" | "english"
  | "nature_and_society" | "history_and_geography" | "science"
  | "informatics" | "technology" | "physical_education"
  | "music" | "art" | "experiential_activities";
type Grade = "1" | "2" | "3" | "4" | "5";
type Difficulty = "easy" | "medium" | "hard";

interface Envelope<T> {
  success: boolean;
  message: string;
  data?: T | null;
  error?: { code: string; details?: unknown };
}

interface Meta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

interface User {
  id: string;
  username: string;
  email: string;
  role: Role;
}

interface AuthData {
  token: string;
  user: User;
}

interface Student {
  id: string;
  name: string;
  code: string;
  image: string;
  public_id: string;
  class_id: string;
  status: StudentStatus;
  guardian: User | null;
  created_at: string;
  updated_at: string;
}

interface Class {
  id: string;
  name: string;
  description: string;
  image: string;
  public_id: string;
  students: Student[];
  created_by: string;
  created_at: string;
  updated_at: string;
}

interface Question {
  id: string;
  type: QuestionType;
  subject: Subject;
  grade: Grade;
  difficulty: Difficulty;
  question: string;
  options?: string[];
  correct_index?: number;   // vắng mặt = đang ẩn đáp án (parent chưa nộp)
  correct_bool?: boolean;   // vắng mặt = đang ẩn
  explanation?: string;     // vắng mặt hoặc "" = đang ẩn
  created_by: string;
  created_at: string;
  updated_at: string;
}

interface QuestionSet {
  id: string;
  title: string;
  question_type: QuestionType;
  description: string | null;
  questions: Question[];
  created_at: string;
  updated_at: string;
}

interface Homework {
  id: string;
  class_id: string;
  title: string;
  description: string | null;
  due_date: string; // YYYY-MM-DD
  questions: Question[];
  created_at: string;
  updated_at: string;
}

interface StudentAnswerView {
  question: Question;
  selected_index?: number | null;
  selected_bool?: boolean | null;
  is_correct: boolean;
  score: number;
  max_score: number;
}

interface HomeworkSubmission {
  id: string;
  homework: Homework;
  student_id: string;
  is_submitted: boolean;
  student_answers: StudentAnswerView[];
  total_score: number;
  max_score: number; // luôn 100
  submitted_by: string;
  submitted_at: string;
  created_at: string;
  updated_at: string;
}

interface ParentsPage { parents: User[]; meta: Meta }
interface ClassesPage { classes: Class[]; meta: Meta }
interface QuestionsPage { questions: Question[]; meta: Meta }
interface QuestionSetsPage { question_sets: QuestionSet[]; meta: Meta }
interface HomeworksPage { homeworks: Homework[]; meta: Meta }
interface SubmissionsPage { homework_submissions: HomeworkSubmission[]; meta: Meta }
```

Helper parent “đã nộp bài này”:

```ts
function homeworkAnswersVisible(hw: Homework): boolean {
  return hw.questions.some(
    (q) => q.correct_index !== undefined || q.correct_bool !== undefined || !!q.explanation
  );
}
```

---

## 7. Catalog endpoint (đủ set, không thêm)

Mọi path dưới đây **sau** base `/api/v1`.

| Method | Path | Role | Content-Type | HTTP OK | `data` |
|---|---|---|---|---|---|
| POST | `/auth/login` | public | JSON | 200 | `AuthData` |
| POST | `/auth/register` | public | JSON | 200 | `AuthData` — **không dùng** |
| POST | `/auth/logout` | đã login | — | 200 | `null` |
| POST | `/admin/users/create` | admin | JSON | 200 | `User` |
| GET | `/users/parents` | teacher | — | 200 | `ParentsPage` |
| POST | `/students/claim` | teacher | JSON | 200 | `null` |
| GET | `/students/by-guardian` | parent | — | 200 | **một** `Student` (không phải mảng) |
| POST | `/questions` | teacher | multipart | **201** | `null` |
| GET | `/questions` | teacher | — | 200 | `QuestionsPage` |
| GET | `/questions/:id` | teacher | — | 200 | `Question` |
| PUT | `/questions/:id` | teacher | multipart | 200 | `null` |
| DELETE | `/questions/:id` | teacher | — | 200 | `null` |
| POST | `/question-sets` | teacher | JSON | 200 | `null` |
| GET | `/question-sets` | teacher | — | 200 | `QuestionSetsPage` |
| GET | `/question-sets/:id` | teacher | — | 200 | `QuestionSet` |
| PUT | `/question-sets/:id` | teacher | JSON | 200 | `null` |
| DELETE | `/question-sets/:id` | teacher | — | 200 | `null` |
| POST | `/classes` | teacher | multipart | **201** | `null` |
| GET | `/classes` | teacher | — | 200 | `ClassesPage` |
| GET | `/classes/:id` | teacher | — | 200 | `Class` |
| PUT | `/classes/:id` | teacher | multipart | 200 | `null` |
| DELETE | `/classes/:id` | teacher | — | 200 | `null` |
| POST | `/homeworks` | teacher | JSON | 200 | `null` |
| GET | `/homeworks` | teacher | — | 200 | `HomeworksPage` |
| GET | `/homeworks/class/:class_id` | teacher | — | 200 | `HomeworksPage` |
| GET | `/homeworks/student/by-guardian` | parent | — | 200 | `HomeworksPage` |
| GET | `/homeworks/:id` | teacher | — | 200 | `Homework` |
| PUT | `/homeworks/:id` | teacher | JSON | 200 | `null` |
| DELETE | `/homeworks/:id` | teacher | — | 200 | `null` |
| POST | `/homework-submissions` | parent | JSON | 200 | `null` |
| GET | `/homework-submissions` | teacher | — | 200 | `SubmissionsPage` |
| GET | `/homework-submissions/homework/:homework_id` | teacher | — | 200 | `SubmissionsPage` |
| GET | `/homework-submissions/student/by-guardian/:homework_id` | parent | — | 200 | **một** `HomeworkSubmission` |
| GET | `/homework-submissions/:id` | teacher | — | 200 | `HomeworkSubmission` |

Thứ tự path quan trọng: `.../student/by-guardian` và `.../class/:id` **trước** `/:id`. FE **không** được gọi `GET /homeworks/:id` hay `GET /homework-submissions/:id` bằng token parent (403).

Không có: list teacher, list admin, GET user by id, update user, upload riêng, nộp lại, xóa bài nộp.

---

## 8. Chi tiết từng API

### 8.1 `POST /auth/login` — public

```json
{ "email": "gv.lan@demo.local", "password": "Demo@123" }
```

`200` `data`:

```json
{
  "token": "eyJ...",
  "user": {
    "id": "66f0aaaaaaaaaaaaaaaaaaaa",
    "username": "nguyen_thi_lan",
    "email": "gv.lan@demo.local",
    "role": "teacher"
  }
}
```

| HTTP | code |
|---|---|
| 400 | `INVALID_REQUEST_BODY` `INVALID_EMAIL` `INVALID_PASSWORD` |
| 401 | `INVALID_CREDENTIALS` |

### 8.2 `POST /auth/logout` — Bearer bất kỳ role

Không body. `200` `data: null`. FE xóa token dù API lỗi.

### 8.3 `POST /admin/users/create` — admin

```json
{
  "username": "phu_huynh_moi",
  "email": "ph.moi@demo.local",
  "password": "Demo@123",
  "role": "parent"
}
```

`role` chỉ `teacher` \| `parent`. `200` `data` = `User` (không token). Người mới tự login.

| HTTP | code |
|---|---|
| 400 | `INVALID_USERNAME` `INVALID_EMAIL` `INVALID_PASSWORD` `INVALID_ROLE` |
| 409 | `EMAIL_ALREADY_EXISTS` `USER_ALREADY_EXISTS` |

Form: username, email, mật khẩu, chọn vai trò Giáo viên / Phụ huynh. Không tạo admin.

### 8.4 `GET /users/parents` — teacher

Query: `q` (regex không phân biệt hoa thường trên `username` **hoặc** `email`), `page_index`, `page_size`.

`200` `data`:

```json
{
  "parents": [
    { "id": "66f0...", "username": "phu_huynh_01", "email": "ph.01@demo.local", "role": "parent" }
  ],
  "meta": { "page": 1, "limit": 10, "total": 20, "total_pages": 2 }
}
```

`id` = `parent_id` khi claim. Parent đã có con **vẫn hiện**. Gắn trùng → claim `409`.

UI: combobox/search `q`, phân trang, hiện email + tên.

### 8.5 `POST /students/claim` — teacher

```json
{ "code": "HS260916-M7KQ", "parent_id": "66f0bbbbbbbbbbbbbbbbbbbb" }
```

- Học sinh `active`, thuộc lớp của giáo viên đang login.
- `parent_id` phải là user `parent`.
- Chưa gắn → tạo. Đã gắn người khác → **đổi** nếu parent mới chưa có con.
- Cùng parent đã gắn → `200` no-op.
- Parent mới đã có con khác → `409` `PARENT_ALREADY_HAS_STUDENT`.

| HTTP | code | UI |
|---|---|---|
| 400 | `INVALID_STUDENT_CODE` | Thiếu mã |
| 400 | `INVALID_PARENT_ID` | Thiếu / sai id |
| 400 | `NOT_A_PARENT` | User không phải phụ huynh |
| 403 | `STUDENT_INACTIVE` | Em đã nghỉ |
| 403 | `STUDENT_NOT_IN_CLASS` | Không phải lớp của bạn |
| 404 | `STUDENT_NOT_FOUND` | Sai mã |
| 404 | `PARENT_NOT_FOUND` | Không có user |
| 409 | `PARENT_ALREADY_HAS_STUDENT` | Phụ huynh đã có con |

UI trên trang lớp: nhập/chọn mã em + chọn parent từ `GET /users/parents` → Gắn / Đổi.

### 8.6 `GET /students/by-guardian` — parent

Không query. `200` **một** `Student`. `guardian` thường `null`.

`404` `GUARDIAN_NOT_FOUND` → “Tài khoản chưa được gắn học sinh. Nhờ giáo viên gắn mã của con.” Không hiện form nộp.

---

### 8.7 Lớp — teacher — `multipart/form-data`

#### `POST /classes` → `201` `data: null`

| Field | Bắt buộc | |
|---|---|---|
| `name` | có | tên lớp |
| `description` | có | mô tả |
| `image` | không | file ảnh, ≤ 5MB |
| `students` | không | **lặp cùng key**: mỗi value = **tên** học sinh |

```ts
const fd = new FormData();
fd.append("name", "Lớp 3A");
fd.append("description", "Lớp buổi sáng");
fd.append("students", "Nguyễn Văn An");
fd.append("students", "Trần Thị Bình");
if (file) fd.append("image", file);
```

Sau đó `GET /classes` hoặc `GET /classes/:id` để lấy `id` / `code`.

`400` `INVALID_CLASS` (thiếu tên), `INVALID_IMAGE` (ảnh > 5MB).

#### `GET /classes`

Query: `name` (regex), `page_index`, `page_size`.  
`data.classes[]` + `data.meta`. Mỗi em có `guardian` nếu đã gắn.

#### `GET /classes/:id`

Một `Class`. Không phải lớp mình → `403` `CLASS_FORBIDDEN`. Sai id → `404` `CLASS_NOT_FOUND`.

Ví dụ `data`:

```json
{
  "id": "66f0cccccccccccccccccccc",
  "name": "Lớp 3A",
  "description": "Lớp buổi sáng",
  "image": "https://res.cloudinary.com/...",
  "public_id": "classes/abc",
  "created_by": "66f0aaaaaaaaaaaaaaaaaaaa",
  "students": [
    {
      "id": "66f0dddddddddddddddddddd",
      "name": "Nguyễn Văn An",
      "code": "HS260916-M7KQ",
      "image": "",
      "public_id": "",
      "class_id": "66f0cccccccccccccccccccc",
      "status": "active",
      "guardian": {
        "id": "66f0bbbbbbbbbbbbbbbbbbbb",
        "username": "phu_huynh_01",
        "email": "ph.01@demo.local",
        "role": "parent"
      },
      "created_at": "2026-09-16T08:00:00Z",
      "updated_at": "2026-09-16T08:00:00Z"
    }
  ],
  "created_at": "2026-09-16T08:00:00Z",
  "updated_at": "2026-09-16T08:00:00Z"
}
```

`guardian: null` = chưa gắn. Hiện mã lớn, nút copy.

#### `PUT /classes/:id`

Field optional: `name`, `description`, `image`, `students`.

**Không gửi `students`** → roster **không đổi**.

**Có gửi `students`** → đây là **danh sách mới đầy đủ**:

| Value | Server hiểu |
|---|---|
| Đúng 24 hex (`id` từ GET) | **Giữ** em đó |
| Không phải ObjectID | **Tên mới** → tạo em mới + mã mới |
| `id` cũ không gửi | Em **nghỉ** (ẩn), không xóa |

**Cấm** gửi `code` (`HS26…`) — bị hiểu là tên mới.

```ts
// Giữ An, Bình; thêm C; không gửi id em nghỉ
fd.append("students", an.id);
fd.append("students", binh.id);
fd.append("students", "Lê Văn C");
```

Khi user bấm Lưu roster, **luôn** append mọi `id` còn giữ. Quên = nghỉ hết lớp.

Lỗi: `STUDENT_NOT_FOUND`, `STUDENT_NOT_IN_CLASS`, `INVALID_CLASS`, `INVALID_IMAGE`.

#### `DELETE /classes/:id`

Hiện **không xóa** nếu lớp còn bất kỳ bài giao → `409` `CLASS_IN_USE` (“Không xóa được vì lớp còn bài tập.”).

---

### 8.8 Câu hỏi — teacher — `multipart/form-data`

Chỉ câu **của giáo viên đang login**.

#### `POST /questions` → `201` | `PUT /questions/:id` → `200` | `data: null`

| Field | Trắc nghiệm | Đúng/sai |
|---|---|---|
| `type` | `multiple_choice` | `true_false` |
| `subject` | enum | enum |
| `grade` | `"1"`…`"5"` | `"1"`…`"5"` |
| `difficulty` | `easy`/`medium`/`hard` | |
| `question` | nội dung | nội dung |
| `options` | lặp key, ≥ 2 đáp án | **không gửi** |
| `correct_index` | 0-based (0 = đáp án đầu) | **không gửi** |
| `correct_bool` | **không gửi** | `true` hoặc `false` (string form) |
| `explanation` | tuỳ | tuỳ |

```ts
// MC
fd.append("type", "multiple_choice");
fd.append("subject", "mathematics");
fd.append("grade", "3");
fd.append("difficulty", "easy");
fd.append("question", "15 + 27 bằng bao nhiêu?");
fd.append("options", "32");
fd.append("options", "42");
fd.append("options", "41");
fd.append("options", "52");
fd.append("correct_index", "1"); // "42"
fd.append("explanation", "15+27=42");

// T/F
fd.append("type", "true_false");
fd.append("subject", "mathematics");
fd.append("grade", "3");
fd.append("difficulty", "easy");
fd.append("question", "Số 0 là số chẵn.");
fd.append("correct_bool", "true");
fd.append("explanation", "0 chia hết cho 2");
```

`type` lạ → `400` `INVALID_QUESTION_TYPE`.

Không xóa / không đổi loại–độ khó–đáp án–options nếu câu đã nằm trong bộ, bài giao, hoặc bài nộp → `409` `QUESTION_IN_USE`. UI: disable sửa đáp án khi lỗi này, hoặc báo “Câu đang được dùng.”

#### `GET /questions`

Query: `question_type`, `question_name` (regex nội dung), `subject`, `grade`, `difficulty`, `page_index`, `page_size`.

`data.questions[]` + `data.meta`. Teacher **luôn** thấy đáp án.

#### `GET /questions/:id`

Một `Question` (có đáp án).

```json
{
  "id": "66f0eeeeeeeeeeeeeeeeeeee",
  "type": "multiple_choice",
  "subject": "mathematics",
  "grade": "3",
  "difficulty": "easy",
  "question": "15 + 27 bằng bao nhiêu?",
  "options": ["32", "42", "41", "52"],
  "correct_index": 1,
  "explanation": "15+27=42",
  "created_by": "66f0aaaaaaaaaaaaaaaaaaaa",
  "created_at": "2026-09-16T08:00:00Z",
  "updated_at": "2026-09-16T08:00:00Z"
}
```

T/F: không `options` / `correct_index`; có `correct_bool`.

#### `DELETE /questions/:id`

`409` `QUESTION_IN_USE` nếu đang được tham chiếu. `403` `QUESTION_FORBIDDEN` nếu không phải chủ.

---

### 8.9 Bộ câu hỏi — teacher — JSON

Mọi câu trong bộ **cùng `question_type`**. Xóa bộ **không** xóa câu. Bộ **không** dùng để giao bài (giao bài chọn từng `question.id`).

#### `POST /question-sets` → `200` `data: null`

```json
{
  "title": "Ôn toán lớp 3",
  "question_type": "multiple_choice",
  "description": "Tuần 1",
  "questions": ["66f0eeeeeeeeeeeeeeeeeeee", "66f0ffffffffffffffffffffff"]
}
```

`questions` = id câu, ít nhất 1. Câu khác loại → `400` `QUESTION_TYPE_MISMATCH`. `question_type` lạ → `400` `INVALID_QUESTION_TYPE`. Thiếu title → `400` `INVALID_TITLE`.

UI tạo bộ: filter `GET /questions?question_type=...` rồi multi-select.

#### `GET /question-sets`

Query: `title`, `question_type`, `page_index`, `page_size`.  
`data.question_sets[]` — mỗi bộ đã embed `questions[]`.

#### `GET /question-sets/:id`

Một `QuestionSet`.

#### `PUT /question-sets/:id`

```json
{
  "title": "Ôn toán lớp 3 (sửa)",
  "description": "Tuần 2",
  "questions": ["66f0eeeeeeeeeeeeeeeeeeee"]
}
```

**Không** đổi `question_type`. `questions` là **list thay thế** (không merge). Field pointer: bỏ field = không đổi.

#### `DELETE /question-sets/:id`

Xóa bộ, câu vẫn còn trong ngân hàng.

---

### 8.10 Bài tập — JSON

#### `POST /homeworks` — teacher → `200` `data: null`

```json
{
  "class_id": "66f0cccccccccccccccccccc",
  "title": "Toán cuối tuần",
  "description": "Làm trong 30 phút",
  "questions": ["66f0eeeeeeeeeeeeeeeeeeee", "66f0ffffffffffffffffffffff"],
  "due_date": "2026-10-05"
}
```

- Chỉ lớp của mình.
- `questions` = id từng câu, ≥ 1. **Trộn MC + T/F được.** Không gửi id bộ.
- `due_date` bắt buộc `YYYY-MM-DD`.

| HTTP | code |
|---|---|
| 400 | `INVALID_TITLE` `INVALID_CLASS_ID` `INVALID_QUESTIONS` `INVALID_DUE_DATE` |
| 403 | `HOMEWORK_FORBIDDEN` |

UI: chọn lớp (`GET /classes`), chọn câu (`GET /questions`, checkbox), date picker → format `YYYY-MM-DD`. Có thể mở bộ để **xem** rồi tick từng câu — không có API “giao cả bộ”.

#### `GET /homeworks` — teacher

`data.homeworks[]` + `data.meta`. Mỗi bài **có đáp án**.

#### `GET /homeworks/class/:class_id` — teacher

Cùng shape, lọc theo lớp (tab trong trang lớp).

#### `GET /homeworks/:id` — teacher

Một `Homework` (có đáp án). Không phải chủ → `403` `HOMEWORK_FORBIDDEN`.

#### `PUT /homeworks/:id` — teacher

Cùng field create. Field rỗng thường = không đổi (title/due/class chỉ đổi khi gửi khác rỗng). Đổi `questions` khi **đã có bài nộp** → `409` `HOMEWORK_IN_USE`. UI: nếu đã có nộp, khóa danh sách câu + nút xóa.

#### `DELETE /homeworks/:id` — teacher

Đã có nộp → `409` `HOMEWORK_IN_USE`.

#### `GET /homeworks/student/by-guardian` — parent

Không gửi `student_id`. Bài của **lớp con**.

- **Chưa nộp:** `correct_index`, `correct_bool`, `explanation` **vắng mặt**.
- **Đã nộp:** hiện đủ đáp án + lời giải.

`data.homeworks[]` + `data.meta`.

Ví dụ một bài chưa nộp:

```json
{
  "id": "66f0abcabcabcabcabcabcab",
  "class_id": "66f0cccccccccccccccccccc",
  "title": "Toán cuối tuần",
  "description": "Làm trong 30 phút",
  "due_date": "2026-10-05",
  "questions": [
    {
      "id": "66f0eeeeeeeeeeeeeeeeeeee",
      "type": "multiple_choice",
      "subject": "mathematics",
      "grade": "3",
      "difficulty": "easy",
      "question": "15 + 27 bằng bao nhiêu?",
      "options": ["32", "42", "41", "52"],
      "created_by": "66f0aaaaaaaaaaaaaaaaaaaa",
      "created_at": "2026-09-16T08:00:00Z",
      "updated_at": "2026-09-16T08:00:00Z"
    }
  ],
  "created_at": "2026-09-16T08:00:00Z",
  "updated_at": "2026-09-16T08:00:00Z"
}
```

UI list: badge “Chưa nộp” / “Đã nộp” theo `homeworkAnswersVisible`. Nút “Làm bài” vs “Xem điểm”.

Parent **cấm** dùng `GET /homeworks/:id` (403 và lộ đáp án).

---

### 8.11 Nộp bài và điểm

#### `POST /homework-submissions` — parent → `200` `data: null`

Không gửi `student_id`. Server lấy con từ guardian.

Phải **đủ mọi câu** của bài, mỗi `question_id` đúng một lần.

```json
{
  "homework_id": "66f0abcabcabcabcabcabcab",
  "student_answers": [
    { "question_id": "66f0eeeeeeeeeeeeeeeeeeee", "selected_index": 1 },
    { "question_id": "66f0ffffffffffffffffffffff", "selected_bool": true }
  ]
}
```

| Loại câu | Bắt buộc | Cấm thiếu |
|---|---|---|
| `multiple_choice` | `selected_index` (number, 0-based) | `selected_bool` không dùng |
| `true_false` | `selected_bool` (boolean) | `selected_index` không dùng |

Không gửi `is_correct`, `score` (server bỏ qua nếu có).

| HTTP | code | UI |
|---|---|---|
| 400 | `INVALID_SUBMISSION` | Thiếu homework_id / answers rỗng |
| 400 | `QUESTION_MISMATCH` | Thiếu/thừa/trùng câu |
| 400 | `INVALID_STUDENT_ANSWER` | Sai field theo loại |
| 400 | `DUE_DATE_EXPIRED` | Hết hạn |
| 403 | `STUDENT_NOT_IN_CLASS` | Con không thuộc lớp bài |
| 404 | `GUARDIAN_NOT_FOUND` | Chưa gắn con |
| 404 | `HOMEWORK_NOT_FOUND` | Sai bài |
| 409 | `SUBMISSION_ALREADY_EXISTS` | Đã nộp — chuyển trang điểm |

UI làm bài: radio options / Đúng-Sai; bắt buộc trả lời hết mới enable Nộp; confirm một lần. Thành công → `/parent/homeworks/:id/result`.

#### `GET /homework-submissions` — teacher

Mọi bài nộp giáo viên đó. `data.homework_submissions[]` + `data.meta`. Điểm đã chấm.

#### `GET /homework-submissions/homework/:homework_id` — teacher

Cùng shape, **màn điểm cả lớp** (dùng cái này, không dùng list global trừ dashboard).

Bảng: tên em (map `student_id` → `GET /classes/:class_id` `students[]`), `total_score` / 100, `submitted_at`. Click xem chi tiết `GET /homework-submissions/:id`.

#### `GET /homework-submissions/:id` — teacher

Một submission. Không phải của giáo viên → `404` `SUBMISSION_NOT_FOUND`. **Cấm** parent gọi path này.

#### `GET /homework-submissions/student/by-guardian/:homework_id` — parent

`homework_id` trên path (id **bài giao**, không phải id bài nộp).

`200` **một** `HomeworkSubmission` (không wrap list). Có đáp án + điểm.

**Chỉ gọi sau khi đã nộp** (`homeworkAnswersVisible` hoặc POST vừa `200`). Gọi khi chưa nộp có thể **500**.

Nếu `submitted_by` khác parent đang login và khác `""` → `404` `SUBMISSION_NOT_FOUND`.

```json
{
  "id": "66f0ffffffffffffffffffffff",
  "homework": {
    "id": "66f0abcabcabcabcabcabcab",
    "class_id": "66f0cccccccccccccccccccc",
    "title": "Toán cuối tuần",
    "description": "Làm trong 30 phút",
    "due_date": "2026-10-05",
    "questions": [],
    "created_at": "2026-09-16T08:00:00Z",
    "updated_at": "2026-09-16T08:00:00Z"
  },
  "student_id": "66f0dddddddddddddddddddd",
  "is_submitted": true,
  "student_answers": [
    {
      "question": {
        "id": "66f0eeeeeeeeeeeeeeeeeeee",
        "type": "multiple_choice",
        "subject": "mathematics",
        "grade": "3",
        "difficulty": "easy",
        "question": "15 + 27 bằng bao nhiêu?",
        "options": ["32", "42", "41", "52"],
        "correct_index": 1,
        "explanation": "15+27=42",
        "created_by": "66f0aaaaaaaaaaaaaaaaaaaa",
        "created_at": "2026-09-16T08:00:00Z",
        "updated_at": "2026-09-16T08:00:00Z"
      },
      "selected_index": 1,
      "selected_bool": null,
      "is_correct": true,
      "score": 12.5,
      "max_score": 12.5
    }
  ],
  "total_score": 85.5,
  "max_score": 100,
  "submitted_by": "66f0bbbbbbbbbbbbbbbbbbbb",
  "submitted_at": "2026-09-18T10:00:00Z",
  "created_at": "2026-09-18T10:00:00Z",
  "updated_at": "2026-09-18T10:00:00Z"
}
```

UI: điểm lớn `total_score / max_score`; từng câu đúng/sai, đáp án chọn vs đáp án đúng, lời giải.

---

## 9. Màn hình — AI phải dựng đủ

Layout: header tên + role + Đăng xuất. Sidebar theo role. Trống/lỗi/loading trên mọi list.

### 9.1 Login `/login`

Email + mật khẩu. Submit → `POST /auth/login` → lưu auth → redirect:

- `admin` → `/admin/users`
- `teacher` → `/teacher/classes` (hoặc `/teacher/questions`)
- `parent` → `/parent`

Lỗi credentials: “Email hoặc mật khẩu không đúng.” Không có link đăng ký.

### 9.2 Admin `/admin/users`

Form tạo user (username, email, password, role teacher/parent) → `POST /admin/users/create`. Toast thành công + hiện `id` (để biết). Không có API list user cho admin — không bịa bảng toàn hệ thống. Có thể ghi chú: giáo viên lấy phụ huynh qua `GET /users/parents`.

### 9.3 Teacher — ngân hàng câu

List: filter loại / môn / khối / độ khó / tên + phân trang. Card hiện loại, môn, khối, độ khó, nội dung, đáp án. Nút Tạo / Sửa / Xóa.

Tạo/sửa: form multipart. MC: editor options (thêm/xóa dòng), radio `correct_index`. T/F: radio Đúng/Sai. Không UI matching.

### 9.4 Teacher — bộ câu

List + tạo: chọn type rồi chọn câu cùng type. Chi tiết hiện câu. Sửa title/mô tả/list câu. Xóa bộ không xóa câu.

### 9.5 Teacher — lớp

List (search `name`) + tạo (tên, mô tả, ảnh, textarea/chips tên học sinh).

Chi tiết:

- Ảnh, mô tả, bảng học sinh: tên, **mã** (copy), phụ huynh (email hoặc “Chưa gắn”).
- Form gắn/đổi: mã em + search parent (`GET /users/parents?q=`) → `POST /students/claim`.
- Sửa roster: list hiện tại (mỗi em một `id` ẩn) + ô thêm tên; bỏ tick = không gửi `id` đó.
- Tab bài của lớp: `GET /homeworks/class/:class_id`.
- Xóa lớp: confirm; nếu `CLASS_IN_USE` giải thích còn bài giao.

### 9.6 Teacher — giao bài / điểm

Tạo: lớp, title, mô tả, date, multi-select câu (hiện type để trộn). Không nút “chọn cả bộ” như một action API — chỉ helper UI tick câu trong bộ.

Chi tiết bài: đề + đáp án. Sửa/xóa nếu chưa có nộp.

Điểm lớp: `GET /homework-submissions/homework/:homework_id`. Bảng điểm; chi tiết từng bài nộp.

### 9.7 Parent home `/parent`

1. `GET /students/by-guardian` — card con (tên, mã, lớp id).
2. `GET /homeworks/student/by-guardian` — list bài: title, hạn, badge nộp.

`GUARDIAN_NOT_FOUND`: empty state, không crash.

### 9.8 Parent làm bài `/parent/homeworks/:id`

1. Load list (hoặc nhớ từ list) — **không** `GET /homeworks/:id`.
2. Tìm homework theo id. Nếu `homeworkAnswersVisible` → redirect result.
3. Render câu: MC radio `options`; T/F Đúng/Sai. **Không** hiện `correct_*` / `explanation` (chúng vốn vắng).
4. Nộp đủ câu → `POST /homework-submissions` → result.
5. `SUBMISSION_ALREADY_EXISTS` / hết hạn: thông báo, về list hoặc result.

### 9.9 Parent điểm `/parent/homeworks/:id/result`

Nếu chưa nộp → về làm bài.  
`GET /homework-submissions/student/by-guardian/:homework_id` (id là homework id). Hiện điểm + từng câu.

---

## 10. Luồng API theo việc (copy đúng thứ tự)

**Admin tạo tài khoản**

`POST /auth/login` → `POST /admin/users/create`

**Giáo viên soạn → lớp → gắn → giao → điểm**

1. `POST /auth/login`
2. `POST /questions` (multipart) → `GET /questions`
3. (Tuỳ) `POST /question-sets`
4. `POST /classes` → `GET /classes/:id` (lấy `students[].id`, `code`)
5. `GET /users/parents` → `POST /students/claim` `{ code, parent_id }`
6. `POST /homeworks`
7. `GET /homework-submissions/homework/:homework_id`

**Phụ huynh xem đề → nộp → điểm**

1. `POST /auth/login`
2. `GET /students/by-guardian`
3. `GET /homeworks/student/by-guardian` (ẩn đáp án nếu chưa nộp)
4. `POST /homework-submissions` (một lần, đủ câu)
5. `GET /homework-submissions/student/by-guardian/:homework_id` **chỉ sau khi nộp**

---

## 11. Bảng `error.code`

| HTTP | code | Nghĩa / gợi ý UI |
|---|---|---|
| 400 | `INVALID_REQUEST_BODY` `BAD_REQUEST` | Dữ liệu gửi sai |
| 400 | `INVALID_USERNAME` | Tên đăng nhập bắt buộc |
| 400 | `INVALID_EMAIL` | Email bắt buộc / không hợp lệ |
| 400 | `INVALID_PASSWORD` | Mật khẩu bắt buộc |
| 400 | `PASSWORD_MISMATCH` | (ít dùng) |
| 400 | `INVALID_ROLE` | Chỉ giáo viên hoặc phụ huynh |
| 400 | `INVALID_STUDENT_CODE` | Thiếu mã học sinh |
| 400 | `INVALID_PARENT_ID` | Thiếu phụ huynh |
| 400 | `NOT_A_PARENT` | Tài khoản không phải phụ huynh |
| 400 | `INVALID_CLASS` | Tên lớp bắt buộc |
| 400 | `INVALID_IMAGE` | Ảnh quá 5MB |
| 400 | `INVALID_QUESTION_TYPE` | Chỉ trắc nghiệm / đúng-sai |
| 400 | `INVALID_QUESTION` | Nội dung câu không hợp lệ |
| 400 | `INVALID_OPTIONS` | Đáp án MC không hợp lệ |
| 400 | `INVALID_CORRECT_ANSWER` | Đáp án đúng không hợp lệ |
| 400 | `INVALID_SUBJECT` | Môn không hợp lệ |
| 400 | `INVALID_GRADE` | Khối không hợp lệ |
| 400 | `INVALID_TITLE` | Thiếu tiêu đề |
| 400 | `INVALID_QUESTIONS` | List câu không hợp lệ / bài chưa có câu |
| 400 | `QUESTION_TYPE_MISMATCH` | Câu không cùng loại với bộ |
| 400 | `INVALID_HOMEWORK` | Bài tập không hợp lệ |
| 400 | `INVALID_CLASS_ID` | Lớp không hợp lệ |
| 400 | `INVALID_DUE_DATE` | Hạn phải `YYYY-MM-DD` |
| 400 | `INVALID_SUBMISSION` | Bài nộp không hợp lệ |
| 400 | `QUESTION_MISMATCH` | Phải trả lời đúng mọi câu, không trùng |
| 400 | `INVALID_STUDENT_ANSWER` | MC cần `selected_index`, T/F cần `selected_bool` |
| 400 | `DUE_DATE_EXPIRED` | Đã hết hạn nộp |
| 401 | `UNAUTHORIZED` | Hết hạn phiên — về login |
| 401 | `INVALID_CREDENTIALS` | Sai email/mật khẩu |
| 403 | `FORBIDDEN` | Sai vai trò |
| 403 | `CLASS_FORBIDDEN` | Không quản lý lớp này |
| 403 | `QUESTION_FORBIDDEN` | Không quản lý câu này |
| 403 | `QUESTION_SET_FORBIDDEN` | Không quản lý bộ này |
| 403 | `HOMEWORK_FORBIDDEN` | Không quản lý bài này |
| 403 | `STUDENT_NOT_IN_CLASS` | Học sinh không thuộc lớp / không phải lớp bạn |
| 403 | `STUDENT_INACTIVE` | Học sinh đã nghỉ |
| 404 | `NOT_FOUND` `USER_NOT_FOUND` `CLASS_NOT_FOUND` `STUDENT_NOT_FOUND` `PARENT_NOT_FOUND` `QUESTION_NOT_FOUND` `QUESTION_SET_NOT_FOUND` `HOMEWORK_NOT_FOUND` `SUBMISSION_NOT_FOUND` | Không tìm thấy |
| 404 | `GUARDIAN_NOT_FOUND` | Phụ huynh chưa được gắn con |
| 409 | `ALREADY_EXISTS` `EMAIL_ALREADY_EXISTS` `USER_ALREADY_EXISTS` | Email đã dùng |
| 409 | `STUDENT_CODE_EXISTS` | Trùng mã (hiếm) |
| 409 | `PARENT_ALREADY_HAS_STUDENT` | Phụ huynh đã có con khác |
| 409 | `QUESTION_IN_USE` | Câu đang dùng trong bộ/bài/nộp |
| 409 | `HOMEWORK_IN_USE` | Bài đã có em nộp — không sửa đề / không xóa |
| 409 | `SUBMISSION_ALREADY_EXISTS` | Đã nộp rồi |
| 409 | `CLASS_IN_USE` | Lớp còn bài giao — không xóa |
| 500 | `INTERNAL_SERVER_ERROR` | Lỗi hệ thống — thử lại |

401 → xóa token, về login. 403 → “Bạn không có quyền.” Toast `message` nếu chưa map.

---

## 12. Tài khoản demo

Mật khẩu chung: `Demo@123`

| Role | Email | Gợi ý test |
|---|---|---|
| admin | `admin@demo.local` | Tạo user |
| teacher | `gv.lan@demo.local` | Lớp 3A, nhiều câu/bài |
| teacher | `gv.minh@demo.local` | Lớp 5A |
| parent | `ph.01@demo.local` … `ph.20@demo.local` | Mỗi người một con |

Seed chỉ chạy khi DB **chưa có** `gv.lan@demo.local`.

Login form: có thể hiện 3 nút điền nhanh demo (admin / gv.lan / ph.01) cho tiện.

---

## 13. Cấm đoán (FE AI đọc kỹ)

1. Không endpoint nào ngoài catalog §7.
2. Không `GET /homeworks/:id` và không `GET /homework-submissions/:id` với role parent.
3. Không gọi điểm parent trước khi nộp.
4. Sửa lớp: gửi `students[].id`, không gửi `code`.
5. Gửi roster PUT = **toàn bộ** em còn giữ. Thiếu id = cho em nghỉ.
6. Ẩn đáp án parent: field **vắng** (omitempty), đừng so `=== null`.
7. Create homework / bộ / claim / nộp = `200` + `data: null`. Create câu / lớp = `201` + `data: null`. Refresh list sau đó.
8. Không màn đăng ký. Không matching. Không nộp lại / nộp trễ. Không học sinh login.
9. `due_date` không convert timezone.
10. Không tin `is_correct`/`score` lúc POST — không gửi.
11. `GET /students/by-guardian` trả **object**, không phải `Student[]`.
12. `GET .../by-guardian/:homework_id` trả **object**, không phải page.
13. Bộ câu chỉ ôn; giao bài bằng list `question.id`.
14. CORS chỉ 5173 + Netlify — dev phải port 5173 hoặc nhờ backend mở thêm.

---

## 14. Checklist nghiệm thu

- [ ] Login 3 role, guard route, logout xóa token.
- [ ] Admin tạo teacher + parent; role lạ bị chặn trên UI.
- [ ] Teacher CRUD câu MC và T/F (multipart), filter, xóa bị `QUESTION_IN_USE` khi đang dùng.
- [ ] Teacher CRUD bộ cùng loại; xóa bộ câu vẫn còn.
- [ ] Teacher tạo lớp bằng tên; chi tiết có mã; sửa roster bằng `id` + tên mới.
- [ ] Teacher search parent + claim + đổi parent; parent đã có con báo 409.
- [ ] Teacher giao bài trộn 2 loại; sửa/xóa bị chặn khi đã có nộp.
- [ ] Teacher xem điểm lớp + chi tiết từng nộp.
- [ ] Parent thấy 1 con; list bài; chưa nộp không thấy đáp án.
- [ ] Parent nộp đủ câu một lần; lần hai 409.
- [ ] Parent xem điểm đúng path `.../student/by-guardian/:homework_id`.
- [ ] Hết hạn: không nộp được, UI báo rõ.
- [ ] Không crash `GUARDIAN_NOT_FOUND`.
- [ ] UI tiếng Việt, enum đúng value API.

Làm xong: app chạy trên `http://localhost:5173`, trỏ `VITE_API_BASE_URL` tới backend, test được bằng tài khoản demo.
