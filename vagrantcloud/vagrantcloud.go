package vagrantcloud

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jochasinga/requests"
	"github.com/pkg/errors"
)

var urlFormat = "https://vagrantcloud.com/%s"

type Provider string

const (
	// Providers
	Parallels     Provider = "parallels"
	Hyperv        Provider = "hyperv"
	Libvirt       Provider = "libvirt"
	Virtualbox    Provider = "virtualbox"
	VmwareDesktop Provider = "vmware_desktop"
)

type BoxJson struct {
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Name             string `json:"name"`
	Versions         []struct {
		Version             string `json:"version"`
		Status              string `json:"status"`
		DescriptionHTML     string `json:"description_html"`
		DescriptionMarkdown string `json:"description_markdown"`
		Providers           []struct {
			Name         string `json:"name"`
			URL          string `json:"url"`
			Checksum     string `json:"checksum"`
			ChecksumType string `json:"checksum_type"`
		} `json:"providers"`
	} `json:"versions"`
}

/*
Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"
  config.vm.box_version = "2004.01"
end

boxName: centos/7
version: 2004.01
*/

func GetBox(boxName, version string, p Provider) (io.ReadCloser, error) {
	addMimeType := func(r *requests.Request) {
		r.Header.Add("Accept", "*/*")
	}

	resp, err := requests.Get(fmt.Sprintf(urlFormat, boxName), addMimeType)
	if err != nil {
		return nil, errors.Errorf("failed to get boxes: %+v", err)
	}
	defer resp.Body.Close()

	var box BoxJson
	if err := json.NewDecoder(resp.Body).Decode(&box); err != nil {
		return nil, errors.Errorf("failed to decode boxes: %+v", err)
	}

	for _, v := range box.Versions {
		if v.Version == version {
			for _, provider := range v.Providers {
				if provider.Name == string(p) {
					resp, err := requests.Get(provider.URL, addMimeType)
					if err != nil {
						return nil, errors.Errorf("failed to get box error: %w", err)
					}
					return resp.Body, nil
				}
			}
		}
	}

	return nil, nil
}

func NewBoxReader(reader io.Reader) (*tar.Reader, error) {
	greader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, errors.Errorf("failed to NewGzipReader: %+v", err)
	}
	defer greader.Close()

	treader := tar.NewReader(greader)
	if err != nil {
		return nil, errors.Errorf("failed to NewTarReader: %+v", err)
	}
	return treader, err
}
