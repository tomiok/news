package collector

func (s *Service) GetNewsByUID(uid string) (*Article, error) {
	article, err := s.Storage.getArticleByUID(uid)

	if err != nil {
		return nil, err
	}

	return article, nil
}
