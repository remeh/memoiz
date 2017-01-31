package mind

type Description string
type ImageUrl string

// Enrichers are the engine allowing to
// add many information to a memo in order
// to send emails with a lot of information.
type Enricher interface {
	// Fetch the data needed to enrich the
	// given text using the given Category,
	// then analyzes the fetched data in order
	// to return a small description and an image
	// Url.
	Enrich(string, Category) (bool, Description, ImageUrl, error)
}

func Enrich(text string, cat Category) (bool, Description, ImageUrl, error) {

	es := make([]Enricher, 0)

	switch cat {
	case Artist, Actor, Movie, Person, Place, Serie, VideoGame, Food:
		es = append(es, &Wikipedia{})
	}

	for _, e := range es {
		if found, desc, imgUrl, err := e.Enrich(text, cat); err != nil {
			return false, "", "", err
		} else if !found {
			continue
		} else {
			return found, desc, imgUrl, nil
		}
	}

	return false, "", "", nil
}
