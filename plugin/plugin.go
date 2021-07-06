// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"strings"

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
func New(privilegedUsers []string, userPermissionsRaw map[string]string) validator.Plugin {
	userPermissions := make(map[string][]string)

	// parse list of envs each user is allowed to promote builds to
	for user, envString := range userPermissionsRaw {
		userPermissions[user] = strings.Split(envString, ";")
	}

	return &plugin{
		privilegedUsers: privilegedUsers,
		userPermissions: userPermissions,
	}
}

type plugin struct {
	privilegedUsers []string
	userPermissions map[string][]string
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	logrus.Debugf("Received %s request from %s", req.Build.Event, req.Build.Trigger)
	logrus.Debugf("Targeted env is %s", req.Build.Target)
	// check if this event requires auth
	if stringInSlice(req.Build.Event, restrictedEvents) {
		// check if user is privilged to promote to any env
		if stringInSlice(req.Build.Trigger, p.privilegedUsers) {
			logrus.Debugf(
				"User %s has been authorized to %s %s env as a privileged user",
				req.Build.Trigger, req.Build.Event, req.Build.Target,
			)
			return nil
		}

		// check if user has any per-env permission
		if allowedEnvs, userHasPermissions := p.userPermissions[req.Build.Trigger]; userHasPermissions {
			// check if user is allowed to promote to a requested env
			if stringInSlice(req.Build.Target, allowedEnvs) {
				logrus.Debugf(
					"User %s has been authorized to %s %s env according to user level permissions",
					req.Build.Trigger, req.Build.Event, req.Build.Target,
				)
				return nil
			}
		}

		logrus.Debugf("user %s not allowed to %s to %s", req.Build.Trigger, req.Build.Event, req.Build.Target)
		return validator.ErrSkip
	}

	return nil
}
