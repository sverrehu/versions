# Versions Proxy

A proxy providing version information from various repositories.
Mainly created to overcome the shortcomings of the [Artifactory](https://jfrog.com/artifactory/) vs [Renovate](https://github.com/renovatebot/renovate) combo where I currently work, in which Artifactory in the only application with access to the Internet.

Unfortunately, Artifactory has proven to be quite bad at providing what Renovate needs to suggest updates in our internal Git repo.

* New versions are not reported for GitHub and GitLab artifacts, unless someone already downloaded.
* Release timestamps are not reported for Docker Hub images, making it impossible to use Renovate's [quarantine feature](https://docs.renovatebot.com/key-concepts/minimum-release-age/).

If this proxy runs alongside Artifactory (i.e. on the host with access to the Internet), and Renovate is properly configured to use it via [custom data sources](https://docs.renovatebot.com/modules/datasource/custom/), these problems will be solved.

## Usage

```text
$ ./versions -h
Usage: versions [option ...]

  -p, --port=PORT    web server port to listen to
  -c, --config=FILE  configuration file in YAML format
```

You will most likely just give the `--config` option and have the port specified in the configuration file.
Please use the [default configuration file](internal/config/config.default.yaml) as a template for your own configuration.

## TODO

In semi-prioritized order:

* Gracefully handle error codes from downstream.
  Currently, everything bad from the backends lead to 500.
  Need to report back on rate limiting errors, at least.
* Implement authentication support for Maven and Docker Hub.
* Add more data sources.
