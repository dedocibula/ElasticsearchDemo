package models

const (
	index       = "dba"
	quizType    = "question"
	rankingType = "ranking"
)

var (
	atomicCounter = 0
)

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
	result, err := q.em.LiteralSearchELK(index, quizType)
	if err != nil {
		return q.buildResult(false, err.Error())
	} else if result != answer {
		return q.buildResult(false,
			"Sorry, your answer wasn't quite correct. Please try again.")
	} else {
		return q.buildResult(true, "That's the correct answer. Well done.")
	}
}

func (q QuizManager) buildResult(correct bool, message string) Result {
	return Result{
		Ok:      correct,
		Message: message,
	}
}
