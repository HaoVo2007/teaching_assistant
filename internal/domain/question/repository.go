package question

import "context"

type QuestionRepository interface {
	Create(ctx context.Context, q *Question) error
}
