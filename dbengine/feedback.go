package dbengine

import "log"

type Feedback struct {
	Comment string `json:"comment"`
	Boxes   string `json:"boxes"`
}

func GetFeedbacks() (feedbacks []Feedback, err error) {
	rows, err := DbConnection.Query(getFeedbacksQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		row := Feedback{}
		err := rows.Scan(&row.Comment, &row.Boxes)
		if err != nil {
			continue
		}

		log.Println(row)
		feedbacks = append(feedbacks, row)
	}

	return feedbacks, rows.Err()
}

func AddFeedback(feedback *Feedback) error {
	_, err := DbConnection.Exec(addFeedbackQuery, feedback.Comment, feedback.Boxes)
	return err
}
