package storage

type Storage struct {
	usersStorage map[string][]string
}

func (s *Storage) AddText(user string, text string) error {
	s.usersStorage[user] = append(s.usersStorage[user], text)
	return nil
}

func (s *Storage) GetTexts(user string) ([]string, error) {
	return s.usersStorage[user], nil
}
