package domain

type Permissions struct {
	SuperAdmin string `json:"superAdmin"`
	Admin      string `json:"admin"`
	Authors    string `json:"authors"`
	Lyrics     string `json:"lyrics"`
	Users      string `json:"users"`
}
