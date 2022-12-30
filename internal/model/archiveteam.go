package model

type AvailableProjects struct {
	Default  string      `json:"auto_project"`
	Projects []ATProject `json:"projects"`
}

type Project struct {
}

type ATProjectResponse struct {
	AutoProject              string      `json:"auto_project"`
	DISABLEDBroadcastMessage string      `json:"DISABLED_broadcast_message"`
	TrackerBannerHTML        string      `json:"tracker_banner_html"`
	Projects                 []ATProject `json:"projects"`
	XXXprojects              []ATProject `json:"XXXprojects"`
	Warrior                  ATWarrior   `json:"warrior"`
}

type ATProject struct {
	Description string `json:"description"`
	Leaderboard string `json:"leaderboard"`
	Logo        string `json:"logo"`
	MarkerHTML  string `json:"marker_html"`
	Name        string `json:"name"`
	Repository  string `json:"repository"`
	Title       string `json:"title"`
}

type ATWarrior struct {
	SeesawVersion string `json:"seesaw_version"`
}
