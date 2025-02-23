package app

type Store interface {
	Save(shortURL string, originalURL string) int
	Get(id string) (string, bool)
	GetAll() []ShortResult
}
