package mind

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

func (s *Stub) Store() error {
	return nil
}

func (s *Stub) Categories() Categories {
	return Categories{Unknown}
}
