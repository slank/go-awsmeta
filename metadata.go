package awsmeta

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/juju/errors"
)

const DEFAULT_API_VERSION = "latest"
const DEFAULT_TIMEOUT = 2
const AWS_MD_URL = "http://169.254.169.254"

type MetaDataServer struct {
	ApiVersion string
	Timeout    int
}

func (s *MetaDataServer) Get(path string) (string, error) {
	var api_version = DEFAULT_API_VERSION
	if s.ApiVersion == "" {
		api_version = s.ApiVersion
	}

	var timeout = DEFAULT_TIMEOUT
	if s.Timeout != 0 {
		timeout = s.Timeout
	}

	md_url := fmt.Sprintf("%s/%s/%s", AWS_MD_URL, api_version, path)

	var h http.Client
	h.Timeout = time.Duration(timeout) * time.Second
	resp, err := h.Get(md_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", errors.Errorf("(not found)")
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Shortcuts for metadata values
var ShortNames = map[string]string{
	"az":          "meta-data/placement/availability-zone",
	"instance-id": "meta-data/instance-id",
}
