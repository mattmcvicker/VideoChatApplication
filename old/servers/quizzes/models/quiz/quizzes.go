package quizzes

//Quiz structs contain an array of questions as well as the topic the quiz is associated with
type Quiz struct {
	TopicID   int64       `json:"topicid"`
	Questions []*Question `json:"questions"`
}

//Question is a struct which contains question data
type Question struct {
	QuestionID   int64  `json:"questionid"`
	QuestionBody string `json:"questionbody"`
}

//QuizEdit structs contain arrays of questions to add and remove from a quiz.
type QuizEdit struct {
	QuestionsToAdd    []string `json:"questionstoadd"`
	QuestionsToRemove []int64  `json:"questionstoremove"`
}

//UserAnswers organizes user answers into a neat struct, easy to process
type UserAnswers struct {
	UserID  int64     `json:"userid"`
	TopicID int64     `json:"topicid"`
	Answers []*Answer `json:"answers"`
}

//Answer contains data about specific answers
type Answer struct {
	QuestionID int64 `json:"questionid"`
	Response   bool  `json:"response"`
}
