### 1.7.0
* Add `--always-auth` `-a` flag to not check keyring for cached creds (but still try to store them)
* Add `--session-duration` `-t` flag to request a different token lifetime, up to 12 hours (may not be granted)
* Add support for session duration in aliases
* Add `max-duration` subcommand to return the maximum possible session duration request for a profile
* `each` subcommand no longer prints each role to stdout as it iterates

### 1.6.1
* Bugfix for `each` subcommand with exit codes

### 1.6.0
* [internal] migration to go modules

### 1.5.0
* Better handling of pipes (in and out) and exit codes

### 1.4.3
* Allow `--region` `-r` to be used at top level

### 1.4.2
* Update the aws-sdk-go, gomock, and logrus versions.

### 1.4.1
* Add new environment variable AWS_KEYCLOAK_PROFILE, which is set to the aws role (eg. power-devx)

### 1.4.0
* Add `each` subcommand to run something across many envs
#### 1.4.0-a
* Bugfix for `each` filtering

### 1.3.6
* Make deps more stable

### 1.3.5
* Add `--filter|-f` param to `list` subcommand to filter roles by regex (eg. '-f admin' will show only admin roles)
* Bump `keyring` version (for linux folks)

### 1.3.4
* Subcommand `env` only display AWS environment vars
* Running without any keycloak-config will automatically download one

### 1.3.3
* Add `list` subcommand, which displays all available roles
* Reduce auth success page auto-close timeout to 1 second
* Remove `aws` subcommand, since it's easier to invoke `aws` after `--` (as with all other commands)

### 1.3.2
* Improve `open` command to allow `aws-keycloak open <profile>`
* Fix `open` command for govcloud (different signin URL)
* Default to less output when when opening browser
* Set both `AWS_REGION` and `AWS_DEFAULT_REGION` env vars

### 1.3.1
* Govcloud bugfix
* Can specify default AWS region in keycloak config (needed for govcloud keycloak)

### 1.3.0
* *Breaking changes*
* Alias support
* Aliases can specify default regions
* Support for govcloud

### 1.2.3
* Support for interactive commands via stdin

### 1.2.2
* Alias `open` subcommand as `login`

### 1.2.1
* New `open` subcommand opens a browser to the logged in AWS console

### 1.2.0
* *Breaking changes*
* Now the shell environment is based to the child command, instead of being stripped out.
