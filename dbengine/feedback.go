package dbengine

import "errors"

type Feedback struct {
	Comment string `json:"comment"`
	Boxes   string `json:"boxes"`
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
		row := Feedback{}
		err := rows.Scan(&row.Comment, &row.Boxes)
		if err != nil {
			return feedbacks, err
		}

		feedbacks = append(feedbacks, row)
	}

	return feedbacks, rows.Err()
}

func AddFeedback(feedback *Feedback) error {
	_, err := DbConnection.Exec(addFeedbackQuery, feedback.Comment, feedback.Boxes)
	return err
}
