# versions

**This is work in early progress**

A proxy providing version information from various repositories.
Mainly created to overcome the shortcomings of the Artifactory vs [Renovate](https://github.com/renovatebot/renovate) combo where I currently work, in which Artifactory in the only application with access to the Internet.
Unfortunately, Artifactory has proven to be quite bad at providing what Renovate needs to suggest updates in our internal Git repo.

* New versions are not reported for GitHub and GitLab artifacts, unless someone already downloaded.
* Release timestamp is not reported for Dockerhub images, making it impossible to use Renovate's [quarantine feature](https://docs.renovatebot.com/key-concepts/minimum-release-age/).

If this proxy runs alongside Artifactory (i.e. on the host with access to the Internet), and Renovate is properly configured to use it via [custom data sources](https://docs.renovatebot.com/modules/datasource/custom/), these problems will be solved.

## TODO

In semi-prioritized order:

* Support authentication with the data sources, in order to enable more queries (avoid rate limiting for anonymous).
* Gracefully handle error codes from downstream.
  Currently, everything bad from the backends lead to 500.
  Need to report back on rate limiting errors, at least.
* Support pagination for data sources; currently only the first page is parsed.
* Use native Renovate format for output.
* Persistent and perhaps shared cache.
* ~~Where applicable, include release notes.~~
* Add more data sources.
  This is the last item because I think that when everything above is in place, it will be possible to see what can be extracted to avoid too much duplicate code.
  There is duplication already, but it is my intention to clean it up, eventually.

## Want to contribute?

I'm a self-made guy.
I don't like external dependencies.
(I made an exception for the YAML dependency since I'm new to Go.
I already implemented a JSON encoder/decoder in my old language, Java, which was faster than any other implementation at that time.)
Every contribution that doesn't add external dependencies without good reason, is welcome.
Contributions that make me better at Go, are especially welcome.
