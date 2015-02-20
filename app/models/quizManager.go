package models

import "time"

const (
	index       = "dba"
	quizType    = "question"
	rankingType = "ranking"
)

type QuizManager struct {
	em          *ELKManager
	validations []func(Attempt) (bool, string)
}

func NewQuizManager() *QuizManager {
	rm := NewResourceManager()
	em := NewELKManager(rm)
	return &QuizManager{em: em}
}

func (q *QuizManager) Dispose() {
	q.em.Dispose()
	q.em = nil
	q.validations = nil
}

func (q *QuizManager) Validate(attempt Attempt) Result {
	QuizMonitorInstance().Publish(attempt)

	q.setupValidations()
	r := q.validateAsync(attempt)
	if result := <-r; result.Ok {
		record := q.createELKRecord(attempt.Name)
		return q.submitELKRecord(record)
	} else {
		return result
	}
}

func (q QuizManager) GetResults() []ELKRecord {
	records, err := q.em.SelectRecordsELK(index, rankingType)
	if err != nil {
		return make([]ELKRecord, 0)
	} else {
		return records
	}
}

func (q QuizManager) ClearResults() bool {
	err := q.em.ClearTypeELK(index, rankingType)
	return err == nil
}

func (q QuizManager) validateAsync(attempt Attempt) chan Result {
	in := make(chan Result, len(q.validations))
	out := make(chan Result, 1)

	for _, validationFunc := range q.validations {
		go func(f func(Attempt) (bool, string)) {
			ok, message := f(attempt)
			in <- q.buildResult(ok, message)
		}(validationFunc)
	}

	go func() {
		result := q.buildResult(true)
		for i := 0; i < len(q.validations); i++ {
			result = q.mergeResults(result, <-in)
		}
		out <- result
		close(in)
		close(out)
	}()

	return out
}

func (q *QuizManager) setupValidations() {
	if q.validations == nil {
		q.validations = []func(Attempt) (bool, string){
			q.validateAnswer,
			q.validateName,
		}
	}
}

func (q QuizManager) validateAnswer(attempt Attempt) (bool, string) {
	result, err := q.em.LiteralSearchELK(index, quizType)
	if err != nil {
		return false, err.Error()
	} else if attempt.Answer != result {
		return false, "Sorry, your answer wasn't quite correct. Please try again."
	} else {
		return true, ""
	}
}

func (q QuizManager) validateName(attempt Attempt) (bool, string) {
	exists, err := q.em.ExistsRecordELK(index, rankingType, attempt.Name)
	if err != nil {
		return false, err.Error()
	} else if exists {
		return false, "Sorry, this nickname appears to be taken."
	} else {
		return true, ""
	}
}

func (q QuizManager) mergeResults(r1, r2 Result) Result {
	return Result{
		Ok:       r1.Ok && r2.Ok,
		Messages: append(r1.Messages, r2.Messages...),
	}
}

func (q QuizManager) buildResult(correct bool, messages ...string) Result {
	return Result{
		Ok:       correct,
		Messages: messages,
	}
}

func (q QuizManager) createELKRecord(nickname string) ELKRecord {
	return ELKRecord{
		Nickname:  nickname,
		Timestamp: time.Now(),
	}
}

func (q QuizManager) submitELKRecord(record ELKRecord) Result {
	err := q.em.RecordSuccessELK(index, rankingType, record)
	if err != nil {
		return q.buildResult(false, "Sorry, this nickname appears to be taken.")
	} else {
		return q.buildResult(true, "That's the correct answer. Well done.")
	}
}
