package harbor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go.f110.dev/xerrors"
)

const (
	userAgent = "harbor-client/1.0"
)

type roundTripper struct {
	http.RoundTripper
	username string
	password string
}

func (rt *roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.SetBasicAuth(rt.username, rt.password)
	r.Header.Set("User-Agent", userAgent)

	return rt.RoundTripper.RoundTrip(r)
}

type Harbor struct {
	host     string
	username string
	password string

	client *http.Client
}

func New(host, username, password string) *Harbor {
	h := &Harbor{
		host:     host,
		username: username,
		password: password,
		client:   &http.Client{},
	}
	h.client.Transport = &roundTripper{
		RoundTripper: http.DefaultTransport.(*http.Transport).Clone(),
		username:     username,
		password:     password,
	}

	return h
}

func (h *Harbor) SetTransport(t http.RoundTripper) {
	h.client.Transport.(*roundTripper).RoundTripper = t
}

func (h *Harbor) ListProjects() ([]Project, error) {
	req, err := h.newRequest(http.MethodGet, "projects", nil)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		// Succeeded
	case http.StatusUnauthorized:
		return nil, xerrors.New("unauthorized")
	default:
		return nil, xerrors.Newf("harbor: list project. unknown status code: %d", res.StatusCode)
	}

	projects := make([]Project, 0)
	if err := json.NewDecoder(res.Body).Decode(&projects); err != nil {
		return nil, xerrors.WithStack(err)
	}
	return projects, nil
}

func (h *Harbor) ExistProject(name string) (bool, error) {
	req, err := h.newRequest(http.MethodHead, "projects?project_name="+name, nil)
	if err != nil {
		return false, xerrors.WithStack(err)
	}
	res, err := h.client.Do(req)
	if err != nil {
		return false, xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, xerrors.Newf("harbor: exists project. unknown status code: %d", res.StatusCode)
	}
}

func (h *Harbor) NewProject(p *NewProjectRequest) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return xerrors.WithStack(err)
	}

	req, err := h.newRequest(http.MethodPost, "projects", bytes.NewReader(buf.Bytes()))
	if err != nil {
		return xerrors.WithStack(err)
	}
	res, err := h.client.Do(req)
	if err != nil {
		return xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
	// Succeeded
	case http.StatusConflict:
		return xerrors.Newf("%s already exists", p.ProjectName)
	default:
		return xerrors.Newf("harbor: new project. unknown status code: %d", res.StatusCode)
	}

	return nil
}

func (h *Harbor) DeleteProject(projectId int) error {
	req, err := h.newRequest(http.MethodDelete, fmt.Sprintf("projects/%d", projectId), nil)
	if err != nil {
		return xerrors.WithStack(err)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return xerrors.Newf("invalid project id: %d", projectId)
	case http.StatusNotFound:
		return xerrors.New("project not found")
	default:
		return xerrors.Newf("harbor: delete project. unknown status code: %d", res.StatusCode)
	}
}

func (h *Harbor) CreateRobotAccount(projectId int, robotRequest *NewRobotAccountRequest) (*RobotAccount, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(robotRequest); err != nil {
		return nil, xerrors.WithStack(err)
	}

	req, err := h.newRequest(http.MethodPost, fmt.Sprintf("projects/%d/robots", projectId), bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	res, err := h.client.Do(req)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
	// Succeeded
	case http.StatusUnauthorized:
		return nil, xerrors.New("create robot account: not logged in")
	default:
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, xerrors.WithStack(err)
		}
		return nil, xerrors.Newf("create robot acount: unknown status code: %d %s", res.StatusCode, string(b))
	}

	result := &RobotAccount{}
	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		return nil, xerrors.WithStack(err)
	}

	return result, nil
}

func (h *Harbor) DeleteRobotAccount(projectId, robotId int) error {
	req, err := h.newRequest(http.MethodDelete, fmt.Sprintf("projects/%d/robots/%d", projectId, robotId), nil)
	if err != nil {
		return xerrors.WithStack(err)
	}
	res, err := h.client.Do(req)
	if err != nil {
		return xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
	// Succeeded
	case http.StatusNotFound:
		return xerrors.New("robot account is not found")
	default:
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return xerrors.WithStack(err)
		}
		return xerrors.Newf("delete robot account: unknown status code: %d %s", res.StatusCode, string(b))
	}

	return nil
}

func (h *Harbor) GetRobotAccounts(projectId int) ([]*RobotAccount, error) {
	req, err := h.newRequest(http.MethodGet, fmt.Sprintf("projects/%d/robots", projectId), nil)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	res, err := h.client.Do(req)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		// Succeeded
	case http.StatusBadRequest, http.StatusNotFound:
		return nil, xerrors.New("get robot accounts: project id is not found or invalid")
	case http.StatusUnauthorized:
		return nil, xerrors.New("get robot accounts: not logged in")
	default:
		return nil, xerrors.Newf("get robot accounts: unknown status code: %d", res.StatusCode)
	}

	result := make([]*RobotAccount, 0)
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, xerrors.WithStack(err)
	}

	return result, nil
}

func (h *Harbor) newRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, fmt.Sprintf("%s/api/v2.0/%s", h.host, endpoint), body)
	if err != nil {
		return nil, xerrors.WithStack(err)
	}
	r.SetBasicAuth(h.username, h.password)
	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "appliaction/json")

	return r, nil
}

type Project struct {
	Id       int             `json:"project_id,omitempty"`
	OwnerId  int             `json:"owner_id,omitempty"`
	Name     string          `json:"name"`
	Metadata ProjectMetadata `json:"metadata"`
}

type NewProjectRequest struct {
	ProjectName  string          `json:"project_name"`
	CVEWhitelist CVEWhitelist    `json:"cve_whitelist,omitempty"`
	CountLimit   int             `json:"count_limit,omitempty"`
	StorageLimit int             `json:"storage_limit,omitempty"`
	Metadata     ProjectMetadata `json:"metadata,omitempty"`
}

type ProjectMetadata struct {
	Public               string `json:"public,omitempty"`
	EnableContentTrust   string `json:"enable_content_trust,omitempty"`
	AutoScan             string `json:"auto_scan,omitempty"`
	Severity             string `json:"severity,omitempty"`
	ReuseSysCVEWhitelist string `json:"reuse_sys_cve_whitelist,omitempty"`
	PreventVUL           string `json:"prevent_vul,omitempty"`
}

type CVEWhitelist struct {
	Items []CVEItem `json:"items,omitempty"`
}

type CVEItem struct {
	CVEId string `json:"cve_id"`
}

type RobotAccount struct {
	Id           int    `json:"id,omitempty"`
	ProjectId    int    `json:"project_id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Token        string `json:"token,omitempty"`
	Disabled     bool   `json:"disabled,omitempty"`
	ExpiresAt    int    `json:"expires_at,omitempty"`
	CreationTime string `json:"creation_time,omitempty"`
	UpdateTime   string `json:"update_time,omitempty"`
}

type NewRobotAccountRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Access      []Access `json:"access,omitempty"`
}

type Access struct {
	Action   string `json:"action"`
	Resource string `json:"resource"`
}