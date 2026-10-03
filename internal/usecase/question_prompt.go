package usecase

import (
	"fmt"
	"teaching_assistant/internal/domain/question"
)

const generateQuestionSystemPrompt = `Bạn là giáo viên tiểu học Việt Nam (lớp 1–5), soạn câu hỏi ôn tập theo chương trình phổ thông.

Trả về DUY NHẤT một JSON object, không markdown, không giải thích ngoài JSON.

Schema bắt buộc:
{
  "questions": [
    {
      "question": "string — nội dung đề",
      "options": ["string"],
      "correct_index": 0,
      "correct_bool": true,
      "explanation": "string — lời giải ngắn cho giáo viên"
    }
  ]
}

Quy tắc:
- Đúng số lượng câu được yêu cầu.
- Câu rõ ràng, phù hợp lứa tuổi, không nội dung nhạy cảm.
- Trắc nghiệm (multiple_choice): đúng 4 lựa chọn khác nhau; correct_index là số nguyên 0–3; không dùng correct_bool.
- Đúng/sai (true_false): không gửi options; correct_bool là true hoặc false; không dùng correct_index.
- Môn english: đề và đáp án bằng tiếng Anh đơn giản, explanation bằng tiếng Việt.
- Các môn còn lại: toàn bộ bằng tiếng Việt, chính tả đúng.
- Không lặp ý giữa các câu.`

func generateQuestionUserPrompt(questionType, grade, subject, difficulty string, quantity int) string {
	typeLabel := "trắc nghiệm (multiple_choice)"
	typeRule := "Mỗi câu có đúng 4 options. correct_index trỏ tới đáp án đúng (0-based)."
	if question.QuestionType(questionType) == question.QuestionTypeTrueFalse {
		typeLabel = "đúng/sai (true_false)"
		typeRule = "Không có options. correct_bool là true hoặc false."
	}

	return fmt.Sprintf(`Soạn %d câu hỏi %s.

Tham số (gắn vào từng câu, không đổi):
- type: %s
- grade: %s (lớp %s)
- subject: %s (%s)
- difficulty: %s (%s)

%s

Chỉ trả JSON {"questions":[...]}.`,
		quantity,
		typeLabel,
		questionType,
		grade, grade,
		subject, subjectLabel(subject),
		difficulty, difficultyLabel(difficulty),
		typeRule,
	)
}

func subjectLabel(subject string) string {
	switch question.Subject(subject) {
	case question.SubjectVietnamese:
		return "Tiếng Việt"
	case question.SubjectMathematics:
		return "Toán"
	case question.SubjectEthics:
		return "Đạo đức"
	case question.SubjectEnglish:
		return "Tiếng Anh"
	case question.SubjectNatureAndSociety:
		return "Tự nhiên và xã hội"
	case question.SubjectHistoryAndGeography:
		return "Lịch sử và địa lý"
	case question.SubjectScience:
		return "Khoa học"
	case question.SubjectInformatics:
		return "Tin học"
	case question.SubjectTechnology:
		return "Công nghệ"
	case question.SubjectPhysicalEducation:
		return "Thể dục"
	case question.SubjectMusic:
		return "Âm nhạc"
	case question.SubjectArt:
		return "Mỹ thuật"
	case question.SubjectExperientialActivities:
		return "Hoạt động trải nghiệm"
	default:
		return subject
	}
}

func difficultyLabel(difficulty string) string {
	switch question.Difficulty(difficulty) {
	case question.DifficultyEasy:
		return "dễ"
	case question.DifficultyMedium:
		return "vừa"
	case question.DifficultyHard:
		return "khó"
	default:
		return difficulty
	}
}
