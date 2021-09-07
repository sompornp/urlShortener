package utils

import (
	"github.com/sompornp/urlShortener/model"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestContainBlacklist(t *testing.T) {
	blacklists := []model.Blacklist{
		model.Blacklist{
			IsRegex: false,
			Url:     "http://www.abc.com",
		},
		model.Blacklist{
			IsRegex: false,
			Url:     "http://www",
		},
		model.Blacklist{
			IsRegex: true,
			Url:     "http://www.news.com/*",
		},
	}
	require.True(t, ContainBlacklist(blacklists, "http://www.abc.com"))
	require.True(t, ContainBlacklist(blacklists, "http://www"))
	require.True(t, ContainBlacklist(blacklists, "http://www.news.com/abc/def"))

	require.False(t, ContainBlacklist(blacklists, "http://www.abc"))
	require.False(t, ContainBlacklist(blacklists, "http://ww"))
	require.False(t, ContainBlacklist(blacklists, "http://www.news."))
}

func TestBuildBlacklistUrl(t *testing.T) {
	blacklists := BuildBlacklistUrl("regex: http://news.bbc.co.uk*, http://www.ku.ac.th, http://wwww")
	require.Equal(t, 3,
		len(blacklists))

	require.True(t, blacklists[0].IsRegex)
	require.False(t, blacklists[1].IsRegex)
	require.False(t, blacklists[2].IsRegex)

	require.Equal(t, "http://news.bbc.co.uk*", blacklists[0].Url)
	require.Equal(t, "http://www.ku.ac.th", blacklists[1].Url)
	require.Equal(t, "http://wwww", blacklists[2].Url)
}
