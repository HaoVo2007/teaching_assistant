package mapper

import (
	"math"

	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
)

const homeworkMaxScore = 100.0

func difficultyWeight(difficulty string) int {
	switch question.Difficulty(difficulty) {
	case question.DifficultyHard:
		return 3
	case question.DifficultyMedium:
		return 2
	default:
		return 1
	}
}

func allocateMaxScores(questions []*question.Question) map[string]float64 {
	scores := make(map[string]float64, len(questions))
	if len(questions) == 0 {
		return scores
	}

	weights := make([]int, len(questions))
	totalWeight := 0
	for i, q := range questions {
		if q == nil {
			continue
		}
		w := difficultyWeight(q.Difficulty)
		weights[i] = w
		totalWeight += w
	}
	if totalWeight == 0 {
		return scores
	}

	var allocated float64
	lastIdx := -1
	for i, q := range questions {
		if q == nil || weights[i] == 0 {
			continue
		}
		lastIdx = i
	}

	for i, q := range questions {
		if q == nil || weights[i] == 0 {
			continue
		}
		id := q.ID.Hex()
		if i == lastIdx {
			scores[id] = round2(homeworkMaxScore - allocated)
			continue
		}
		point := round2(homeworkMaxScore * float64(weights[i]) / float64(totalWeight))
		scores[id] = point
		allocated += point
	}
	return scores
}

func isStudentAnswerCorrect(q *question.Question, answer homeworksubmission.StudentAnswer) bool {
	if q == nil {
		return false
	}
	switch question.QuestionType(q.Type) {
	case question.QuestionTypeMultipleChoice:
		return q.CorrectIndex != nil && answer.SelectedIndex != nil && *q.CorrectIndex == *answer.SelectedIndex
	case question.QuestionTypeTrueFalse:
		return q.CorrectBool != nil && answer.SelectedBool != nil && *q.CorrectBool == *answer.SelectedBool
	default:
		return false
	}
}

func MapHomeworkSubmissionToResponse(
	submission *homeworksubmission.HomeworkSubmission,
	hw *homework.Homework,
	questionsByID map[string]*question.Question,
) *response.HomeworkSubmissionResponse {
	orderedQuestions := homeworkQuestions(hw, submission, questionsByID)
	maxScores := allocateMaxScores(orderedQuestions)

	answers := make([]response.StudentAnswer, 0, len(submission.StudentAnswers))
	var total float64
	for _, answer := range submission.StudentAnswers {
		q := questionsByID[answer.QuestionID]
		isCorrect := isStudentAnswerCorrect(q, answer)
		maxScore := maxScores[answer.QuestionID]
		score := 0.0
		if isCorrect {
			score = maxScore
		}
		total += score

		mappedQuestion := response.QuestionResponse{}
		if q != nil {
			mappedQuestion = *MapQuestionToResponse(q)
		}

		answers = append(answers, response.StudentAnswer{
			Question:      mappedQuestion,
			SelectedIndex: answer.SelectedIndex,
			SelectedBool:  answer.SelectedBool,
			IsCorrect:     isCorrect,
			Score:         score,
			MaxScore:      maxScore,
		})
	}

	return &response.HomeworkSubmissionResponse{
		ID:             submission.ID.Hex(),
		HomeworkID:     submission.HomeworkID,
		StudentName:    submission.StudentName,
		IsSubmitted:    submission.IsSubmitted,
		StudentAnswers: answers,
		TotalScore:     round2(total),
		MaxScore:       homeworkMaxScore,
		SubmittedAt:    submission.SubmittedAt,
		CreatedAt:      submission.CreatedAt,
		UpdatedAt:      submission.UpdatedAt,
	}
}

func MapHomeworkSubmissionsToResponses(
	submissions []*homeworksubmission.HomeworkSubmission,
	homeworksByID map[string]*homework.Homework,
	questionsByID map[string]*question.Question,
) []*response.HomeworkSubmissionResponse {
	responses := make([]*response.HomeworkSubmissionResponse, 0, len(submissions))
	for _, submission := range submissions {
		if submission == nil {
			continue
		}
		responses = append(responses, MapHomeworkSubmissionToResponse(
			submission,
			homeworksByID[submission.HomeworkID],
			questionsByID,
		))
	}
	return responses
}

func homeworkQuestions(
	hw *homework.Homework,
	submission *homeworksubmission.HomeworkSubmission,
	questionsByID map[string]*question.Question,
) []*question.Question {
	ids := make([]string, 0)
	if hw != nil && len(hw.Questions) > 0 {
		ids = hw.Questions
	} else {
		for _, answer := range submission.StudentAnswers {
			ids = append(ids, answer.QuestionID)
		}
	}

	items := make([]*question.Question, 0, len(ids))
	for _, id := range ids {
		if q := questionsByID[id]; q != nil {
			items = append(items, q)
		}
	}
	return items
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
