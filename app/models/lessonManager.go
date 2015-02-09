package models

type LessonManager struct {
	lessons []Lesson
}

func NewLessonManager() *LessonManager {
	return &LessonManager{}
}

func (l *LessonManager) GenerateLessons() []Lesson {
	if len(l.lessons) == 0 {
		l.lessons = []Lesson{
			Lesson{1, "Lesson #1 - Basic CRUD operations", "btn-primary"},
			Lesson{2, "Lesson #2 - Advanced Searching API", "btn-success"},
			Lesson{3, "Lesson #3 - Aggreggation and Highlighting", "btn-warning"},
		}
	}
	return l.lessons
}
