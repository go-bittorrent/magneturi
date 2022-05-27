package magneturi

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

const (
	magnetPrefix = "magnet:?"
)

// Magnet is a parsed URI.
type Magnet struct {
	// DisplayName is a "dp".
	DisplayName string
	// ExactTopics is a "xt".
	ExactTopics []string
	// ExactLength is a "xl".
	ExactLength int64
	// AcceptableSources is "as".
	AcceptableSources []string
	// ExactSource is a "xs".
	ExactSource []string
	// KeywordTopic is a "kt".
	KeywordTopic []string
	// ManifestTopic is a "mt".
	ManifestTopic string
	// Trackers is a "tr".
	Trackers []string
	// AdditionParams is a "x." and other unparsed params.
	AdditionParams map[string][]string
}

// ErrUnsupportedFormat ...
var ErrUnsupportedFormat = errors.New("uri doesn't starts from prefix")

// Parse magnet URI.
func Parse(rawURI string) (*Magnet, error) {
	m := Magnet{
		AdditionParams: make(map[string][]string),
	}

	prefixSplit := strings.Split(rawURI, magnetPrefix)
	if len(prefixSplit) <= 1 {
		return nil, ErrUnsupportedFormat
	}

	// Create maps for unique values.
	exactTopics := make(map[string]struct{})
	trackers := make(map[string]struct{})
	acceptableSources := make(map[string]struct{})

	params := strings.Split(prefixSplit[1], "&")

	for _, param := range params {
		paramSplit := strings.Split(param, "=")

		if len(paramSplit) < 2 || paramSplit[1] == "" {
			continue
		}

		key := paramSplit[0]
		value := paramSplit[1]

		switch key {
		case "dn":
			decoded, err := url.QueryUnescape(value)
			if err != nil {
				return nil, err
			}

			m.DisplayName = decoded
		case "xt":
			if _, exist := exactTopics[value]; !exist {
				exactTopics[value] = struct{}{}
			}
		case "kt":
			m.KeywordTopic = append(m.KeywordTopic, strings.Split(value, "+")...)
		case "mt":
			m.ManifestTopic = value
		case "tr":
			decoded, err := url.QueryUnescape(value)
			if err != nil {
				return nil, err
			}

			if _, exist := trackers[decoded]; !exist {
				trackers[decoded] = struct{}{}
			}
		case "as":
			decoded, err := url.QueryUnescape(value)
			if err != nil {
				return nil, err
			}

			if _, exist := acceptableSources[decoded]; !exist {
				acceptableSources[decoded] = struct{}{}
			}
		case "xl":
			size, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}

			m.ExactLength = size
		case "xs":
			decoded, err := url.QueryUnescape(value)
			if err != nil {
				return nil, err
			}

			m.ExactSource = append(m.ExactSource, decoded)
		default:
			decoded, err := url.QueryUnescape(value)
			if err != nil {
				return nil, err
			}

			m.AdditionParams[key] = append(m.AdditionParams[key], decoded)
		}
	}

	m.ExactTopics = make([]string, 0, len(exactTopics))
	for topic := range exactTopics {
		m.ExactTopics = append(m.ExactTopics, topic)
	}

	m.Trackers = make([]string, 0, len(trackers))
	for tracker := range trackers {
		m.Trackers = append(m.Trackers, tracker)
	}

	m.AcceptableSources = make([]string, 0, len(acceptableSources))
	for as := range acceptableSources {
		m.AcceptableSources = append(m.AcceptableSources, as)
	}

	return &m, nil
}

// Encoded returns encoded magnet URI.
func (m *Magnet) Encoded() string {
	r := magnetPrefix + "dn=" + url.QueryEscape(m.DisplayName)

	for _, xt := range m.ExactTopics {
		r += "&xt=" + xt
	}

	if m.ExactLength > 0 {
		r += "&xl=" + strconv.FormatInt(m.ExactLength, 10)
	}

	for _, as := range m.AcceptableSources {
		r += "&as=" + url.QueryEscape(as)
	}

	for _, xs := range m.ExactSource {
		r += "&xs=" + url.QueryEscape(xs)
	}

	if len(m.KeywordTopic) > 0 {
		r += "&kt=" + strings.Join(m.KeywordTopic, "+")
	}

	if m.ManifestTopic != "" {
		r += "&mt=" + m.ManifestTopic
	}

	for _, tr := range m.Trackers {
		r += "&tr=" + url.QueryEscape(tr)
	}

	for key, param := range m.AdditionParams {
		for _, paramValue := range param {
			r += "&" + key + "=" + paramValue
		}
	}

	return r
}
