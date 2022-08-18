# build/cli

We ship the assignmentctl CLI with two different distros,

1. Ubuntu, for use with MikTeX as a LaTeX distribution and automatic building of documents, and
2. Alpine Linux, such that the image is smaller and more easily used in CI runners where you do not need the entire MikTeX installation,
   e.g. by only running the `assignmentctl bundle` and `assignmentctl ci release` commands.

Container images are shipped at

```plain
# Ubuntu image
ghcr.io/zoomoid/assignments/cli:latest

# Alpine Linux image
ghcr.io/zoomoid/assignments/cli:alpine
```

## Building containers from scratch

If you'd like to build the containers on your own, make sure that you have `docker` (or an equivalent image builder) installed and run these commands
from the CLI's project root, i.e., `$REPOSITORY/cli`.

```bash
# Build the ubuntu-based image
$ docker build -t assignmentctl:latest -f build/cli/ubuntu/Dockerfile

# Build the alpine-based image
$ docker build -t assignmentctl:alpine -f build/cli/alpine/Dockerfile
```

Note that the Ubuntu image does not yet contain any LaTeX distro! To build this image also on your own, see [build/runner](../runner/README.md).
