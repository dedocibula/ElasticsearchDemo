package models

var uiHelper *UIHelper

type UIHelper struct {
	lessons       []Lesson
	attemptFields []AttemptField
}

func UIHelperInstance() *UIHelper {
	if uiHelper == nil {
		uiHelper = &UIHelper{}
	}
	return uiHelper
}

func (u *UIHelper) GenerateLessons() []Lesson {
	if len(u.lessons) == 0 {
		u.lessons = []Lesson{
			Lesson{1, "Lesson #1 - Basic CRUD operations", "btn-primary"},
			Lesson{2, "Lesson #2 - Advanced Searching API", "btn-success"},
			Lesson{3, "Lesson #3 - Aggreggation and Highlighting", "btn-warning"},
		}
	}
	return u.lessons
}

func (u *UIHelper) GenerateAttemptFields() []AttemptField {
	if len(u.attemptFields) == 0 {
		u.attemptFields = []AttemptField{
			AttemptField{"failures", "Failed attempts"},
			AttemptField{"successes", "Successful attempts"},
		}
	}
	return u.attemptFields
}
