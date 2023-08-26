package storage

type Storage struct {
	usersStorage map[string][]string
}

func (s *Storage) AddText(user string, text string) {
	s.usersStorage[user] = append(s.usersStorage[user], text)
}

func (s *Storage) GetTexts(user string) []string {
	return s.usersStorage[user]
}
