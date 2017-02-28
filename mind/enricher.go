package mind

type EnrichResult struct {
	Title            string
	Content          string
	ContentCopyright string
	ContentSource    string
	ImageUrl         string
	ImageCopyright   string
	ImageSource      string
	Format           EnrichFormat
}

func (e EnrichResult) Enriched() bool {
	if len(e.Content) > 0 || len(e.ImageUrl) > 0 {
		return true
	}
	return false
}

type EnrichFormat string

var (
	EnrichStandard      EnrichFormat = "standard"
	EnrichUrlNoImage    EnrichFormat = "url_no_image"
	EnrichUrlFocusImage EnrichFormat = "url_focus_image"
)

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

// Enrich builds an EnrichResult containing description and image
// to complete the information guessed from the given text.
// The current purpose is to use these information to send contextual
// and interesting emails.
func Enrich(text string, cat Category) (bool, EnrichResult, error) {

	es := make([]Enricher, 0)

	// looks whether the text contains an Url
	// ---------------------

	url := rxUrl.FindString(text)
	if len(url) != 0 {
		// we have an URL: starts with the Url enricher
		// because it'll probably be enough.
		es = append(es, &Url{url: url})
	}

	// other engines
	// ----------------------

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
