package models

const (
	index = "dba"
	_type = "question"
)

type Result struct {
	Ok      bool
	Message string
}

type QuizManager struct {
	em *ELKManager
}

func NewQuizManager() *QuizManager {
	rm := NewResourceManager()
	em := NewELKManager(rm)
	return &QuizManager{em: em}
}

func (q *QuizManager) Dispose() {
	q.em.Dispose()
	q.em = nil
}

func (q QuizManager) Validate(answer int) Result {
	result, err := q.em.LiteralSearchELK(index, _type)
	if err != nil {
		return Result{
			Ok:      false,
			Message: err.Error(),
		}
	} else if result != answer {
		return Result{
			Ok:      false,
			Message: "Sorry, your answer wasn't quite correct. Please try again.",
		}
	} else {
		return Result{
			Ok:      true,
			Message: "That's the correct answer. Well done.",
		}
	}
}
