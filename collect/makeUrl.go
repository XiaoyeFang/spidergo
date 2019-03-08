package collect

import (
	"log"
	"net/http"
	"net/url"
	"crawler-apkmirror/config"
)

func MakeCheckUrl(link string) (*http.Request, error) {
	r, err := url.Parse(link)
	if err != nil {
		log.Fatalln(err)
	}
	q := r.Query()
	r.RawQuery = q.Encode()
	req, err := http.NewRequest(
		"GET",
		r.String(),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", config.HTTP_ACCEPT)
	req.Header.Set("User-Agent", config.HTTP_USER_AGENT)

	return req, err

}
