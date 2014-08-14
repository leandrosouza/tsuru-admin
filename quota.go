// Copyright 2014 tsuru-admin authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tsuru/tsuru/cmd"
	"github.com/tsuru/tsuru/quota"
	"net/http"
	"strconv"
)

type viewUserQuota struct{}

func (viewUserQuota) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "view-user-quota",
		MinArgs: 1,
		Usage:   "view-user-quota <user-email>",
		Desc:    "Displays the current usage and limit of the user",
	}
}

func (viewUserQuota) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/users/" + context.Args[0] + "/quota")
	if err != nil {
		return err
	}
	request, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var quota quota.Quota
	err = json.NewDecoder(resp.Body).Decode(&quota)
	if err != nil {
		return err
	}
	fmt.Fprintf(context.Stdout, "User: %s\n", context.Args[0])
	fmt.Fprintf(context.Stdout, "Apps owned: %d\n", quota.InUse)
	fmt.Fprintf(context.Stdout, "Limit of apps: %d\n", quota.Limit)
	return nil
}

type changeUserQuota struct{}

func (changeUserQuota) Info() *cmd.Info {
	desc := `Changes the limit of apps that a user can create

The new limit must be an integer, it may also be "unlimited".`
	return &cmd.Info{
		Name:    "change-user-quota",
		MinArgs: 2,
		Usage:   "change-user-quota <user-email> <new-limit>",
		Desc:    desc,
	}
}

func (changeUserQuota) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/users/" + context.Args[0] + "/quota")
	if err != nil {
		return err
	}
	limit, err := parseLimit(context.Args[1])
	if err != nil {
		return err
	}
	body := bytes.NewBufferString("limit=" + limit)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	fmt.Fprintln(context.Stdout, "Quota successfully updated.")
	return nil
}

type changeAppQuota struct{}

func (changeAppQuota) Info() *cmd.Info {
	desc := `Changes the limit of units that an app can have

The new limit must be an integer, it may also be "unlimited".`
	return &cmd.Info{
		Name:    "change-app-quota",
		MinArgs: 2,
		Usage:   "change-app-quota <user-email> <new-limit>",
		Desc:    desc,
	}
}

func (changeAppQuota) Run(context *cmd.Context, client *cmd.Client) error {
	url, err := cmd.GetURL("/apps/" + context.Args[0] + "/quota")
	if err != nil {
		return err
	}
	limit, err := parseLimit(context.Args[1])
	if err != nil {
		return err
	}
	body := bytes.NewBufferString("limit=" + limit)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = client.Do(request)
	if err != nil {
		return err
	}
	fmt.Fprintln(context.Stdout, "Quota successfully updated.")
	return nil
}

func parseLimit(value string) (string, error) {
	if value == "unlimited" {
		return "-1", nil
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return "", errors.New(`invalid limit. It must be either an integer or "unlimited"`)
	}
	return strconv.Itoa(n), nil
}