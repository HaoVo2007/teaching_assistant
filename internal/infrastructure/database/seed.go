package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	seedMarkerEmail = "gv.lan@demo.local"
	seedPassword    = "Demo@123"
)

func Seed(ctx context.Context, db *mongo.Database) error {
	count, err := db.Collection("users").CountDocuments(ctx, bson.M{"email": seedMarkerEmail})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	now := time.Date(2026, 9, 16, 8, 0, 0, 0, time.UTC)
	hash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	password := string(hash)

	adminID := primitive.NewObjectID()
	teacherLanID := primitive.NewObjectID()
	teacherMinhID := primitive.NewObjectID()
	teacherLan := teacherLanID.Hex()
	teacherMinh := teacherMinhID.Hex()

	parentIDs := make([]primitive.ObjectID, 20)
	for i := range parentIDs {
		parentIDs[i] = primitive.NewObjectID()
	}

	users := []any{
		newUser(adminID, "admin", "admin@demo.local", password, user.RoleAdmin, now),
		newUser(teacherLanID, "nguyen_thi_lan", seedMarkerEmail, password, user.RoleUser, now),
		newUser(teacherMinhID, "tran_van_minh", "gv.minh@demo.local", password, user.RoleUser, now),
	}
	for i, id := range parentIDs {
		users = append(users, newUser(
			id,
			fmt.Sprintf("phu_huynh_%02d", i+1),
			fmt.Sprintf("ph.%02d@demo.local", i+1),
			password,
			user.RoleParent,
			now,
		))
	}
	if _, err := db.Collection("users").InsertMany(ctx, users); err != nil {
		return err
	}

	class3A := primitive.NewObjectID()
	class5A := primitive.NewObjectID()

	names3A := studentNames3A()
	students3A := make([]*student.Student, 0, len(names3A))
	studentIDs3A := make([]string, 0, len(names3A))
	for i, name := range names3A {
		id := primitive.NewObjectID()
		status := student.StudentStatusActive
		if i >= 28 {
			status = student.StudentStatusInactive
		}
		students3A = append(students3A, &student.Student{
			ID:        id,
			Code:      fmt.Sprintf("HS260916-%02d", i+1),
			Name:      name,
			ClassID:   class3A.Hex(),
			Status:    status,
			CreatedBy: teacherLan,
			CreatedAt: now.Add(time.Duration(i) * time.Minute),
			UpdatedAt: now,
		})
		studentIDs3A = append(studentIDs3A, id.Hex())
	}

	names5A := []string{
		"Nguyễn Nhật Minh", "Trần Khánh Linh", "Lê Đức Anh", "Phạm Mỹ Duyên",
		"Hoàng Gia Bảo", "Vũ Ngọc Hà", "Đặng Thanh Tùng", "Bùi Phương Anh",
		"Ngô Hải Nam", "Đỗ Thu Trang", "Lý Quốc Huy", "Mai Anh Thư",
	}
	students5A := make([]*student.Student, 0, len(names5A))
	studentIDs5A := make([]string, 0, len(names5A))
	for i, name := range names5A {
		id := primitive.NewObjectID()
		students5A = append(students5A, &student.Student{
			ID:        id,
			Code:      fmt.Sprintf("HS260916-B%02d", i+1),
			Name:      name,
			ClassID:   class5A.Hex(),
			Status:    student.StudentStatusActive,
			CreatedBy: teacherMinh,
			CreatedAt: now,
			UpdatedAt: now,
		})
		studentIDs5A = append(studentIDs5A, id.Hex())
	}

	allStudents := make([]any, 0, len(students3A)+len(students5A))
	for _, s := range students3A {
		allStudents = append(allStudents, s)
	}
	for _, s := range students5A {
		allStudents = append(allStudents, s)
	}
	if _, err := db.Collection("students").InsertMany(ctx, allStudents); err != nil {
		return err
	}

	classes := []any{
		&class.Class{
			ID:          class3A,
			Name:        "Lớp 3A",
			Description: "Lớp 3A năm học 2026-2027 — 30 học sinh, khối 3",
			Students:    studentIDs3A,
			CreatedBy:   teacherLan,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		&class.Class{
			ID:          class5A,
			Name:        "Lớp 5A",
			Description: "Lớp 5A năm học 2026-2027 — giáo viên Trần Văn Minh",
			Students:    studentIDs5A,
			CreatedBy:   teacherMinh,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
	if _, err := db.Collection("classes").InsertMany(ctx, classes); err != nil {
		return err
	}

	guardians := make([]any, 0, 24)
	roles := []string{"mother", "father", "guardian", "mother", "father"}
	for i := 0; i < 18; i++ {
		guardians = append(guardians, &student.Guardian{
			ID:        primitive.NewObjectID(),
			ParentID:  parentIDs[i].Hex(),
			StudentID: studentIDs3A[i],
			Role:      roles[i%len(roles)],
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	guardians = append(guardians,
		&student.Guardian{
			ID:        primitive.NewObjectID(),
			ParentID:  parentIDs[0].Hex(),
			StudentID: studentIDs3A[1],
			Role:      "mother",
			CreatedAt: now,
			UpdatedAt: now,
		},
		&student.Guardian{
			ID:        primitive.NewObjectID(),
			ParentID:  parentIDs[18].Hex(),
			StudentID: studentIDs5A[0],
			Role:      "father",
			CreatedAt: now,
			UpdatedAt: now,
		},
	)
	if _, err := db.Collection("guardians").InsertMany(ctx, guardians); err != nil {
		return err
	}

	questions := buildQuestions(teacherLan, teacherMinh, now)
	qDocs := make([]any, 0, len(questions))
	for _, q := range questions {
		qDocs = append(qDocs, q)
	}
	if _, err := db.Collection("questions").InsertMany(ctx, qDocs); err != nil {
		return err
	}

	lanMC := idsBy(questions, teacherLan, string(question.QuestionTypeMultipleChoice), 8)
	lanTF := idsBy(questions, teacherLan, string(question.QuestionTypeTrueFalse), 8)
	lanMatch := idsBy(questions, teacherLan, string(question.QuestionTypeMatching), 5)
	minhMC := idsBy(questions, teacherMinh, string(question.QuestionTypeMultipleChoice), 6)

	desc := func(s string) *string { return &s }
	sets := []any{
		&questionset.QuestionSet{
			ID:           primitive.NewObjectID(),
			Title:        "Trắc nghiệm Toán — Lớp 3",
			Description:  desc("Ôn tập phép tính và hình học"),
			QuestionType: questionset.QuestionSetTypeMultipleChoice,
			QuestionIds:  lanMC,
			CreatedBy:    teacherLan,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		&questionset.QuestionSet{
			ID:           primitive.NewObjectID(),
			Title:        "Đúng / Sai Khoa học",
			Description:  desc("Tự nhiên và xã hội"),
			QuestionType: questionset.QuestionSetTypeTrueFalse,
			QuestionIds:  lanTF,
			CreatedBy:    teacherLan,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		&questionset.QuestionSet{
			ID:           primitive.NewObjectID(),
			Title:        "Matching từ vựng tiếng Anh",
			Description:  desc("Ghép từ với nghĩa"),
			QuestionType: questionset.QuestionSetTypeMatching,
			QuestionIds:  lanMatch,
			CreatedBy:    teacherLan,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
		&questionset.QuestionSet{
			ID:           primitive.NewObjectID(),
			Title:        "Trắc nghiệm Lớp 5",
			Description:  desc("Ngân hàng của giáo viên Minh"),
			QuestionType: questionset.QuestionSetTypeMultipleChoice,
			QuestionIds:  minhMC,
			CreatedBy:    teacherMinh,
			CreatedAt:    now,
			UpdatedAt:    now,
		},
	}
	if _, err := db.Collection("question_sets").InsertMany(ctx, sets); err != nil {
		return err
	}

	hwToan := primitive.NewObjectID()
	hwTV := primitive.NewObjectID()
	hwTA := primitive.NewObjectID()
	hwMix := primitive.NewObjectID()
	hwOverdue := primitive.NewObjectID()
	hw5A := primitive.NewObjectID()

	mcIDs := questionIDsByType(questions, teacherLan, string(question.QuestionTypeMultipleChoice))
	tfIDs := questionIDsByType(questions, teacherLan, string(question.QuestionTypeTrueFalse))
	minhIDs := questionIDsByType(questions, teacherMinh, string(question.QuestionTypeMultipleChoice))

	homeworks := []any{
		&homework.Homework{
			ID: hwToan, ClassID: class3A.Hex(), Title: "Bài tập Toán tuần 1",
			Description: desc("Phép cộng, trừ, nhân — 8 câu"),
			Questions:   take(mcIDs, 0, 8), DueDate: now.Add(10 * 24 * time.Hour),
			CreatedBy: teacherLan, CreatedAt: now, UpdatedAt: now,
		},
		&homework.Homework{
			ID: hwTV, ClassID: class3A.Hex(), Title: "Tiếng Việt — Đúng hay sai",
			Description: desc("Từ loại và chính tả"),
			Questions:   take(tfIDs, 0, 6), DueDate: now.Add(14 * 24 * time.Hour),
			CreatedBy: teacherLan, CreatedAt: now, UpdatedAt: now,
		},
		&homework.Homework{
			ID: hwTA, ClassID: class3A.Hex(), Title: "English vocabulary",
			Description: desc("Trắc nghiệm từ vựng"),
			Questions:   take(mcIDs, 8, 6), DueDate: now.Add(7 * 24 * time.Hour),
			CreatedBy: teacherLan, CreatedAt: now, UpdatedAt: now,
		},
		&homework.Homework{
			ID: hwMix, ClassID: class3A.Hex(), Title: "Ôn tập giữa kỳ",
			Description: desc("Trộn trắc nghiệm và đúng/sai"),
			Questions:   append(append([]string{}, take(mcIDs, 2, 4)...), take(tfIDs, 2, 4)...),
			DueDate:     now.Add(21 * 24 * time.Hour),
			CreatedBy:   teacherLan, CreatedAt: now, UpdatedAt: now,
		},
		&homework.Homework{
			ID: hwOverdue, ClassID: class3A.Hex(), Title: "Bài kiểm tra tuần trước",
			Description: desc("Đã quá hạn — dùng xem bài nộp cũ"),
			Questions:   take(mcIDs, 4, 5), DueDate: now.Add(-5 * 24 * time.Hour),
			CreatedBy: teacherLan, CreatedAt: now.Add(-12 * 24 * time.Hour), UpdatedAt: now,
		},
		&homework.Homework{
			ID: hw5A, ClassID: class5A.Hex(), Title: "Toán lớp 5 — phân số",
			Description: desc("Bài giao cho lớp 5A"),
			Questions:   take(minhIDs, 0, 5), DueDate: now.Add(8 * 24 * time.Hour),
			CreatedBy: teacherMinh, CreatedAt: now, UpdatedAt: now,
		},
	}
	if _, err := db.Collection("homeworks").InsertMany(ctx, homeworks); err != nil {
		return err
	}

	qByID := map[string]*question.Question{}
	for _, q := range questions {
		qByID[q.ID.Hex()] = q
	}

	submissions := make([]any, 0, 80)
	submissions = append(submissions, seedSubmissions(hwToan.Hex(), teacherLan, take(mcIDs, 0, 8), students3A[:24], qByID, now.Add(2*time.Hour))...)
	submissions = append(submissions, seedSubmissions(hwTV.Hex(), teacherLan, take(tfIDs, 0, 6), students3A[:16], qByID, now.Add(26*time.Hour))...)
	submissions = append(submissions, seedSubmissions(hwTA.Hex(), teacherLan, take(mcIDs, 8, 6), students3A[4:20], qByID, now.Add(30*time.Hour))...)
	submissions = append(submissions, seedSubmissions(hwOverdue.Hex(), teacherLan, take(mcIDs, 4, 5), students3A[:28], qByID, now.Add(-6*24*time.Hour))...)
	submissions = append(submissions, seedSubmissions(hw5A.Hex(), teacherMinh, take(minhIDs, 0, 5), students5A[:8], qByID, now.Add(3*time.Hour))...)
	if _, err := db.Collection("homework_submissions").InsertMany(ctx, submissions); err != nil {
		return err
	}

	log.Printf("seeded demo data (password %s): teacher %s / %s, 30 students class 3A codes HS260916-01..30, parents ph.01@demo.local..",
		seedPassword, seedMarkerEmail, "gv.minh@demo.local")
	return nil
}

func newUser(id primitive.ObjectID, username, email, password string, role user.Role, now time.Time) *user.User {
	return &user.User{
		ID: id, Username: username, Email: email, Password: password, Role: role,
		CreatedAt: now, UpdatedAt: now,
	}
}

func studentNames3A() []string {
	return []string{
		"Nguyễn Minh An", "Trần Thị Bảo", "Lê Hoàng Cường", "Phạm Ngọc Dung", "Hoàng Văn Đức",
		"Vũ Thị Hà", "Đặng Minh Hải", "Bùi Thị Hoa", "Ngô Văn Hùng", "Đỗ Thị Hương",
		"Lý Minh Khoa", "Mai Thị Lan", "Phan Văn Long", "Trịnh Thị Mai", "Dương Minh Nam",
		"Lương Thị Ngọc", "Cao Văn Phát", "Tô Thị Quỳnh", "Hồ Minh Sơn", "Châu Thị Trang",
		"Võ Văn Thành", "Đinh Thị Uyên", "Lâm Minh Vũ", "Kiều Thị Yến", "Tạ Văn Bình",
		"Từ Thị Cẩm", "Ông Minh Đạt", "Quách Thị Em", "Huỳnh Văn Phong", "Lưu Thị Sen",
	}
}

func ptrInt(v int) *int    { return &v }
func ptrBool(v bool) *bool { return &v }

func mc(createdBy, subject, grade, difficulty, prompt string, options []string, correct int, explanation string, now time.Time) *question.Question {
	return &question.Question{
		ID: primitive.NewObjectID(), Type: string(question.QuestionTypeMultipleChoice),
		Subject: subject, Grade: grade, Difficulty: difficulty, Question: prompt,
		Options: options, CorrectIndex: ptrInt(correct), Explanation: explanation,
		CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now,
	}
}

func tf(createdBy, subject, grade, difficulty, prompt string, correct bool, explanation string, now time.Time) *question.Question {
	return &question.Question{
		ID: primitive.NewObjectID(), Type: string(question.QuestionTypeTrueFalse),
		Subject: subject, Grade: grade, Difficulty: difficulty, Question: prompt,
		CorrectBool: ptrBool(correct), Explanation: explanation,
		CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now,
	}
}

func matching(createdBy, subject, grade, difficulty, prompt string, pairs []question.Pair, explanation string, now time.Time) *question.Question {
	return &question.Question{
		ID: primitive.NewObjectID(), Type: string(question.QuestionTypeMatching),
		Subject: subject, Grade: grade, Difficulty: difficulty, Question: prompt,
		Pairs: pairs, Explanation: explanation,
		CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now,
	}
}

func textPair(left, right string) question.Pair {
	return question.Pair{Left: left, LeftKind: string(question.Text), Right: right, RightKind: string(question.Text)}
}

func buildQuestions(teacherLan, teacherMinh string, now time.Time) []*question.Question {
	g3, g5 := string(question.Grade3), string(question.Grade5)
	math, viet, eng := string(question.SubjectMathematics), string(question.SubjectVietnamese), string(question.SubjectEnglish)
	sci, eth := string(question.SubjectScience), string(question.SubjectEthics)
	easy, med, hard := string(question.DifficultyEasy), string(question.DifficultyMedium), string(question.DifficultyHard)

	qs := []*question.Question{
		mc(teacherLan, math, g3, easy, "15 + 27 bằng bao nhiêu?", []string{"32", "42", "41", "52"}, 1, "15+27=42", now),
		mc(teacherLan, math, g3, easy, "9 × 6 bằng bao nhiêu?", []string{"54", "56", "45", "63"}, 0, "9×6=54", now),
		mc(teacherLan, math, g3, easy, "100 − 37 bằng bao nhiêu?", []string{"73", "63", "67", "74"}, 0, "100-37=63", now),
		mc(teacherLan, math, g3, med, "Một hình chữ nhật dài 8cm, rộng 5cm. Chu vi là?", []string{"13cm", "26cm", "40cm", "18cm"}, 1, "Chu vi = 2×(8+5)=26", now),
		mc(teacherLan, math, g3, med, "1 giờ = ? phút", []string{"30", "60", "90", "100"}, 1, "1 giờ = 60 phút", now),
		mc(teacherLan, math, g3, hard, "An có 48 viên bi, cho bạn 1/3 số bi. An còn lại?", []string{"16", "24", "32", "36"}, 2, "48 - 16 = 32", now),
		mc(teacherLan, math, g3, med, "Số liền sau của 199 là?", []string{"198", "200", "190", "201"}, 1, "199+1=200", now),
		mc(teacherLan, math, g3, easy, "Hình nào có 3 cạnh?", []string{"Vuông", "Tròn", "Tam giác", "Chữ nhật"}, 2, "Tam giác có 3 cạnh", now),
		mc(teacherLan, viet, g3, easy, "Từ nào viết đúng chính tả?", []string{"ngành nghề", "ngành nghê", "ngành nghe", "nghanh nghề"}, 0, "ngành nghề", now),
		mc(teacherLan, viet, g3, med, "Từ nào là danh từ?", []string{"chạy", "học sinh", "đẹp", "rất"}, 1, "học sinh là danh từ", now),
		mc(teacherLan, viet, g3, med, "Dấu câu thích hợp cuối câu cảm: \"Trời ơi đẹp quá\"", []string{".", "?", "!", ","}, 2, "Câu cảm dùng dấu chấm than", now),
		mc(teacherLan, viet, g3, hard, "Từ nào trái nghĩa với \"chăm chỉ\"?", []string{"siêng năng", "lười biếng", "cần cù", "kiên trì"}, 1, "trái nghĩa là lười biếng", now),
		mc(teacherLan, eng, g3, easy, "\"Cat\" nghĩa là gì?", []string{"Chó", "Mèo", "Gà", "Cá"}, 1, "cat = mèo", now),
		mc(teacherLan, eng, g3, easy, "Chọn từ đúng: This ___ a book.", []string{"am", "is", "are", "be"}, 1, "This is a book", now),
		mc(teacherLan, eng, g3, med, "Số 12 trong tiếng Anh là?", []string{"twenty", "eleven", "twelve", "ten"}, 2, "twelve", now),
		mc(teacherLan, eng, g3, med, "\"Good morning\" dùng khi nào?", []string{"Buổi tối", "Buổi sáng", "Buổi đêm", "Tạm biệt"}, 1, "chào buổi sáng", now),
		mc(teacherLan, sci, g3, easy, "Cây xanh cần gì để quang hợp?", []string{"Chỉ nước", "Ánh sáng mặt trời", "Chỉ đất", "Muối"}, 1, "cần ánh sáng", now),
		mc(teacherLan, sci, g3, med, "Cơ quan nào bơm máu đi khắp cơ thể?", []string{"Phổi", "Gan", "Tim", "Dạ dày"}, 2, "tim", now),
		mc(teacherLan, eth, g3, easy, "Khi được người khác giúp, em nên nói?", []string{"Im lặng", "Cảm ơn", "Đi ngay", "Phàn nàn"}, 1, "nói cảm ơn", now),
		mc(teacherLan, string(question.SubjectInformatics), g3, med, "Chuột máy tính dùng để?", []string{"Nấu ăn", "Điều khiển con trỏ", "Thổi gió", "Sạc pin tai"}, 1, "điều khiển con trỏ", now),

		tf(teacherLan, math, g3, easy, "Số 0 là số chẵn.", true, "0 chia hết cho 2", now),
		tf(teacherLan, math, g3, easy, "Tam giác có 4 góc.", false, "tam giác có 3 góc", now),
		tf(teacherLan, math, g3, med, "1km = 1000m.", true, "đúng", now),
		tf(teacherLan, viet, g3, easy, "\"Hoa\" là danh từ.", true, "hoa là tên sự vật", now),
		tf(teacherLan, viet, g3, med, "Câu hỏi luôn kết thúc bằng dấu chấm.", false, "câu hỏi dùng dấu hỏi", now),
		tf(teacherLan, sci, g3, easy, "Mặt Trời mọc ở hướng Tây.", false, "mọc ở hướng Đông", now),
		tf(teacherLan, sci, g3, med, "Nước sôi ở 100°C (điều kiện thường).", true, "đúng", now),
		tf(teacherLan, eng, g3, easy, "\"I\" đi với động từ \"am\".", true, "I am", now),
		tf(teacherLan, eng, g3, med, "\"They is students\" là câu đúng.", false, "They are students", now),
		tf(teacherLan, eth, g3, easy, "Xếp hàng khi mua đồ là việc nên làm.", true, "lịch sự", now),
		tf(teacherLan, string(question.SubjectPhysicalEducation), g3, easy, "Chơi thể thao giúp cơ thể khỏe.", true, "đúng", now),
		tf(teacherLan, string(question.SubjectMusic), g3, easy, "Nốt Đồ là nốt thấp nhất trong gam Đô trưởng.", true, "C là nốt chủ", now),

		matching(teacherLan, eng, g3, easy, "Ghép từ tiếng Anh với nghĩa tiếng Việt", []question.Pair{
			textPair("Apple", "Quả táo"), textPair("Dog", "Con chó"), textPair("School", "Trường học"),
		}, "từ vựng cơ bản", now),
		matching(teacherLan, math, g3, med, "Ghép hình với số cạnh", []question.Pair{
			textPair("Tam giác", "3"), textPair("Tứ giác", "4"), textPair("Ngũ giác", "5"),
		}, "số cạnh", now),
		matching(teacherLan, viet, g3, easy, "Ghép từ với loại từ", []question.Pair{
			textPair("học sinh", "danh từ"), textPair("chạy", "động từ"), textPair("đẹp", "tính từ"),
		}, "từ loại", now),
		matching(teacherLan, sci, g3, med, "Ghép động vật với nhóm", []question.Pair{
			textPair("Cá", "Sống dưới nước"), textPair("Chim", "Có cánh"), textPair("Mèo", "Thú nuôi"),
		}, "động vật", now),
		matching(teacherLan, string(question.SubjectHistoryAndGeography), g3, hard, "Ghép địa danh", []question.Pair{
			textPair("Hà Nội", "Thủ đô"), textPair("Huế", "Cố đô"), textPair("Hạ Long", "Vịnh"),
		}, "địa lý Việt Nam", now),
	}

	for i := 0; i < 10; i++ {
		a, b := 11+i, 7+i
		qs = append(qs, mc(teacherLan, math, g3, easy,
			fmt.Sprintf("%d + %d = ?", a, b),
			[]string{fmt.Sprintf("%d", a+b-2), fmt.Sprintf("%d", a+b), fmt.Sprintf("%d", a+b+3), fmt.Sprintf("%d", a*b)},
			1, fmt.Sprintf("%d+%d=%d", a, b, a+b), now))
	}

	qs = append(qs,
		mc(teacherMinh, math, g5, med, "1/2 + 1/4 bằng?", []string{"1/6", "2/6", "3/4", "1/4"}, 2, "1/2+1/4=3/4", now),
		mc(teacherMinh, math, g5, hard, "Diện tích hình vuông cạnh 12cm?", []string{"24cm²", "48cm²", "144cm²", "36cm²"}, 2, "12×12=144", now),
		mc(teacherMinh, math, g5, easy, "0,5 viết thành phân số là?", []string{"1/5", "1/2", "5/1", "2/5"}, 1, "0.5=1/2", now),
		mc(teacherMinh, viet, g5, med, "Biểu cảm là kiểu văn nào?", []string{"Kể chuyện", "Miêu tả", "Nghị luận", "Biểu cảm"}, 3, "văn biểu cảm", now),
		mc(teacherMinh, sci, g5, med, "Trái Đất quay quanh?", []string{"Mặt Trăng", "Mặt Trời", "Sao Hỏa", "Sao Kim"}, 1, "quanh Mặt Trời", now),
		mc(teacherMinh, eng, g5, med, "Past tense of go?", []string{"goed", "went", "gone", "going"}, 1, "went", now),
		mc(teacherLan, math, string(question.Grade1), easy, "2 + 3 = ?", []string{"4", "5", "6", "3"}, 1, "2+3=5", now),
		mc(teacherLan, math, string(question.Grade2), easy, "20 − 8 = ?", []string{"12", "18", "10", "28"}, 0, "20-8=12", now),
		mc(teacherLan, math, string(question.Grade4), med, "125 × 4 = ?", []string{"400", "500", "525", "450"}, 1, "125×4=500", now),
		tf(teacherMinh, math, g5, easy, "Số nguyên tố nhỏ nhất là 1.", false, "số nguyên tố nhỏ nhất là 2", now),
		tf(teacherMinh, sci, g5, med, "Con người hít vào khí oxygen.", true, "đúng", now),
	)

	return qs
}

func questionIDsByType(qs []*question.Question, createdBy, qType string) []string {
	ids := make([]string, 0)
	for _, q := range qs {
		if q.CreatedBy == createdBy && q.Type == qType {
			ids = append(ids, q.ID.Hex())
		}
	}
	return ids
}

func idsBy(qs []*question.Question, createdBy, qType string, n int) []string {
	return take(questionIDsByType(qs, createdBy, qType), 0, n)
}

func take(ids []string, start, n int) []string {
	if start >= len(ids) {
		return nil
	}
	end := start + n
	if end > len(ids) {
		end = len(ids)
	}
	out := make([]string, end-start)
	copy(out, ids[start:end])
	return out
}

func seedSubmissions(
	homeworkID, teacherID string,
	questionIDs []string,
	students []*student.Student,
	qByID map[string]*question.Question,
	submittedAt time.Time,
) []any {
	docs := make([]any, 0, len(students))
	for i, st := range students {
		if st.Status != student.StudentStatusActive {
			continue
		}
		answers := make([]homeworksubmission.StudentAnswer, 0, len(questionIDs))
		for j, qid := range questionIDs {
			q := qByID[qid]
			ans := homeworksubmission.StudentAnswer{QuestionID: qid}
			correct := (i+j)%4 != 0
			switch question.QuestionType(q.Type) {
			case question.QuestionTypeMultipleChoice:
				idx := 0
				if q.CorrectIndex != nil {
					idx = *q.CorrectIndex
				}
				if !correct {
					idx = (idx + 1) % 4
				}
				ans.SelectedIndex = ptrInt(idx)
			case question.QuestionTypeTrueFalse:
				val := false
				if q.CorrectBool != nil {
					val = *q.CorrectBool
				}
				if !correct {
					val = !val
				}
				ans.SelectedBool = ptrBool(val)
			default:
				continue
			}
			answers = append(answers, ans)
		}
		docs = append(docs, &homeworksubmission.HomeworkSubmission{
			ID:             primitive.NewObjectID(),
			HomeworkID:     homeworkID,
			StudentName:    st.Name,
			IsSubmitted:    true,
			StudentAnswers: answers,
			SubmittedAt:    submittedAt.Add(time.Duration(i) * time.Minute),
			TeacherID:      teacherID,
			CreatedAt:      submittedAt,
			UpdatedAt:      submittedAt,
		})
	}
	return docs
}
