# build/runner

This docker image bundles MikTeX and the `assignmentctl` executable in an Ubuntu-based container image.

To use the image, run e.g.

```bash
# Run the docker image interactively
$ docker run -ti ghcr.io/zoomoid/assignments/runner:latest assignmentctl build --all 
```

## Building containers from scratch

You can build the container from scratch by running 

```bash
# build the runner container that includes a LaTeX distro
$ docker build -t assignments/runner:latest .
```

> Note that the Dockerfile bases the image on the image `ghcr.io/zoomoid/assignments/cli:latest`.
> If you cannot pull that image, see [build/cli](../cli/README.md), on how to build the CLI locally.
> Instead of changing the runner's dockerfile, build the CLI image locally with a tag, e.g.,
> `assignments/cli:latest`, and then run the above build command with an additional 
> `--build-arg=IMAGE=assignments/cli:latest`.
