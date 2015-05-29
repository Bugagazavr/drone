package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/drone/drone/pkg/settings"
	common "github.com/drone/drone/pkg/types"
)

// Client communicates with a Remote plugin using the
// net/rpc protocol.
type Client struct {
	SkipVerfiy  bool
	Base        string
	ServiceAddr string
	client      *http.Client
}

// New returns a new, remote datastore backend that connects
// via tcp and exchanges data using Go's RPC mechanism.
func New(service *settings.Service) (*Client, error) {
	httpClient := &http.Client{}

	client := &Client{
		SkipVerfiy:  service.SkipVerify,
		Base:        service.Base,
		ServiceAddr: service.Address,
		client:      httpClient,
	}

	return client, nil
}

func (c *Client) Request(method, path string, body []byte) ([]byte, error) {
	var request *http.Request
	var err error
	var url = fmt.Sprintf("%s/%s", c.ServiceAddr, path)

	switch method {
	case "GET":
		request, err = http.NewRequest("GET", url, nil)
	case "POST":
		request, err = http.NewRequest("POST", url, bytes.NewReader(body))
	case "PUT":
		request, err = http.NewRequest("PUT", url, bytes.NewReader(body))
	case "DELETE":
		request, err = http.NewRequest("DELETE", url, bytes.NewReader(body))
	default:
		err = errors.New(fmt.Sprintf("Unsupported method: %s, for url: %s", method, url))
	}

	if err != nil {
		return nil, err
	}

	request.Header.Set("X-Drone-Base", c.Base)
	request.Header.Set("X-Drone-Skip-Verify", strconv.FormatBool(c.SkipVerfiy))

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)

	return content, err
}

func (c *Client) Login(token, secret string) (*common.CaseUser, error) {
	var user *common.CaseUser

	body, err := json.Marshal(struct {
		Token  string `json:"token"`
		Secret string `json:"secret"`
	}{
		Token:  token,
		Secret: secret,
	})

	if err != nil {
		return nil, err
	}

	content, err := c.Request("POST", "auth", body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &user); err != nil {
		return nil, err
	}

	return user, nil
}

// Repo fetches the named repository from the remote system.
func (c *Client) Repo(u *common.CaseUser, owner, repo string) (*common.Repo, error) {
	var repository *common.Repo

	body, err := json.Marshal(struct {
		User  *common.CaseUser `json:"user"`
		Owner string           `json:"owner"`
		Repo  string           `json:"repo"`
	}{
		User:  u,
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		return nil, err
	}

	content, err := c.Request("POST", "repo", body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &repository); err != nil {
		return nil, err
	}

	return repository, nil
}

func (c *Client) Perm(u *common.CaseUser, owner, repo string) (*common.Perm, error) {
	var permission *common.Perm

	body, err := json.Marshal(struct {
		User  *common.CaseUser `json:"user"`
		Owner string           `json:"owner"`
		Repo  string           `json:"repo"`
	}{
		User:  u,
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		return nil, err
	}

	content, err := c.Request("POST", "perm", body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(content, &permission); err != nil {
		return nil, err
	}

	fmt.Println(permission)
	return permission, nil
}

func (c *Client) Script(u *common.User, r *common.Repo, b *common.Commit) ([]byte, error) {
	return nil, nil
}

func (c *Client) Status(u *common.User, r *common.Repo, b *common.Commit) error {
	return nil
}

func (c *Client) Activate(u *common.User, r *common.Repo, k *common.Keypair, link string) error {
	return nil
}

func (c *Client) Deactivate(u *common.User, r *common.Repo, link string) error {
	return nil
}

func (c *Client) Netrc(u *common.User) (*common.Netrc, error) {
	url_, err := url.Parse(c.Base)
	if err != nil {
		return nil, err
	}
	netrc := &common.Netrc{}
	netrc.Login = u.Token
	netrc.Password = "x-oauth-basic"
	netrc.Machine = url_.Host
	return netrc, nil
}

func (c *Client) Orgs(u *common.User) ([]string, error) {
	return []string{}, nil
}

func (c *Client) Hook(r *http.Request) (*common.Hook, error) {
	hook := new(common.Hook)
	header := make(http.Header)
	copyHeader(r.Header, header)

	return hook, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
