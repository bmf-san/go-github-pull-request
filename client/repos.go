package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// Repository is a repository.
type Repository struct {
	Name string
}

// HandlerRepo is a handler for repos.
type HandlerRepo struct {
	httpClient *HTTPClient
}

// GetOrgRepoParams is parameters for GetOrgRepos.
type GetOrgRepoParams struct {
	PerPage int
	Page    int
}

// NewHandlerRepo creates a handler.
func NewHandlerRepo(httpClient *HTTPClient) *HandlerRepo {
	return &HandlerRepo{
		httpClient: httpClient,
	}
}

// GetOrgsRepos get organization repositories.
func (hr *HandlerRepo) GetOrgRepos(org string, params GetOrgRepoParams) ([]*Repository, error) {
	url := &url.URL{}
	query := url.Query()
	query.Add("per_page", strconv.Itoa(params.PerPage))
	query.Add("page", strconv.Itoa(params.Page))
	// TODO:

	req, err := hr.httpClient.NewRequest(http.MethodGet, nil, nil, fmt.Sprintf("orgs/%v/repos?%v", org, query.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := hr.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var r []*Repository
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, err
	}

	return r, nil
}
