package storage

type Manager interface {
	AddText(user string, text string)
	GetTexts(user string) []string
}

func CreateRamStorage() Manager {
	return &Storage{make(map[string][]string)}
}
