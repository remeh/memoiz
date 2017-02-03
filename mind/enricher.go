package mind

type EnrichResult struct {
	Content          string
	ImageUrl         string
	ContentCopyright string
	ImageCopyright   string
}

type EnrichResults []EnrichResult

// Enrichers are the engine allowing to
// add many information to a memo in order
// to send emails with a lot of information.
type Enricher interface {
	// Fetch the data needed to enrich the
	// given text using the given Category,
	// then analyzes the fetched data in order
	// to return a small description and an image
	// Url.
	Enrich(string, Category) (bool, EnrichResult, error)
}

func Enrich(text string, cat Category) (bool, EnrichResult, error) {

	es := make([]Enricher, 0)

	switch cat {
	// TODO(remy): imdb, allociné
	// TODO(remy): yelp
	// TODO(remy): bandcamp ?
	case Artist, Actor, Movie, Person, Place, Serie, VideoGame, Food:
		es = append(es, &Wikipedia{})
	}

	// TODO(remy): long text can probably be sent as is in an email.

	for _, e := range es {
		if found, result, err := e.Enrich(text, cat); err != nil {
			return false, EnrichResult{}, err
		} else if !found {
			continue
		} else {
			return found, result, nil
		}
	}

	return false, EnrichResult{}, nil
}
