package domain

import (
	"github.com/pkg/errors"
)

var ErrQuestionsIndexOutOfRange = errors.New("Index out of range")
var ErrQuestionsNoNext = errors.New("No next question")
var ErrQuestionsIdExists = errors.New("Duplicate question id")
var ErrQuestionsEmpty = errors.New("Cannot do this operation on an empty Questions struct")

// stores a list of Question objects
// not a map because we need to keep the order of the questions
// maps in Go do not keep order
type Questions struct {
	text []string
	ids  []string
}

func NewQuestions() *Questions {
	return &Questions{}
}

func (qs *Questions) Add(id, question string) error {
	for _, existingId := range qs.ids {
		if existingId == id {
			return ErrQuestionsIdExists
		}
	}
	qs.text = append(qs.text, question)
	qs.ids = append(qs.ids, id)
	return nil
}

func (qs *Questions) At(index int) (q *Question, err error) {
	if index >= len(qs.ids) {
		return q, ErrQuestionsIndexOutOfRange
	}
	return &Question{
		Id:   qs.ids[index],
		Text: qs.text[index],
	}, nil
}

func (qs *Questions) Next(qid string) (q *Question, err error) {
	returnNext := false
	for index, id := range qs.ids {
		if returnNext {
			return qs.At(index)
		}
		if id == qid {
			returnNext = true
			continue
		}
	}
	return q, ErrQuestionsNoNext
}

func (qs *Questions) Last() (q *Question, err error) {
	if len(qs.ids) == 0 {
		return q, ErrQuestionsEmpty
	}

	return &Question{
		Id:   qs.ids[len(qs.ids)-1],
		Text: qs.text[len(qs.text)-1],
	}, nil
}
