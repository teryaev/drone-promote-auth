// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/drone/drone-go/plugin/validator"
)

var (
	restrictedEvents = []string{
		"promote",
		"rollback",
	}
)

func stringInSlice(str string, slice []string) bool {
	for _, item := range slice {
		if str == item {
			return true
		}
	}
	return false
}

// New returns a new validator plugin.
func New(allowedUsers []string) validator.Plugin {
	return &plugin{
		allowedUsers: allowedUsers,
	}
}

type plugin struct {
	allowedUsers []string
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	logrus.Debugf("Received %s request from %s", req.Build.Event, req.Build.Trigger)
	if stringInSlice(req.Build.Event, restrictedEvents) &&
		(!stringInSlice(req.Build.Trigger, p.allowedUsers)) {
		logrus.Debugf("user %s not allowed to %s", req.Build.Trigger, req.Build.Event)
		return validator.ErrSkip
	}

	return nil
}
