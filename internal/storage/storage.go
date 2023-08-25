package storage

type Storage struct {
	usersStorage map[string][]string
}

func Create() *Storage {
	return &Storage{make(map[string][]string)}
}

func (s *Storage) AddText(user string, text string) {
	s.usersStorage[user] = append(s.usersStorage[user], text)
}
