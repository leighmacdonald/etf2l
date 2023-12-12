package etf2l

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Ban struct {
	CurrentPage  int           `json:"current_page"`
	Data         []interface{} `json:"data"`
	FirstPageUrl string        `json:"first_page_url"`
	From         *string       `json:"from"`
	LastPage     int           `json:"last_page"`
	LastPageUrl  string        `json:"last_page_url"`
	Links        []struct {
		Url    *string `json:"url"`
		Label  string  `json:"label"`
		Active bool    `json:"active"`
	} `json:"links"`
	NextPageUrl *string `json:"next_page_url"`
	Path        string  `json:"path"`
	PerPage     int     `json:"per_page"`
	PrevPageUrl *string `json:"prev_page_url"`
	To          *string `json:"to"`
	Total       int     `json:"total"`
}

type BansResponse struct {
	Bans Ban `json:"bans"`
}

func Bans() (Ban, error) {
	return Ban{}, nil
}
