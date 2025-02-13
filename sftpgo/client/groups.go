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

// GetGroups - Returns list of groups
func (c *Client) GetGroups() ([]sdk.Group, error) {
	var result []sdk.Group
	limit := 100

	for {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/groups?limit=%d&offset=%d", c.HostURL, limit, len(result)), nil)
		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req, http.StatusOK)
		if err != nil {
			return nil, err
		}

		var groups []sdk.Group
		err = json.Unmarshal(body, &groups)
		if err != nil {
			return nil, err
		}
		result = append(result, groups...)
		if len(groups) < limit {
			break
		}
	}

	return result, nil
}

// CreateGroup - creates a new group
func (c *Client) CreateGroup(group sdk.Group) (*sdk.Group, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v2/groups", c.HostURL), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}

	var newGroup sdk.Group
	err = json.Unmarshal(body, &newGroup)
	return &newGroup, err
}

// GetGroup - Returns a specifc group
func (c *Client) GetGroup(name string) (*sdk.Group, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v2/groups/%s", c.HostURL, url.PathEscape(name)), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.doRequest(req, http.StatusOK)
	if err != nil {
		return nil, err
	}

	var group sdk.Group
	err = json.Unmarshal(body, &group)
	return &group, err
}

// UpdateGroup - Updates an existing group
func (c *Client) UpdateGroup(group sdk.Group) error {
	rb, err := json.Marshal(group)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v2/groups/%s", c.HostURL, url.PathEscape(group.Name)),
		bytes.NewBuffer(rb))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req, http.StatusOK)
	return err
}

// DeleteGroup - Deletes a group
func (c *Client) DeleteGroup(name string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v2/groups/%s", c.HostURL, url.PathEscape(name)), nil)
	if err != nil {
		return err
	}
	_, err = c.doRequest(req, http.StatusOK)
	return err
}
