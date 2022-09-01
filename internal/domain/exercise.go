package domain

import "time"

type Exercise struct {
	ID          int        `gorm:"primary_key" json:"id"`
	Title       string     `gorm:"size:255" json:"title"`
	Description string     `gorm:"size:255" json:"description"`
	Questions   []Question `gorm:"foreignkey:ID" json:"questions"`
}

type Question struct {
	ID            int       `gorm:"primary_key" json:"id"`
	ExerciseID    int       `gorm:"index" json:"exercise_id"`
	Body          string    `gorm:"size:255" json:"body"`
	OptionA       string    `gorm:"size:255" json:"option_a"`
	OptionB       string    `gorm:"size:255" json:"option_b"`
	OptionC       string    `gorm:"size:255" json:"option_c"`
	OptionD       string    `gorm:"size:255" json:"option_d"`
	CorrectAnswer string    `gorm:"size:255" json:"correct_answer"`
	Score         int       `gorm:"size:255" json:"score"`
	CreatorID     int       `gorm:"index" json:"creator_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Answer struct {
	ID         int       `gorm:"primary_key" json:"id"`
	ExerciseID int       `gorm:"index" json:"exercise_id"`
	QuestionID int       `gorm:"index" json:"question_id"`
	Answer     string    `gorm:"size:255" json:"answer"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
