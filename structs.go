package structs

type Activities struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Sid   string `json:"sid"`
	Score int    `json:"score"`
}

type Players struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type PlayerActivities struct {
	Activity_id int `json:"activity_id"`
	PlayerId    int `json:"player_id"`
	Score       int `json:"score"`
	Time        int `json:"times"`
}
