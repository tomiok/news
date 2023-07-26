package feed

const (
	Argentina = "Argentina"
	CABA      = "CABA"
)

// GetNewsByUID give a UID (stored in DB) return an *Article.
func (s *Service) GetNewsByUID(uid string) (*Article, error) {
	article, err := s.Storage.getArticleByUID(uid)

	if err != nil {
		return nil, err
	}

	return article, nil
}

// GetFeed will return a slice of articles. A pair of locations will be given, if is empty, a default one will be added.
func (s *Service) GetFeed(locations ...string) ([]Article, error) {
	feed, err := s.Storage.GetDBFeed(locations...)

	if err != nil {
		return nil, err
	}

	return feed, nil
}
