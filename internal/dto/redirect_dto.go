package dto

type Redirect struct {
	Code    string `json:"code"`
	IP      string `json:"ip"`
	Country string `json:"country"`
	Device  string `json:"device"`
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Referer string `json:"referer"`
}
