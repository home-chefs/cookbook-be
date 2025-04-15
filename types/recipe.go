package types

import (
	"kalos-cookbook/util"
	"time"
)

type Recipe struct {
	ID         int           `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	Name       string        `json:"name"`
	TimeToCook util.Duration `json:"time_to_cook"`
	Source     string        `json:"source"`
	VideoLink  string        `json:"video_link"`
	Products   []Product     `json:"products"`
	Directions string        `json:"directions"`
	Labels     []string      `json:"labels"`
}
