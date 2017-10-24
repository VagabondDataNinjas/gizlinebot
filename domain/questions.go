package domain

import (
	"encoding/json"
	"strconv"

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
	text    []string
	ids     []string
	weights []int
}

func NewQuestions() *Questions {
	return &Questions{}
}

func (qs *Questions) MarshalJSON() ([]byte, error) {
	q := make(map[string]map[string]string)
	for index, id := range qs.ids {
		q[id] = map[string]string{
			"text":   qs.text[index],
			"weight": strconv.Itoa(qs.weights[index]),
		}
	}
	return json.Marshal(q)
}

func (qs *Questions) Add(id, question string, weight int) error {
	for _, existingId := range qs.ids {
		if existingId == id {
			return ErrQuestionsIdExists
		}
	}
	qs.text = append(qs.text, question)
	qs.ids = append(qs.ids, id)
	qs.weights = append(qs.weights, weight)
	return nil
}

func (qs *Questions) At(index int) (q *Question, err error) {
	if index >= len(qs.ids) {
		return q, ErrQuestionsIndexOutOfRange
	}
	return &Question{
		Id:     qs.ids[index],
		Text:   qs.text[index],
		Weight: qs.weights[index],
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
		Id:     qs.ids[len(qs.ids)-1],
		Text:   qs.text[len(qs.text)-1],
		Weight: qs.weights[len(qs.weights)-1],
	}, nil
}
