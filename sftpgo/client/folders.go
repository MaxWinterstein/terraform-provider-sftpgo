// Copyright (C) 2023 Nicola Murino
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sftpgo/sdk"
)

// GetFolders - Returns list of folders
func (c *Client) GetFolders() ([]sdk.BaseVirtualFolder, error) {
	var result []sdk.BaseVirtualFolder
	limit := 100

	for {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/folders?limit=%d&offset=%d", c.HostURL, limit, len(result)), nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req, http.StatusOK)
		if err != nil {
			return nil, err
		}

		var folders []sdk.BaseVirtualFolder
		err = json.Unmarshal(body, &folders)
		if err != nil {
			return nil, err
		}
		result = append(result, folders...)
		if len(folders) < limit {
			break
		}
	}

	return result, nil
}

// CreateFolder - creates a new folder
func (c *Client) CreateFolder(folder sdk.BaseVirtualFolder) (*sdk.BaseVirtualFolder, error) {
	rb, err := json.Marshal(folder)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/folders", c.HostURL), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var newFolder sdk.BaseVirtualFolder
	err = json.Unmarshal(body, &newFolder)
	return &newFolder, err
}

// GetFolder - Returns a specifc folder
func (c *Client) GetFolder(name string) (*sdk.BaseVirtualFolder, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/folders/%s", c.HostURL, url.PathEscape(name)), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var folder sdk.BaseVirtualFolder
	err = json.Unmarshal(body, &folder)
	return &folder, err
}

// UpdateFolder - Updates an existing folder
func (c *Client) UpdateFolder(folder sdk.BaseVirtualFolder) error {
	rb, err := json.Marshal(folder)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/folders/%s", c.HostURL, url.PathEscape(folder.Name)),
		bytes.NewBuffer(rb))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, http.StatusOK)
	return err
}

// DeleteFolder - Deletes a folder
func (c *Client) DeleteFolder(name string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/folders/%s", c.HostURL, url.PathEscape(name)), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, http.StatusOK)
	return err
}
