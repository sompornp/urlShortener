package utils

import (
	"github.com/sompornp/urlShortener/constant"
	"github.com/sompornp/urlShortener/model"
	"regexp"
	"strings"
)

func ContainBlacklist(blacklists []model.Blacklist, url string) bool {
	for _, b := range blacklists {
		if b.IsRegex {
			match, _ := regexp.MatchString(b.Url, url)
			if match {
				return true
			}
		} else if b.Url == url {
			return true
		}
	}
	return false
}

func BuildBlacklistUrl(blacklistUrl string) []model.Blacklist {
	blacklists := []model.Blacklist{}
	b := strings.Split(blacklistUrl, ",")
	for i := range b {
		btext := strings.TrimSpace(b[i])
		if btext != "" {
			blacklist := model.Blacklist{IsRegex: strings.HasPrefix(btext, constant.RegexPrefix), Url: btext}
			if blacklist.IsRegex {
				blacklist.Url = strings.TrimSpace(blacklist.Url[len(constant.RegexPrefix):])
			}
			blacklists = append(blacklists, blacklist)
		}
	}
	return blacklists
}
