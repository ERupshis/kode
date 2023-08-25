package jsonmsg

//go:generate easyjson -all jsonmsg.go
type Input struct {
	Text string `json:"data"`
}

type Output struct {
	Texts []string `json:"dataArray"`
}
