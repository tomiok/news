package collector

// GetNewsByUID give a UID (stored in DB) return an *Article.
func (s *Service) GetNewsByUID(uid string) (*Article, error) {
	article, err := s.Storage.getArticleByUID(uid)

	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetFeed will return ALWAYS a slice of articles. A location will be given, if is empty, a default one will be added.
func (s *Service) GetFeed(location string) []Article {
	return nil
}
