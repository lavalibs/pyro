package lavalink

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/lavalibs/pyro/lavalink/types"
)

// HTTP handles HTTP requests to Lavalink
type HTTP struct {
	Client   *http.Client
	Endpoint string
	Password string
}

// LoadTrack from the given endpoint
func (h *HTTP) LoadTrack(identifier string) (track types.TrackResponse, err error) {
	res, err := h.DoJSON(http.MethodGet, "/loadtracks", nil, func(q url.Values) {
		q.Add("identifier", identifier)
	})
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&track)
	res.Body.Close()
	return
}

// DecodeTrack from the given endpoint
func (h *HTTP) DecodeTrack(identifier string) (track types.TrackInfo, err error) {
	res, err := h.DoJSON(http.MethodGet, "/decodetrack", nil, func(q url.Values) {
		q.Add("track", identifier)
	})
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&track)
	res.Body.Close()
	return
}

// DecodeTracks from the given endpoint
func (h *HTTP) DecodeTracks(identifiers []string) (tracks []types.TrackInfo, err error) {
	res, err := h.DoJSON(http.MethodPost, "/decodetracks", &identifiers, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&tracks)
	res.Body.Close()
	return
}

// DoJSON performs a JSON HTTP request to the Lavalink node
func (h *HTTP) DoJSON(method, path string, d interface{}, genQuery func(q url.Values)) (res *http.Response, err error) {
	var (
		b []byte
		r io.Reader
	)

	if d != nil {
		b, err = json.Marshal(d)
		if err != nil {
			return
		}

		r = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, h.Endpoint, r)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", h.Password)
	req.URL.Path = path

	if genQuery != nil {
		q := req.URL.Query()
		genQuery(q)
		req.URL.RawQuery = q.Encode()
	}

	if h.Client == nil {
		h.Client = http.DefaultClient
	}
	res, err = h.Client.Do(req)
	return
}
