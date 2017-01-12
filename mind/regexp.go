package mind

import "regexp"

// rxUrl looks for an URL in an arbitrary text
var rxUrl = regexp.MustCompile(`((https?:\/\/)?([0-9a-zA-Z]+\.)*[-_0-9a-zA-Z]+\.[-_0-9a-zA-Z]+)\/([-_0-9a-zA-Z\.\/\-\=\_\.\!\+\,@])*(\?[0-9a-zA-Z\%\&\-\=\_\.\!\+\,@]*)*`)

// rxDomain retrieves the domain (removing the TLD) of the URL (only
// if there is a trailing '/')
var rxDomain *regexp.Regexp = regexp.MustCompile(`([a-zA-Z0-9]*)\.[a-zA-Z0-9]*\/`)
