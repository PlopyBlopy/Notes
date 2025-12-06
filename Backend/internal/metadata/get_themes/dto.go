package getthemes

type output struct {
	Themes []theme `json:"themes"`
}

type theme struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}
