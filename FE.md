    You are a senior UI/UX designer and frontend architect.

    Your task is to DESIGN A COMPLETE, PRODUCTION-READY UI/UX SYSTEM for an
    elementary-school homework management web application.

    IMPORTANT:
    - This is a UI/UX DESIGN task first.
    - Do not invent backend APIs.
    - Do not invent features that are not described below.
    - The UI must be consistent across all roles.
    - Vietnamese is the primary UI language.
    - The application must feel modern, friendly, clean and trustworthy for teachers and parents.
    - Avoid an overly corporate/admin-dashboard appearance.
    - The design should be suitable for elementary education.
    - Prioritize usability over visual decoration.

    ==================================================
    1. PRODUCT
    ==================================================

    Product type:
    Elementary School Teaching Assistant / Homework Management System.

    Roles:

    1. Admin
    2. Teacher
    3. Parent

    There is NO student login.

    Main purpose:

    ADMIN
    - Create teacher accounts
    - Create parent accounts

    TEACHER
    - Manage question bank
    - Manage question sets
    - Manage classes
    - Manage students in classes
    - Link students with parents
    - Create and manage homework
    - View student submissions and scores

    PARENT
    - View their child
    - View assigned homework
    - Complete homework
    - Submit homework once
    - View score and detailed answers after submission

    ==================================================
    2. VISUAL DIRECTION
    ==================================================

    Create a modern educational SaaS interface.

    Visual personality:
    - Friendly
    - Clean
    - Bright
    - Professional
    - Calm
    - Easy for non-technical teachers and parents
    - Suitable for elementary school

    Avoid:
    - Excessive gradients
    - Excessive glassmorphism
    - Excessive shadows
    - Neon colors
    - Gaming-style UI
    - Overly dense enterprise tables
    - Tiny text
    - Excessive animations

    Use:
    - Large readable typography
    - Rounded cards
    - Soft borders
    - Moderate shadows
    - Clear hierarchy
    - Friendly icons
    - Generous spacing
    - Clear status badges

    Suggested visual language:
    - Primary color: educational blue/indigo
    - Success: green
    - Warning: amber
    - Error: red
    - Neutral: slate/gray

    Do not hard-code colors everywhere.
    Create semantic design tokens.

    ==================================================
    3. DESIGN SYSTEM
    ==================================================

    Create a complete design system before designing individual pages.

    Define:

    Typography:
    - Display
    - H1
    - H2
    - H3
    - Body
    - Small
    - Caption
    - Label

    Spacing:
    - 4
    - 8
    - 12
    - 16
    - 20
    - 24
    - 32
    - 40
    - 48

    Border radius:
    - Small
    - Medium
    - Large
    - Full/pill

    Components:
    - Button
    - Icon button
    - Input
    - Textarea
    - Select
    - Combobox
    - Search input
    - Date picker
    - Checkbox
    - Radio
    - Switch if needed
    - Badge
    - Avatar
    - Card
    - Table
    - Pagination
    - Tabs
    - Breadcrumb
    - Dropdown menu
    - Tooltip
    - Toast
    - Alert
    - Modal/Dialog
    - Drawer
    - Confirmation dialog
    - Empty state
    - Error state
    - Skeleton loader
    - Spinner
    - Progress bar
    - File upload
    - Question card
    - Student card
    - Homework card
    - Score card

    Every reusable component must have:
    - Default
    - Hover
    - Focus
    - Disabled
    - Loading
    - Error
    - Success states where applicable

    ==================================================
    4. GLOBAL APPLICATION LAYOUT
    ==================================================

    Create a shared application shell.

    Desktop:

    ------------------------------------------------
    | Sidebar | Header                             |
    |         |------------------------------------|
    |         | Breadcrumb / Page title            |
    |         |                                    |
    |         | Main content                       |
    |         |                                    |
    ------------------------------------------------

    Sidebar:
    - Logo
    - Product name
    - Navigation
    - Active navigation state
    - Role-specific menu
    - User profile section
    - Logout

    Header:
    - Page title / breadcrumb
    - Optional search
    - User avatar
    - Username
    - Role
    - Logout menu

    Mobile:
    - Collapsible sidebar
    - Hamburger button
    - Compact header
    - Bottom navigation only if appropriate
    - Tables become cards where necessary

    Sidebar navigation MUST change according to role.

    ==================================================
    5. LOGIN
    ==================================================

    Route:

    /login

    Design a centered login page.

    Content:
    - Logo
    - Product name
    - Friendly illustration or subtle educational visual
    - "Đăng nhập"
    - Email
    - Mật khẩu
    - Show/hide password
    - Login button

    DO NOT create registration.

    Include:
    - Loading state
    - Invalid credentials state
    - Network error state

    Error:
    "Email hoặc mật khẩu không đúng."

    Optional development-only demo login section:
    - Admin
    - Giáo viên
    - Phụ huynh

    This should look like a developer/demo helper, not a production feature.

    ==================================================
    6. ADMIN UI
    ==================================================

    Route:

    /admin/users

    Admin only needs to create users.

    Design page:

    Header:
    "Quản lý tài khoản"

    Main:
    - Summary card explaining what Admin can do
    - Create user form/card

    Form:
    - Username
    - Email
    - Password
    - Role

    Role options:
    - Giáo viên
    - Phụ huynh

    DO NOT show:
    - Admin
    - Student

    After successful creation:
    - Success toast
    - Show created user ID in a small success result panel
    - Copy ID button

    States:
    - Loading
    - Validation error
    - API error
    - Success
    - Empty/default state

    ==================================================
    7. TEACHER DASHBOARD
    ==================================================

    Create a teacher home/dashboard even if it only contains information
    available from existing resources.

    Do not invent backend statistics.

    The dashboard can contain:
    - Quick navigation cards
    - Recent classes if available
    - Recent homework if available
    - Quick actions

    Quick actions:
    - Tạo câu hỏi
    - Tạo bộ câu
    - Tạo lớp
    - Giao bài

    If there is no API data:
    use useful navigation cards instead of fake statistics.

    ==================================================
    8. QUESTION BANK
    ==================================================

    Route:

    /teacher/questions

    Create a professional question bank interface.

    Top:
    - Page title
    - Search
    - Filter
    - Create button

    Filters:
    - Loại câu
    - Môn học
    - Khối
    - Độ khó
    - Tên/nội dung câu hỏi

    Question types:
    - Trắc nghiệm
    - Đúng / Sai

    DO NOT support matching questions.

    Question card/table should show:
    - Question type
    - Subject
    - Grade
    - Difficulty
    - Question content
    - Correct answer
    - Actions

    Actions:
    - View
    - Edit
    - Delete

    ==================================================
    9. QUESTION DETAIL
    ==================================================

    Create a reusable Question Detail Modal or Drawer.

    Show:
    - Question
    - Type
    - Subject
    - Grade
    - Difficulty
    - Options
    - Correct answer
    - Explanation

    For MC:
    show radio-style answer options.

    For True/False:
    show Đúng / Sai.

    ==================================================
    10. CREATE QUESTION
    ==================================================

    Route:

    /teacher/questions/new

    Create a clean form.

    Fields:
    - Type
    - Subject
    - Grade
    - Difficulty
    - Question content
    - Explanation

    For multiple choice:
    - Dynamic option list
    - Add option
    - Remove option
    - Select correct answer

    For true/false:
    - Select Đúng/Sai

    Buttons:
    - Hủy
    - Lưu câu hỏi

    Show:
    - Validation errors
    - Loading
    - Success
    - API error

    Important:
    The form must clearly indicate which answer is correct to the teacher.

    ==================================================
    11. EDIT QUESTION
    ==================================================

    Route:

    /teacher/questions/:id/edit

    Reuse the create-question layout.

    If the backend returns QUESTION_IN_USE:
    show a clear warning:

    "Câu hỏi đang được sử dụng và không thể thay đổi nội dung hoặc đáp án."

    Disable dangerous fields where appropriate.

    ==================================================
    12. DELETE QUESTION MODAL
    ==================================================

    Create confirmation dialog.

    Title:
    "Xóa câu hỏi?"

    Message:
    "Bạn có chắc chắn muốn xóa câu hỏi này?"

    Buttons:
    - Hủy
    - Xóa

    If QUESTION_IN_USE:
    show:
    "Câu hỏi đang được sử dụng trong bộ câu, bài tập hoặc bài nộp nên không thể xóa."

    ==================================================
    13. QUESTION SETS
    ==================================================

    Route:

    /teacher/question-sets

    List page.

    Features:
    - Search
    - Filter by type if appropriate
    - Create button
    - Cards/table

    Each set:
    - Title
    - Description
    - Type
    - Number of questions
    - Actions

    Actions:
    - View
    - Edit
    - Delete

    ==================================================
    14. CREATE QUESTION SET
    ==================================================

    Route:

    /teacher/question-sets/new

    Form:
    - Title
    - Description
    - Type

    After selecting type:
    show compatible questions only.

    Question selector:
    - Search
    - Filter
    - Checkbox
    - Question preview
    - Selected counter

    Show selected questions in a separate panel.

    Desktop:
    Two-column layout:

    Available questions | Selected questions

    Mobile:
    Stack vertically.

    Do NOT create a "assign entire set" action.
    The set is only a reusable collection for reviewing/selecting questions.

    ==================================================
    15. QUESTION SET DETAIL
    ==================================================

    Route:

    /teacher/question-sets/:id

    Show:
    - Set title
    - Description
    - Type
    - Questions

    Actions:
    - Edit
    - Delete

    Question list:
    - Number
    - Type
    - Subject
    - Grade
    - Difficulty
    - Content
    - Preview

    Delete set:
    The UI must communicate that deleting a set does NOT delete its questions.

    ==================================================
    16. CLASSES
    ==================================================

    Route:

    /teacher/classes

    List page.

    Features:
    - Search by class name
    - Create class
    - Class cards

    Each class card:
    - Class image
    - Class name
    - Description
    - Student count
    - View button

    Empty state:
    "Bạn chưa có lớp học nào."

    CTA:
    "Tạo lớp đầu tiên"

    ==================================================
    17. CREATE CLASS
    ==================================================

    Route:

    /teacher/classes/new

    Fields:
    - Class name
    - Description
    - Image
    - Student names

    Student entry should support:
    - Textarea
    - Chips/tags
    - One student per line

    Explain:
    "Mỗi dòng là tên một học sinh."

    Image:
    - Drag & drop
    - Browse
    - Preview
    - Remove
    - Maximum 5MB

    Buttons:
    - Hủy
    - Tạo lớp

    ==================================================
    18. CLASS DETAIL
    ==================================================

    Route:

    /teacher/classes/:id

    Create a rich class management page.

    Header:
    - Class image
    - Class name
    - Description
    - Edit button
    - Delete button

    Tabs:

    1. Học sinh
    2. Bài tập

    Student table:

    - Học sinh
    - Mã học sinh
    - Phụ huynh
    - Trạng thái
    - Actions

    Student code:
    - Copy button
    - Tooltip
    - Success toast after copying

    Parent:
    - Email if linked
    - "Chưa gắn" if not linked

    ==================================================
    19. LINK PARENT MODAL
    ==================================================

    Create modal:

    Title:
    "Gắn phụ huynh"

    Fields:
    - Mã học sinh
    - Search parent

    Parent search:
    - Search by email/name
    - Dropdown results
    - Email
    - Name

    Buttons:
    - Hủy
    - Gắn phụ huynh

    For changing parent:
    same modal can become:

    "Đổi phụ huynh"

    If PARENT_ALREADY_HAS_STUDENT:
    show an inline error explaining the parent already has another child.

    ==================================================
    20. EDIT CLASS ROSTER
    ==================================================

    Create Edit Roster modal/drawer.

    Show existing students as selectable rows.

    Each row:
    - Checkbox
    - Student name
    - Hidden student ID conceptually
    - Status

    Add student:
    - Input name
    - Add button

    Important UX:
    Explain clearly:

    "Bỏ chọn học sinh sẽ đưa học sinh ra khỏi danh sách lớp."

    Do not expose technical IDs to the teacher.

    ==================================================
    21. DELETE CLASS MODAL
    ==================================================

    Confirmation dialog.

    If CLASS_IN_USE:
    show:

    "Lớp này đang có bài tập được giao nên chưa thể xóa."

    ==================================================
    22. HOMEWORK LIST
    ==================================================

    Route:

    /teacher/homeworks

    Page:

    - Search
    - Filter
    - Create homework

    Homework card/table:
    - Title
    - Class
    - Due date
    - Number of questions
    - Status
    - Actions

    Status examples:
    - Chưa đến hạn
    - Đã hết hạn
    - Đã có bài nộp

    Do not invent unsupported backend status fields.
    Derive display only from available data.

    ==================================================
    23. CREATE HOMEWORK
    ==================================================

    Route:

    /teacher/homeworks/new

    Use a step-like or well-organized single-page form.

    Section 1:
    Thông tin bài tập

    - Class
    - Title
    - Description
    - Due date

    Section 2:
    Chọn câu hỏi

    - Search
    - Filters
    - Question list
    - Checkbox
    - Selected counter

    Allow mixing:
    - Multiple choice
    - True/False

    Show question type clearly.

    Optional helper:
    "Chọn câu từ bộ câu"

    But selecting a set is only a helper to select individual questions.
    Do not create an API action for assigning an entire set.

    Bottom sticky action bar:
    - Hủy
    - Tạo bài tập

    ==================================================
    24. HOMEWORK DETAIL
    ==================================================

    Route:

    /teacher/homeworks/:id

    Header:
    - Title
    - Class
    - Due date
    - Description

    Question list:
    - Question number
    - Type
    - Content
    - Correct answer
    - Explanation

    Actions:
    - Edit
    - Delete
    - View submissions

    If submissions already exist:
    Disable:
    - Changing question list
    - Delete homework

    Show warning:

    "Bài tập đã có học sinh nộp nên không thể thay đổi đề hoặc xóa."

    ==================================================
    25. DELETE HOMEWORK MODAL
    ==================================================

    Confirmation modal.

    If no submissions:
    normal delete confirmation.

    If HOMEWORK_IN_USE:
    show locked state and explanation.

    ==================================================
    26. CLASS SUBMISSIONS / SCORES
    ==================================================

    Route:

    /teacher/homeworks/:id/submissions

    Create score dashboard.

    Header:
    - Homework title
    - Class
    - Due date

    Summary area:
    - Number of submissions if available
    - Average only if backend data supports it
    - Other statistics only when actual API data exists

    Main table:
    - Student name
    - Score / 100
    - Submitted time
    - Status

    Click a student:
    open submission detail.

    ==================================================
    27. SUBMISSION DETAIL
    ==================================================

    Create a full-page detail or large drawer.

    Header:
    - Student name
    - Homework
    - Score

    Large score display:

    85.5 / 100

    Then question-by-question review.

    Each question:
    - Question
    - Student answer
    - Correct answer
    - Correct/incorrect indicator
    - Score
    - Max score
    - Explanation

    Use strong visual hierarchy.

    Correct:
    positive state.

    Incorrect:
    error/warning state.

    Do not overwhelm the teacher with colors.

    ==================================================
    28. PARENT HOME
    ==================================================

    Route:

    /parent

    Parent dashboard.

    Top:
    - Greeting
    - Child information

    Child card:
    - Student name
    - Student code
    - Class ID/name if available

    Homework section:
    - Homework cards

    Each card:
    - Title
    - Description
    - Due date
    - Number of questions
    - Submission status
    - CTA

    Statuses:
    - Chưa làm
    - Đã nộp
    - Đã hết hạn

    If GUARDIAN_NOT_FOUND:

    Create a friendly empty state:

    "Chưa có học sinh được gắn với tài khoản này."

    "Nhờ giáo viên gắn mã học sinh của con vào tài khoản phụ huynh."

    Do NOT show homework form.

    ==================================================
    29. PARENT HOMEWORK
    ==================================================

    Route:

    /parent/homeworks/:id

    This is an important user experience.

    Create a distraction-free homework interface.

    Header:
    - Homework title
    - Class
    - Due date
    - Progress

    Example:

    Câu 5 / 10

    Progress bar.

    Question card:

    "Câu 5"

    Question content

    For multiple choice:
    large selectable answer cards.

    For true/false:
    two large buttons/cards:
    - Đúng
    - Sai

    Do NOT show:
    - Correct answer
    - Explanation
    - Score

    until submission.

    Navigation:
    - Câu trước
    - Câu tiếp theo

    Sticky bottom action:
    - Tiến độ
    - Nộp bài

    ==================================================
    30. SUBMIT HOMEWORK CONFIRMATION MODAL
    ==================================================

    Before submission:

    Title:
    "Nộp bài?"

    Show:
    - Number of answered questions
    - Number of unanswered questions

    If incomplete:
    "Nếu chưa trả lời hết, bạn chưa thể nộp bài."

    Disable submit until all questions answered.

    Once complete:

    "Bạn đã trả lời tất cả câu hỏi. Bạn có chắc muốn nộp bài?"

    Buttons:
    - Xem lại
    - Nộp bài

    After submit:
    redirect to result.

    ==================================================
    31. PARENT RESULT
    ==================================================

    Create a result page.

    Large score:

    85.5 / 100

    Friendly result header.

    Show:
    - Homework title
    - Submission time

    Then question review.

    Each question:
    - Question
    - Parent/student selected answer
    - Correct answer
    - Correct/incorrect
    - Score
    - Explanation

    Use a clear learning-oriented layout.

    This page should feel encouraging rather than punitive.

    ==================================================
    32. GLOBAL MODALS
    ==================================================

    Design all required reusable modal types:

    1. Confirm delete
    2. Confirm submit
    3. Create/edit
    4. Link parent
    5. Change parent
    6. View question
    7. View submission
    8. Error dialog
    9. Unsaved changes
    10. Image upload
    11. Logout confirmation if appropriate

    Modal rules:
    - Clear title
    - Short explanation
    - Primary action
    - Secondary action
    - Escape closes when safe
    - Destructive actions visually separated
    - Loading state during submission
    - Prevent duplicate submission

    ==================================================
    33. DRAWERS
    ==================================================

    Use drawers for:
    - Question preview
    - Submission detail on desktop
    - Student detail
    - Filters on mobile

    Do not turn every interaction into a modal.

    Use full pages for complex forms.

    ==================================================
    34. TOAST SYSTEM
    ==================================================

    Create consistent toast notifications.

    Success:
    "Cập nhật thành công."

    "Đã tạo câu hỏi."

    "Đã tạo lớp."

    "Đã giao bài."

    "Đã nộp bài."

    Error:
    Use backend message when appropriate.

    Never show raw HTTP errors.

    Examples:
    - "Bạn không có quyền thực hiện thao tác này."
    - "Câu hỏi đang được sử dụng."
    - "Bài tập đã có học sinh nộp."
    - "Phụ huynh này đã được gắn với học sinh khác."

    ==================================================
    35. ERROR STATES
    ==================================================

    Design reusable error states for:

    401:
    "Phiên đăng nhập đã hết. Vui lòng đăng nhập lại."

    403:
    "Bạn không có quyền truy cập nội dung này."

    404:
    "Không tìm thấy dữ liệu."

    409:
    Context-specific message.

    500:
    "Đã xảy ra lỗi hệ thống. Vui lòng thử lại."

    Network:
    "Không thể kết nối đến máy chủ."

    ==================================================
    36. LOADING STATES
    ==================================================

    Every list and detail page MUST have loading state.

    Use skeletons rather than full-page spinners whenever possible.

    Create:
    - Table skeleton
    - Card skeleton
    - Question skeleton
    - Homework skeleton
    - Student skeleton
    - Form loading state

    Buttons during request:
    - Disable
    - Show spinner
    - Preserve button text context

    Example:
    "Lưu..." instead of "Lưu"

    ==================================================
    37. EMPTY STATES
    ==================================================

    Every list must have a meaningful empty state.

    Examples:

    Questions:
    "Chưa có câu hỏi nào."

    Question sets:
    "Chưa có bộ câu nào."

    Classes:
    "Bạn chưa có lớp học nào."

    Homework:
    "Chưa có bài tập nào."

    Submissions:
    "Chưa có học sinh nào nộp bài."

    Use:
    - Small illustration/icon
    - Explanation
    - CTA when appropriate

    ==================================================
    38. RESPONSIVE DESIGN
    ==================================================

    Desktop:
    Optimized for 1440px.

    Tablet:
    Optimized for 768px.

    Mobile:
    Optimized for 375px and 390px.

    Rules:
    - Sidebar collapses
    - Tables become cards
    - Modals become bottom sheets or full-screen dialogs when appropriate
    - Forms become one-column
    - Sticky actions remain accessible
    - Touch targets minimum 44px

    ==================================================
    39. ACCESSIBILITY
    ==================================================

    Follow accessible UI principles.

    Requirements:
    - Keyboard navigation
    - Visible focus states
    - Proper labels
    - Accessible dialogs
    - Accessible radio/checkbox
    - Sufficient contrast
    - Do not rely only on color to indicate correctness
    - Screen-reader-friendly buttons
    - Tooltips for icon-only buttons

    ==================================================
    40. UX SAFETY
    ==================================================

    Dangerous operations require confirmation:

    - Delete question
    - Delete question set
    - Delete class
    - Delete homework
    - Submit homework

    Do not allow accidental duplicate submission.

    Do not allow destructive actions while API request is running.

    ==================================================
    41. API/DOMAIN CONSTRAINTS
    ==================================================

    The UI MUST respect these backend rules.

    Roles:
    admin | teacher | parent

    Question types:
    multiple_choice | true_false

    Subjects:
    vietnamese
    mathematics
    ethics
    english
    nature_and_society
    history_and_geography
    science
    informatics
    technology
    physical_education
    music
    art
    experiential_activities

    Grades:
    "1" | "2" | "3" | "4" | "5"

    Difficulty:
    easy | medium | hard

    UI labels must be Vietnamese.

    Do not invent:
    - matching
    - student login
    - registration
    - retake
    - late submission
    - assigning an entire question set as a homework API action

    ==================================================
    42. IMPORTANT SECURITY / DATA RULES
    ==================================================

    Parent must NEVER see answers before submission.

    Before submission:
    Do not render:
    - correct_index
    - correct_bool
    - explanation

    After submission:
    These may be shown.

    Parent must not call teacher-only submission endpoints.

    Parent score page should only load after submission.

    ==================================================
    43. PAGE INVENTORY
    ==================================================

    Design ALL of these screens:

    PUBLIC
    1. Login

    ADMIN
    2. Admin user management

    TEACHER
    3. Teacher dashboard
    4. Question bank
    5. Question detail
    6. Create question
    7. Edit question
    8. Question set list
    9. Create question set
    10. Question set detail
    11. Class list
    12. Create class
    13. Class detail
    14. Homework list
    15. Create homework
    16. Homework detail
    17. Homework submissions
    18. Submission detail

    PARENT
    19. Parent dashboard
    20. Homework list/detail
    21. Homework answering
    22. Submit confirmation
    23. Homework result

    SYSTEM STATES
    24. Loading
    25. Empty
    26. Error
    27. Unauthorized
    28. Forbidden
    29. Not found
    30. Network error
    31. Success toast
    32. Validation errors

    MODALS / DRAWERS
    33. Delete confirmation
    34. Submit confirmation
    35. Link parent
    36. Change parent
    37. Question preview
    38. Submission preview
    39. Unsaved changes
    40. Logout confirmation
    41. Image upload

    ==================================================
    44. FINAL DESIGN DELIVERABLE
    ==================================================

    Generate a COMPLETE UI DESIGN SYSTEM and all application screens.

    Do not only design the happy path.

    For every important screen include:
    - Default state
    - Loading state
    - Empty state
    - Error state
    - Validation state
    - Success state
    - Confirmation state where relevant

    Make all pages visually consistent.

    Use realistic Vietnamese sample content.

    Do not use Lorem Ipsum.

    Example:
    "Lớp 3A"
    "Toán cuối tuần"
    "15 + 27 bằng bao nhiêu?"
    "Nguyễn Văn An"

    The final result should look like a real production application,
    not a generic AI-generated dashboard.

    Most important:
    UX clarity > visual effects.
    Consistency > decoration.
    Real educational workflow > fake statistics.
    Respect the backend/domain rules exactly.