package dbengine

import (
	"errors"
	"strconv"
	"strings"
)

type Feedback struct {
	Comment string `json:"comment"`
	Boxes   []int  `json:"boxes"`
}

func GetFeedbacks(count uint) (feedbacks []Feedback, err error) {
	if count == 0 {
		return []Feedback{}, errors.New("count has to be at least 1")
	}

	rows, err := DbConnection.Query(getFeedbacksQuery, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var row Feedback
		var boxes string
		err := rows.Scan(&row.Comment, &boxes)
		if err != nil {
			return feedbacks, err
		}

		for _, box := range strings.Split(boxes, ",") {
			parsed, err := strconv.Atoi(box)
			if err == nil {
				row.Boxes = append(row.Boxes, parsed)
			}
		}

		feedbacks = append(feedbacks, row)
	}

	return feedbacks, rows.Err()
}

func AddFeedback(feedback *Feedback) error {
	var boxes []string
	for _, val := range feedback.Boxes {
		boxes = append(boxes, strconv.Itoa(val))
	}
	_, err := DbConnection.Exec(addFeedbackQuery, feedback.Comment, strings.Join(boxes, ", "))
	return err
}
