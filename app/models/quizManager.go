package models

type Result struct {
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

func (q QuizManager) Validate(answer int) {

}
