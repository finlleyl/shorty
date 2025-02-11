package app

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

type ShortResult struct {
	ID          int    `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

type Storage struct {
	mu   sync.Mutex
	data []ShortResult
	file string
}

func NewStorage(file string) *Storage {
	s := &Storage{file: file}
	s.Load()
	return s
}

func (s *Storage) Save(shortUrl, originalUrl string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID := 1
	if len(s.data) > 0 {
		newID = s.data[len(s.data)-1].ID + 1
	}

	s.data = append(s.data, ShortResult{
		ID:          newID,
		ShortUrl:    shortUrl,
		OriginalUrl: originalUrl,
	})

	s.Flush()
	return newID
}

func (s *Storage) GetAll() []ShortResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.data
}

func (s *Storage) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, data, 0644)
}

func (s *Storage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			s.data = []ShortResult{}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &s.data)
}

func (s *Storage) Get(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range s.data {
		if v.ID, _ = strconv.Atoi(id); v.ShortUrl == id {
			return v.OriginalUrl, true
		}
	}
	return "", false
}
