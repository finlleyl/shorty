package app

import (
	"encoding/json"
	"os"
)

type ShortResult struct {
	ID          int    `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	data []ShortResult
	file string
}

func NewStorage(file string) *Storage {
	s := &Storage{file: file}
	s.Load()
	return s
}

func (s *Storage) Save(shortURL, originalURL string) int {
	newID := 1
	if len(s.data) > 0 {
		newID = s.data[len(s.data)-1].ID + 1
	}

	s.data = append(s.data, ShortResult{
		ID:          newID,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	})

	s.Flush()
	return newID
}

func (s *Storage) GetAll() []ShortResult {
	return s.data
}

func (s *Storage) Flush() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}

	return os.WriteFile(s.file, data, 0644)
}

func (s *Storage) Load() error {
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
	for _, v := range s.data {
		if v.ShortURL == id {
			return v.OriginalURL, true
		}
	}
	return "", false
}
