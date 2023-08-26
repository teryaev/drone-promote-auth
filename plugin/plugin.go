// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"encoding/csv"
	"strings"

	log "github.com/sirupsen/logrus"

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
func New(privilegedUsers []string, userPermissionsRaw string) validator.Plugin {
	//userPermissionsRaw sturcture `0.userName, 1.env, 2.repoName`
	userData, err := csv.NewReader(strings.NewReader(userPermissionsRaw)).ReadAll()
	if err != nil {
		log.Error(err)
	}

	userPermissions := make(map[string]map[string][]string)
	for _, row := range userData {
		userMap, ok := userPermissions[row[0]]
		if !ok {
			userMap = make(map[string][]string)
			userPermissions[row[0]] = userMap
		}
		userEnv, ok := userMap[row[1]]
		if !ok {
			userEnv = make([]string, 0)
			userMap[row[1]] = userEnv
		}
		userMap[row[1]] = append(userEnv, row[2])
	}

	return &plugin{
		privilegedUsers: privilegedUsers,
		userPermissions: userPermissions,
	}
}

type plugin struct {
	privilegedUsers []string
	userPermissions map[string]map[string][]string
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	fields := log.Fields{
		"user":  req.Build.Trigger,
		"event": req.Build.Event,
		"env":   req.Build.Deploy,
		"repo":  req.Repo.Name,
	}
	log.WithFields(fields).WithField("config", req.Config.Data).Info("Validating build")

	// check if this event requires auth
	if stringInSlice(req.Build.Event, restrictedEvents) {
		// check if env is staging and avaialable for all users
		if strings.HasPrefix(req.Build.Deploy, "staging") {
			log.WithFields(fields).Info(
				"User has been authorized for staging env",
			)
			return nil
		}

		// check if user is privilged to promote to any env
		if stringInSlice(req.Build.Trigger, p.privilegedUsers) {
			log.WithFields(fields).Info(
				"User has been authorized as a privileged user",
			)
			return nil
		}

		// check if user has any per-env[repo] permission
		if allowedEnvs, userHasPermissions := p.userPermissions[req.Build.Trigger]; userHasPermissions {
			// check if user is allowed to promote to a requested env

			for env, repos := range allowedEnvs {
				if env == req.Build.Deploy && stringInSlice(req.Repo.Name, repos) {
					log.WithFields(fields).Info(
						"User has been authorized according to user level permissions",
					)
					return nil
				}
			}

		}
		log.WithFields(fields).Info("User is not allowed")
		return validator.ErrSkip
	}

	return nil
}
