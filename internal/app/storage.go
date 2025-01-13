package app

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Save(url string) string {
	id := generateID()
	s.data[id] = url
	return id
}

func (s *Storage) Get(id string) (string, bool) {
	url, exists := s.data[id]
	return url, exists
}
