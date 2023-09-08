⚠️ This plugin has been deprecated. See https://github.com/getgoing/drone-validation instead.

[![Build Status](https://cloud.drone.io/api/badges/teryaev/drone-promote-auth/status.svg)](https://cloud.drone.io/teryaev/drone-promote-auth)


A validation extension to Drone validation plugin to restrict users who can promote builds. _Please note this project requires Drone server version 1.4 or higher._

Docker image -- https://hub.docker.com/r/reptiloid666/drone-promote-auth

## Installation

Create a shared secret:

```console
$ openssl rand -hex 16
bea26a2221fd8090ea38720fc445eca6
```

Download and run the plugin:

```console
$ docker run -d \
  --publish=3000:3000 \
  --env=DRONE_DEBUG=true \
  --env=DRONE_SECRET=bea26a2221fd8090ea38720fc445eca6 \
  --env=PRIVILEGED_USERS=user1,user2,user3 \
  --env=USER_PERMISSIONS=user4:uat;prod,user5:uat \
  --restart=always \
  --name=drone-promote-auth reptiloid666/drone-promote-auth
```

PRIVILEGED_USERS -- comma-separated list of users that are authorized to promote to any env
USER_PERMISSIONS -- comma-separated list of maps where key is username and value is semicolon-separated list of envs the user is authorised to promote to

Update your Drone server configuration to include the plugin address and the shared secret.

```text
DRONE_VALIDATE_PLUGIN_ENDPOINT=http://1.2.3.4:3000
DRONE_VALIDATE_PLUGIN_SECRET=bea26a2221fd8090ea38720fc445eca6
