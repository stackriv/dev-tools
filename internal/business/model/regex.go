package model

type RegexMatch struct {
	Match  string   `json:"match"`
	Groups []string `json:"groups"`
	Start  int      `json:"start"`
	End    int      `json:"end"`
}
