package service

import (
	"time"
	"url-shortnener/pkg/generator"
)

type Store interface {
	Save(code, url string, expiresAt time.Time) error
	Get(code string) (string, time.Time, bool, error)
}

type URLService struct {
	store Store
}

func NewURLService(store Store) *URLService {
	return &URLService{store: store}
}

func (s *URLService) CreateShortURL(
	longURL string,
	ttlMinutes *int,
) (string, error) {

	expiresAt := calculateExpiry(ttlMinutes)

	for {
		code := generator.Generate(6)

		_, _, exists, err := s.store.Get(code)
		if err != nil {
			return "", err
		}

		if !exists {
			err := s.store.Save(code, longURL, expiresAt)
			if err != nil {
				return "", err
			}
			return code, nil
		}
	}
}

func (s *URLService) Resolve(code string) (string, time.Time, bool, error) {
	return s.store.Get(code)
}

func calculateExpiry(ttlMinutes *int) time.Time {
	if ttlMinutes == nil {
		return time.Now().UTC().Add(24 * time.Hour)
	}
	return time.Now().UTC().Add(time.Duration(*ttlMinutes) * time.Minute)
}
