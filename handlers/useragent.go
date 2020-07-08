package handlers

import "regexp"

var (
	mobileRegexp *regexp.Regexp
)

func initUAgentRegex() {
	if mobileRegexp == nil {
		mobileRegexp = regexp.MustCompile("(?i)(iPhone|iPod|iPad|Android|BlackBerry|Windows Phone)")
	}
}

// Return true if given
// useragent is a mobile uagent
func isMobileUseragent(uagent string) bool {
	initUAgentRegex()
	return mobileRegexp.MatchString(uagent)
}
