package quizzes

//Store accesses Quizzes
type Store interface {
	GetQuizByID(topicID int64) (*Quiz, error)
	GetQuestionByID(questionID int64) (*Question, error)
	ClearQuiz(topicID int64) error
	UpdateQuiz(topicID int64, edit *QuizEdit) (*Quiz, error)
	//UpdateQuiz can be used to populate an empty quiz
	AnswerQuestion(questionID int64, userID int64, answer bool)
}
