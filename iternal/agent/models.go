package agent

type Res struct {
	ID      string  `json:"id"`
	Result  float64 `json:"result"`
	Timeout bool    `json:"timeout"`
	Errors  string  `json:"errors"`
}
