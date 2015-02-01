package models

type Lesson struct {
	Number      int
	Description string
	HtmlClass   string
}

type LessonManager struct {
	Lessons []Lesson
}

func (l *LessonManager) GenerateLessons() []Lesson {
	if len(l.Lessons) == 0 {
		l.Lessons = []Lesson{
			Lesson{1, "Lesson #1 - Basic CRUD operations", "btn-primary"},
			Lesson{2, "Lesson #2 - Advanced Searching API", "btn-success"},
			Lesson{3, "Lesson #3 - Aggreggation and Highlighting", "btn-warning"},
		}
	}
	return l.Lessons
}
