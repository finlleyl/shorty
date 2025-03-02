package app

type Store interface {
	Save(shortURL string, originalURL string, userID string) (int, error)
	Get(id string) (string, bool)
	GetAll() []ShortResult
	GetByUserID(userID string) ([]ShortResult, error)
}
