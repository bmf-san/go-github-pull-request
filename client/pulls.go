package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// PullRequest is a pull request.
// see: https://github.com/google/go-github/blob/master/github/pulls.go
type PullRequest struct {
	ID        int64      `json:"id,omitempty"`
	Number    int        `json:"number,omitempty"`
	State     string     `json:"state,omitempty"`
	Title     string     `json:"title,omitempty"`
	Body      string     `json:"body,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"`
	HTMLURL   string     `json:"html_url,omitempty"`
	Assignee  *Assignee  `csv:"-"`
}

// GetPullParams is parameters for GetPull.
type GetPullParams struct {
	PerPage int
	State   string
	Page    int
}

// HandlerPull is a handler for pulls.
type HandlerPull struct {
	httpClient *HTTPClient
}

// NewHandlerPull creates a handler.
func NewHandlerPull(httpClient *HTTPClient) *HandlerPull {
	return &HandlerPull{
		httpClient: httpClient,
	}
}

// GetPulls gets pulls by assignee.
func (hp *HandlerPull) GetPullsByAssignee(assignee string, owner string, repo string, params GetPullParams) ([]*PullRequest, error) {
	url := &url.URL{}
	query := url.Query()
	query.Add("per_page", strconv.Itoa(params.PerPage))
	query.Add("state", params.State)
	query.Add("page", strconv.Itoa(params.Page))

	req, err := hp.httpClient.NewRequest(http.MethodGet, nil, nil, fmt.Sprintf("repos/%v/%v/pulls?%v", owner, repo, query.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body, err := hp.httpClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var pr []*PullRequest
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, err
	}

	var rslt []*PullRequest
	for _, p := range pr {
		if p.Assignee != nil {
			if p.Assignee.Login == assignee {
				rslt = append(rslt, p)
			}
		}
	}

	return rslt, nil
}
