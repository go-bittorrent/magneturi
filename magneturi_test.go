package magneturi_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-bittorrent/magneturi"
)

const (
	badURLEscape = "q%w%erty"
)

func TestParse(t *testing.T) {
	t.Parallel()

	t.Run("Parse", func(t *testing.T) {
		t.Parallel()

		parsed, err := magneturi.Parse("magnet:?xt=urn:ed2k:354B15E68FB8F36D7CD88FF94116CDC1&xt=urn:tree:tiger:7N5OAMRNGMSSEUE3ORHOKWN4WWIQ5X4EBOOTLJY&xt=urn:btih:QHQXPYWMACKDWKP47RRVIV7VOURXFE5Q&xl=10826029&as=http%3A%2F%2Fdownload.wikimedia.org%2Fmediawiki%2F1.15%2Fmediawiki-1.15.1.tar.gz&dn=mediawiki-1.15.1.tar.gz&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&xs=http%3A%2F%2Fcache.example.org%2FXRX2PEFXOOEJFRVUCX6HMZMKS5TWG4K5&xs=dchub%3A%2F%2Fexample.org&mt=http://weblog.foo/all-my-favorites.rss&kt=tag1+tag2&x.some=qwerty&x.some=qwerty2")
		require.NoError(t, err)

		require.EqualValues(t, parsed.DisplayName, "mediawiki-1.15.1.tar.gz")
		require.ElementsMatch(t, parsed.ExactTopics, []string{
			"urn:tree:tiger:7N5OAMRNGMSSEUE3ORHOKWN4WWIQ5X4EBOOTLJY",
			"urn:ed2k:354B15E68FB8F36D7CD88FF94116CDC1",
			"urn:btih:QHQXPYWMACKDWKP47RRVIV7VOURXFE5Q",
		})
		require.EqualValues(t, parsed.ExactLength, 10826029)
		require.ElementsMatch(t, parsed.AcceptableSources, []string{
			"http://download.wikimedia.org/mediawiki/1.15/mediawiki-1.15.1.tar.gz",
		})
		require.ElementsMatch(t, parsed.ExactSource, []string{
			"http://cache.example.org/XRX2PEFXOOEJFRVUCX6HMZMKS5TWG4K5",
			"dchub://example.org",
		})
		require.ElementsMatch(t, parsed.KeywordTopic, []string{
			"tag1",
			"tag2",
		})
		require.EqualValues(t, parsed.ManifestTopic, "http://weblog.foo/all-my-favorites.rss")
		require.ElementsMatch(t, parsed.Trackers, []string{
			"udp://tracker.openbittorrent.com:80/announce",
		})
		require.Len(t, parsed.AdditionParams, 1)
		t.Log(parsed)
	})
	t.Run("TestParseBadURI", func(t *testing.T) {
		t.Parallel()

		_, err := magneturi.Parse("magnet")
		require.Error(t, err)
	})
	t.Run("TestParseBadParam", func(t *testing.T) {
		t.Parallel()

		params, err := magneturi.Parse("magnet:?xt=some&k=&g")
		require.NoError(t, err)
		require.Len(t, params.AdditionParams, 0)
	})
	t.Run("TestParseBadQueryEscape", func(t *testing.T) {
		t.Parallel()

		mustQueryUnescape := [...]string{"dn", "tr", "as", "xs", "xl", "x.some"}
		for _, tag := range mustQueryUnescape {
			_, err := magneturi.Parse("magnet:?" + tag + "=" + badURLEscape)
			require.Error(t, err)
		}
	})
}

func TestMagnet_Encoded(t *testing.T) {
	t.Parallel()

	// heh, I changed the hashes, you won't download anything :)
	testCases := []string{
		"magnet:?xt=urn:ed2k:354B15E68FB8F3347CD88FF94116CDC1&xt=urn:tree:tiger:7N5OAMRNKLSSEUE3ORHOKWN4WWIQ5X4EBOOTLJY&xt=urn:btih:QHQXPYWOPCKDWKP47RRVIV7VOURXFE5Q&xl=10826029&as=http%3A%2F%2Fdownload.wikimedia.org%2Fmediawiki%2F1.15%2Fmediawiki-1.15.1.tar.gz&dn=mediawiki-1.15.1.tar.gz&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&xs=http%3A%2F%2Fcache.example.org%2FXRX2PEFXOOEJFRVUCX6HMZMKS5TWG4K5&xs=dchub%3A%2F%2Fexample.org&mt=http://weblog.foo/all-my-favorites.rss&kt=tag1+tag2&x.some=qwerty&x.some=qwerty2",
		"magnet:?xt=urn:btih:a856772133117dd29d6e79695f4d087fa3dd3214&dn=Depeche%20Mode%20%e2%80%93%20Playing%20The%20Angel%20(Deluxe)%20(2005)&tr=http%3a%2f%2fbt3.t-ru.org%2fann&tr=http%3a%2f%2fretracker.local%2fannounce",
		"magnet:?xt=urn:btih:45OOC4380ACE0E39D1989C58D3CFE71DBE974E2D&tr=http%3A%2F%2Fbt3.t-ru.org%2Fann%3Fmagnet&dn=%D0%92%D0%B5%D1%81%D1%8C%20%D0%A7%D0%B5%D0%BC%D0%BF%D0%B8%D0%BE%D0%BD%D0%B0%D1%82%20%D0%95%D0%B2%D1%80%D0%BE%D0%BF%D1%8B%202020%20%2F%20UEFA%20Euro%202020%20%2F%20Live%20%2F%20%D0%9C%D0%B0%D1%82%D1%87%20%D0%A2%D0%92%2C%20%D0%9F%D0%B5%D1%80%D0%B2%D1%8B%D0%B9%20%D0%BA%D0%B0%D0%BD%D0%B0%D0%BB%2C%20%D0%A0%D0%BE%D1%81%D1%81%D0%B8%D1%8F%201%20%5B11.06-11.07.2021%2C%20%D0%A4%D1%83%D1%82%D0%B1%D0%BE%D0%BB%2C%20HDTVRip%2F720p%2F50fps%2C%20MKV%2FH.264%2C%20RU%5D",
		"magnet:?xt=urn:btih:40BDD8101119EF795C9B31B456701C2513819FCE&tr=http%3A%2F%2Fbt4.t-ru.org%2Fann%3Fmagnet&dn=%D0%91%D1%8D%D1%82%D0%BC%D0%B5%D0%BD%20%2F%20The%20Batman%20(%D0%9C%D1%8D%D1%82%D1%82%20%D0%A0%D0%B8%D0%B2%D0%B7%20%2F%20Matt%20Reeves)%20%5B2022%2C%20%D0%A1%D0%A8%D0%90%2C%20%D0%B1%D0%BE%D0%B5%D0%B2%D0%B8%D0%BA%2C%20%D0%BA%D1%80%D0%B8%D0%BC%D0%B8%D0%BD%D0%B0%D0%BB%2C%20%D0%B4%D1%80%D0%B0%D0%BC%D0%B0%2C%20%D0%BA%D0%BE%D0%BC%D0%B8%D0%BA%D1%81%2C%20UHD%20BDRemux%202160p%2C%20HDR10%2C%20Dolby%20Vision%5D%20Dub%20(iTunes)%20%2B%203x%20M",
	}

	for _, testCase := range testCases {
		parsed, err := magneturi.Parse(testCase)
		require.NoError(t, err)

		caseUnescape, err := url.PathUnescape(testCase)
		require.NoError(t, err)
		encodedUnescape, err := url.PathUnescape(parsed.Encoded())
		require.NoError(t, err)

		require.Len(t, encodedUnescape, len(caseUnescape))
	}
}
