package services

import "net/url"

// Replace URL-safe symbols using QueryUnescape from the net/url package
func UnescapeURL(encodedStr string) (string, error) {
	unescapedStr, err := url.QueryUnescape(encodedStr)
	if err != nil {
		return "", err
	}
	return unescapedStr, nil
}
