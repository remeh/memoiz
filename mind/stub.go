package mind

import "remy.io/scratche/uuid"

type Stub struct {
}

func (s *Stub) TryCache(text string) (bool, error) {
	return false, nil
}

func (s *Stub) Fetch(text string) error {
	return nil
}

func (s *Stub) Analyze() error {
	return nil
}

func (s *Stub) Store(uuid.UUID) error {
	return nil
}

func (s *Stub) Categories() Categories {
	return Categories{Unknown}
}
