# Versions Proxy

A proxy providing version information from various repositories.
Mainly created to overcome the shortcomings of the [Artifactory](https://jfrog.com/artifactory/) vs. [Renovate](https://github.com/renovatebot/renovate) combo in situations where the Artifactory server is the only server with access to the Internet.

Unfortunately, Artifactory has proven to be quite bad at providing what Renovate needs to suggest updates in our internal Git repo.

* New versions are not reported for GitHub and GitLab artifacts unless someone already downloaded.
* Release timestamps are not reported for Docker Hub images, making it impossible to use Renovate's [quarantine feature](https://docs.renovatebot.com/key-concepts/minimum-release-age/).

If this proxy runs alongside Artifactory (i.e., on the host with access to the Internet) and Renovate is properly configured to use it via [custom data sources](https://docs.renovatebot.com/modules/datasource/custom/), these problems will be solved.

## Usage

```text
$ ./versions -h
Usage: versions [option ...]

  -p, --port=PORT    web server port to listen to
  -c, --config=FILE  configuration file in YAML format
```

You will most likely just give the `--config` option and have the port specified in the configuration file.
Please use the [default configuration file](internal/config/config.default.yaml) as a template for your own configuration.

In order to have TLS support, you need to run the Versions Proxy behind a terminating web server.
Let's say you are running on the default port (8086) on a computer with nginx.
In that case, you may add the following inside the `location /` block of the nginx configuration:

```text
location ^~ /_v/ {
  rewrite ^/_v/(.*)$ /$1 break;
  proxy_pass http://127.0.0.1:8086/;
}
```

This will make the Versions Proxy available as "https://myserver.example.com/_v/".
Then you may add custom Renovate data sources like this:

```json
"customDatasources": {
  "github-releases": {
    "defaultRegistryUrlTemplate": "https://myserver.example.com/_v/github-releases/{{packageName}}"
  },
  "github-tags": {
    "defaultRegistryUrlTemplate": "https://myserver.example.com/_v/github-tags/{{packageName}}"
  },
  "gitlab-releases": {
    "defaultRegistryUrlTemplate": "https://myserver.example.com/_v/gitlab-releases/{{packageName}}"
  }
}
```

You may then use these data sources as you would use the built-in `github-releases`, `github-tags` and `gitlab-releases`, just prefix with `custom.`, i.e., `custom.github-releases`.
And voila! Renovate will be able to suggest updates. 

## TODO

In semi-prioritized order:

* Gracefully handle error codes from downstream.
  Currently, everything bad from the backends leads to 500.
  Need to report back on rate limiting errors, at least.
* Implement authentication support for Maven and Docker Hub.
* Add more data sources.
