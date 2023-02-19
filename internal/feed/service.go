package feed

var __locs = []string{"Rosario",
	"Argentina",
	"Entre Rios",
	"Cordoba",
	"MDQ",
	"Rio Negro",
	"Corrientes",
	"San Luis",
	"San Juan",
	"Comodoro",
	"Salta",
	"CABA",
	"La Plata",
	"Jujuy",
	"Quilmes",
	"Bahia Blanca",
	"Catamarca",
}

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
func (s *Service) GetFeed(locations ...string) ([]Article, []string, error) {
	if locations == nil || len(locations) == 0 || len(locations) != 2 {
		locations = []string{Argentina, CABA}
	}

	feed, err := s.Storage.GetDBFeed(locations...)

	if err != nil {
		return nil, nil, err
	}

	return feed, locations, nil
}
