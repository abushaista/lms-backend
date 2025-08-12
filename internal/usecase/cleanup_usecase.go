package usecase

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CleanUpUseCase struct {
	validator *validator.Validate
}

func NewCleanUpUsecase() *CleanUpUseCase {
	return &CleanUpUseCase{
		validator: validator.New(),
	}
}

// func (c *CleanUpUseCase) CleanOperation(url string, operation string) (string, error) {
// 	newUrl := url
// 	switch operation {
// 	case "all":
// 		parsed, err := canonicalURL(url)
// 		if err != nil {
// 			return parsed, err
// 		}
// 		newUrl = parsed
// 		parsed, err = redirection(newUrl)
// 		if err != nil {
// 			return parsed, err
// 		}
// 		newUrl = parsed
// 	case "canonical":
// 		parsed, err := canonicalURL(url)
// 		if err != nil {
// 			return parsed, err
// 		}
// 		newUrl = parsed
// 	case "redirection":
// 		parsed, err := redirection(newUrl)
// 		if err != nil {
// 			return parsed, err
// 		}
// 		newUrl = parsed
// 	}
// 	return newUrl, nil
// }

func (c *CleanUpUseCase) CanonicalURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Lowercase scheme and host
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	if (u.Scheme == "http" && strings.HasSuffix(u.Host, ":80")) ||
		(u.Scheme == "https" && strings.HasSuffix(u.Host, ":443")) {
		u.Host = strings.Split(u.Host, ":")[0]
	}

	if u.Path != "/" {
		u.Path = strings.TrimRight(u.Path, "/")
	}
	u.RawQuery = ""

	u.Fragment = ""

	return u.String(), nil
}

func (c *CleanUpUseCase) Redirection(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = "www.byfood.com"

	return strings.ToLower(u.String()), nil
}
