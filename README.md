# Mattermost plugin minotor

[![Build Status](https://github.com/mattermost/mattermost-plugin-starter-template/actions/workflows/ci.yml/badge.svg)](https://github.com/mattermost/mattermost-plugin-starter-template/actions/workflows/ci.yml)
[![E2E Status](https://github.com/mattermost/mattermost-plugin-starter-template/actions/workflows/e2e.yml/badge.svg)](https://github.com/mattermost/mattermost-plugin-starter-template/actions/workflows/e2e.yml)

Export some Mattermost metrics for Prometheus:

```
# HELP mattermost_kpi_channels_total Total number of channels
# TYPE mattermost_kpi_channels_total gauge
mattermost_kpi_channels_total 2
# HELP mattermost_kpi_last_post_date Timestamp of last post date
# TYPE mattermost_kpi_last_post_date gauge
mattermost_kpi_last_post_date 1.734186913838e+12
# HELP mattermost_kpi_posts_total Total number of posts
# TYPE mattermost_kpi_posts_total gauge
mattermost_kpi_posts_total 39
# HELP mattermost_system_database_status Database component status
# TYPE mattermost_system_database_status gauge
mattermost_system_database_status 1
# HELP mattermost_system_filestore_status Filestore component status
# TYPE mattermost_system_filestore_status gauge
mattermost_system_filestore_status 1
# HELP mattermost_system_status Global status
# TYPE mattermost_system_status gauge
mattermost_system_status 1
# HELP mattermost_usage_posts_total Total number of posts
# TYPE mattermost_usage_posts_total gauge
mattermost_usage_posts_total{group="mattermost"} 40
# HELP mattermost_usage_storage_bytes Storage usage in bytes
# TYPE mattermost_usage_storage_bytes gauge
mattermost_usage_storage_bytes 6e+07
# HELP mattermost_usage_users_total Total number of users
# TYPE mattermost_usage_users_total gauge
mattermost_usage_users_total 7
# HELP promhttp_metric_handler_errors_total Total number of internal errors encountered by the promhttp metric handler.
# TYPE promhttp_metric_handler_errors_total counter
promhttp_metric_handler_errors_total{cause="encoding"} 0
promhttp_metric_handler_errors_total{cause="gathering"} 0
```

## Getting Started

### Installation

Download and install the plugin <https://developers.mattermost.com/integrate/plugins/using-and-managing-plugins/>

### Usage

Configure your Prometheus or Grafana alloy to scrap endpoint: `/plugins/org.ebuildy.plugin-minotor/metrics`

<!> There is no authentication, you must protect this endpoint from public access <!>

## Development

To avoid having to manually install your plugin, build and deploy your plugin using one of the following options. In order for the below options to work, you must first enable plugin uploads via your config.json or API and restart Mattermost.

```json
    "PluginSettings" : {
        ...
        "EnableUploads" : true
    }
```

### Deploying with Local Mode

If your Mattermost server is running locally, you can enable [local mode](https://docs.mattermost.com/administration/mmctl-cli-tool.html#local-mode) to streamline deploying your plugin. Edit your server configuration as follows:

```json
{
    "ServiceSettings": {
        ...
        "EnableLocalMode": true,
        "LocalModeSocketLocation": "/var/tmp/mattermost_local.socket"
    },
}
```

and then deploy your plugin:
```
make deploy
```

You may also customize the Unix socket path:
```bash
export MM_LOCALSOCKETPATH=/var/tmp/alternate_local.socket
make deploy
```

If developing a plugin with a webapp, watch for changes and deploy those automatically:
```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make watch
```

### Deploying with credentials

Alternatively, you can authenticate with the server's API with credentials:
```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```

or with a [personal access token](https://docs.mattermost.com/developer/personal-access-tokens.html):
```bash
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_TOKEN=j44acwd8obn78cdcx7koid4jkr
make deploy
```

### Releasing new versions

The version of a plugin is determined at compile time, automatically populating a `version` field in the [plugin manifest](plugin.json):
* If the current commit matches a tag, the version will match after stripping any leading `v`, e.g. `1.3.1`.
* Otherwise, the version will combine the nearest tag with `git rev-parse --short HEAD`, e.g. `1.3.1+d06e53e1`.
* If there is no version tag, an empty version will be combined with the short hash, e.g. `0.0.0+76081421`.

To disable this behaviour, manually populate and maintain the `version` field.

## How to Release

To trigger a release, follow these steps:

1. **For Patch Release:** Run the following command:
    ```
    make patch
    ```
   This will release a patch change.

2. **For Minor Release:** Run the following command:
    ```
    make minor
    ```
   This will release a minor change.

3. **For Major Release:** Run the following command:
    ```
    make major
    ```
   This will release a major change.

4. **For Patch Release Candidate (RC):** Run the following command:
    ```
    make patch-rc
    ```
   This will release a patch release candidate.

5. **For Minor Release Candidate (RC):** Run the following command:
    ```
    make minor-rc
    ```
   This will release a minor release candidate.

6. **For Major Release Candidate (RC):** Run the following command:
    ```
    make major-rc
    ```
   This will release a major release candidate.

## Q&A

### How do I make a server-only or web app-only plugin?

Simply delete the `server` or `webapp` folders and remove the corresponding sections from `plugin.json`. The build scripts will skip the missing portions automatically.

### How do I include assets in the plugin bundle?

Place them into the `assets` directory. To use an asset at runtime, build the path to your asset and open as a regular file:

```go
bundlePath, err := p.API.GetBundlePath()
if err != nil {
    return errors.Wrap(err, "failed to get bundle path")
}

profileImage, err := ioutil.ReadFile(filepath.Join(bundlePath, "assets", "profile_image.png"))
if err != nil {
    return errors.Wrap(err, "failed to read profile image")
}

if appErr := p.API.SetProfileImage(userID, profileImage); appErr != nil {
    return errors.Wrap(err, "failed to set profile image")
}
```
