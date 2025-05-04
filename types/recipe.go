package types

import "time"

type Recipe struct {
	ID             int       `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	Name           string    `json:"name"`
	TimeToCook     Duration  `json:"time_to_cook"`
	Source         string    `json:"source"`
	VideoLink      string    `json:"video_link"`
	Products       []Product `json:"products"`
	Directions     string    `json:"directions"`
	Labels         []string  `json:"labels"`
	CoverImagePath *string   `json:"cover_image_path"`
}
