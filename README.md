# Archlinux Repository Upload Server

This small program can be used to updating a package on an Archlinux repository
using a single HTTP request. It is intended to be used in conjunction with a
continuous integration / development pipeline with a stage that deploys an
Archlinux package to a repository.

## Usage

This server can be quickly deployed using the available Docker image.

```
$ docker run --rm                             \
      --volume <PATH TO REPOSITORY>:/srv/repo \
      --publish 8080:8080                     \
      --env PACKAGER_TOKEN=<SECRET TOKEN>     \
      gpol/archlinux-repo-upload-server:latest
```

## Example

With the server deployed, the packager can easily update a package using a
simple `curl` command.

```
$ curl http://example.com:8080/upload         \
      -F "file=@curl-7.67.0-3-arm.pkg.tar.xz" \
      -F "arch=armv7h"                        \
      -F "repo=core"                          \
      -F "token=my-secret-token"
```

## TODO

* Add HTTP/S support
