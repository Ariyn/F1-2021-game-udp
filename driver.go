package F1_2021_game_udp

type Driver struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	RaceNumber int    `json:"raceNumber"`
	TeamName   string `json:"teamName"`
	CarIndex   int    `json:"carIndex"`
	IsAi       bool   `json:"isAi"`
}
