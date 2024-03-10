package auth

import "github.com/daidr/doulog-core/lib/conf"

var platMap = map[string]struct{}{
	conf.AuthGitHub: {},
}

func PlatExists(platform string) bool {
	_, ok := platMap[platform]
	return ok
}
