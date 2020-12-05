package quizzes

import (
	"database/sql"
)

//SQLQuizStore points to a database containing quiz info
type SQLQuizStore struct {
	DB *sql.DB
	//access mongo
}

//GetQuizByID returns a quiz struct by ID
func (store *SQLQuizStore) GetQuizByID(topicID int64) (*Quiz, error) {

}

//GetQuestionByID returns a question struct
func (store *SQLQuizStore) GetQuestionByID(questionID int64) (*Question, error) {
	q := "select questionid, questionbody from questions where questionid = ?"
	row, err := store.DB.QueryRow(q, questionID)
	if err != nil {
		return nil, err
	}
	question := &Question{}
	if err = row.Scan(question.QuestionID, question.QuestionBody); err != nil {
		return nil, err
	}

	return question, nil
}

//ClearQuiz deletes a quiz for a topic
func (store *SQLQuizStore) ClearQuiz(topicID int64) error {

}

//UpdateQuiz updates or populates a topic's quiz
func (store *SQLQuizStore) UpdateQuiz(topicID int64, edit *QuizEdit) (*Quiz, error) {

}

//AnswerQuestion creates a user answer to a question
func (store *SQLQuizStore) AnswerQuestion(questionID int64, userID int64, answer bool) {

}
