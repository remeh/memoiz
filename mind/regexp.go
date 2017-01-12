package mind

import "regexp"

// TODO(remy): comment me
var rxUrl = regexp.MustCompile(`((https?:\/\/)?([0-9a-zA-Z]+\.)*[-_0-9a-zA-Z]+\.[-_0-9a-zA-Z]+)\/([-_0-9a-zA-Z\.\/])*(\?[0-9a-zA-Z\%\&\-\=\_\.]*)*`)

// rxDomain retrieves only the domain (removing the TLD)
var rxDomain *regexp.Regexp = regexp.MustCompile(`([a-zA-Z0-9]*)\.[a-zA-Z0-9]*\/`)
