/**
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package redirector

import (
	"net/http"
	"net/url"
	"strings"

	oauthhandlers "github.com/openshift/origin/pkg/oauthserver/oauth/handlers"
)

const (
	// URLToken in the query of the redirectURL gets replaced with the original request URL, escaped as a query parameter.
	// Example use: https://www.example.com/login?then=${url}
	URLToken = "${url}"

	// ServerRelativeURLToken in the query of the redirectURL gets replaced with the server-relative portion of the original request URL, escaped as a query parameter.
	// Example use: https://www.example.com/login?then=${server-relative-url}
	ServerRelativeURLToken = "${server-relative-url}"

	// QueryToken in the query of the redirectURL gets replaced with the original request URL, unescaped.
	// Example use: https://www.example.com/sso/oauth/authorize?${query}
	QueryToken = "${query}"
)

// NewRedirector returns an oauthhandlers.AuthenticationRedirector that redirects to the specified redirectURL.
// Request URLs missing scheme/host, or with relative paths are resolved relative to the baseRequestURL, if specified.
// The following tokens are replaceable in the query of the redirectURL:
//   ${url} is replaced with the current request URL, escaped as a query parameter. Example: https://www.example.com/login?then=${url}
//   ${query} is replaced with the current request query, unescaped. Example: https://www.example.com/sso/oauth/authorize?${query}
func NewRedirector(baseRequestURL *url.URL, redirectURL string) oauthhandlers.AuthenticationRedirector {
	return &redirector{BaseRequestURL: baseRequestURL, RedirectURL: redirectURL}
}

// NewChallenger returns an oauthhandlers.AuthenticationChallenger that returns a Location header to the specified redirectURL.
// Request URLs missing scheme/host, or with relative paths are resolved relative to the baseRequestURL, if specified.
// The following tokens are replaceable in the query of the redirectURL:
//   ${url} is replaced with the current request URL, escaped as a query parameter. Example: https://www.example.com/login?then=${url}
//   ${query} is replaced with the current request query, unescaped. Example: https://www.example.com/sso/oauth/authorize?${query}
func NewChallenger(baseRequestURL *url.URL, redirectURL string) oauthhandlers.AuthenticationChallenger {
	return &redirector{BaseRequestURL: baseRequestURL, RedirectURL: redirectURL}
}

type redirector struct {
	BaseRequestURL *url.URL
	RedirectURL    string
}

// AuthenticationChallenge returns a Location header to the configured RedirectURL (which should return a challenge)
func (r *redirector) AuthenticationChallenge(req *http.Request) (http.Header, error) {
	redirectURL, err := buildRedirectURL(r.RedirectURL, r.BaseRequestURL, req.URL)
	if err != nil {
		return nil, err
	}
	headers := http.Header{}
	headers.Add("Location", redirectURL.String())
	return headers, nil
}

// AuthenticationRedirect redirects to the configured RedirectURL
func (r *redirector) AuthenticationRedirect(w http.ResponseWriter, req *http.Request) error {
	redirectURL, err := buildRedirectURL(r.RedirectURL, r.BaseRequestURL, req.URL)
	if err != nil {
		return nil
	}
	http.Redirect(w, req, redirectURL.String(), http.StatusFound)
	return nil
}

func buildRedirectURL(redirectTemplate string, baseRequestURL, requestURL *url.URL) (*url.URL, error) {
	if baseRequestURL != nil {
		requestURL = baseRequestURL.ResolveReference(requestURL)
	}
	redirectURL, err := url.Parse(redirectTemplate)
	if err != nil {
		return nil, err
	}
	serverRelativeRequestURL := &url.URL{
		Path:     requestURL.Path,
		RawQuery: requestURL.RawQuery,
	}
	redirectURL.RawQuery = strings.Replace(redirectURL.RawQuery, QueryToken, requestURL.RawQuery, -1)
	redirectURL.RawQuery = strings.Replace(redirectURL.RawQuery, URLToken, url.QueryEscape(requestURL.String()), -1)
	redirectURL.RawQuery = strings.Replace(redirectURL.RawQuery, ServerRelativeURLToken, url.QueryEscape(serverRelativeRequestURL.String()), -1)
	return redirectURL, nil
}
