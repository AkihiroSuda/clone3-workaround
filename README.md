# `clone3-workaround`: Workaround for running `ubuntu:21.10`, `fedora:35`, and other glibc >= 2.34 distros on Docker <= 20.10.9

Old container engines such as [Docker <= 20.10.9](https://github.com/moby/moby/pull/42836) cannot run
[glibc >= 2.34](https://github.com/bminor/glibc/commit/d8ea0d0168b190bdf138a20358293c939509367f) images such as `ubuntu:21.10` and `fedora:35`:

```console
$ docker run -it  --rm ubuntu:21.10
root@862f014171b5:/# apt-get update
Get:1 http://security.ubuntu.com/ubuntu impish-security InRelease [90.7 kB]
Get:2 http://archive.ubuntu.com/ubuntu impish InRelease [270 kB]
Get:3 http://security.ubuntu.com/ubuntu impish-security/main amd64 Packages [620 B]
Get:4 http://archive.ubuntu.com/ubuntu impish-updates InRelease [90.7 kB]
Get:5 http://archive.ubuntu.com/ubuntu impish-backports InRelease [90.7 kB]
Get:6 http://archive.ubuntu.com/ubuntu impish/universe amd64 Packages [16.7 MB]
Get:7 http://archive.ubuntu.com/ubuntu impish/restricted amd64 Packages [110 kB]
Get:8 http://archive.ubuntu.com/ubuntu impish/main amd64 Packages [1793 kB]
Get:9 http://archive.ubuntu.com/ubuntu impish/multiverse amd64 Packages [256 kB]
Get:10 http://archive.ubuntu.com/ubuntu impish-updates/main amd64 Packages [620 B]
Fetched 19.4 MB in 7s (2893 kB/s)
Reading package lists... Done
E: Problem executing scripts APT::Update::Post-Invoke 'rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true'
E: Sub-process returned an error code
```

```console
$ docker run -it --rm fedora:35 dnf update
[root@849f3703c4b5 /]# dnf install -y hello
Fedora 35 - x86_64                                                                                                                                                                                 0.0  B/s |   0  B     00:00
Errors during downloading metadata for repository 'fedora':
  - Curl error (6): Couldn't resolve host name for https://mirrors.fedoraproject.org/metalink?repo=fedora-35&arch=x86_64 [getaddrinfo() thread failed to start]
Error: Failed to download metadata for repo 'fedora': Cannot prepare internal mirrorlist: Curl error (6): Couldn't resolve host name for https://mirrors.fedoraproject.org/metalink?repo=fedora-35&arch=x86_64 [getaddrinfo() thread failed to start]
```

`clone3-workaround` provides a workaround for this issue, by loading an additional seccomp profile that hides `clone3(2)` syscall from glibc, so that
the `clone()` wrapper of glibc works in the legacy-compatible mode.

No need to upgrade Docker. No need to specify custom `docker run --security-opt` flags.

## Target container engines
`clone3-workaround` should be useful for the following containe engines.

- [Docker, prior to 20.10.10](https://github.com/moby/moby/pull/42836)
- [containerd CRI-mode, prior to 1.5.6](https://github.com/containerd/containerd/pull/6013)
- [containerd CRI-mode, prior to 1.4.10](https://github.com/containerd/containerd/pull/6014)

Newer container engines DO NOT need `clone3-workaround`.

Also note that some distributor vendors have already cherry-picked the Docker 20.10.10 patch to older versions.
e.g., [`docker.io/20.10.7-0ubuntu5~20.04.1`](https://bugs.launchpad.net/cloud-images/+bug/1943049) DO NOT need `clone3-workaround`, although its version number is smaller than `20.10.10`.

## Build

Run `make` to build `clone3-workaround` binary.

Dependencies:
- Go
- libseccomp-dev

## Usage

Mount or copy `clone3-workaround` to the container, and run `clone3-workaround COMMAND [ARGUMENTS...]` to run the command with the workaround.

```console
$ docker run -it --rm -v $(pwd)/clone3-workaround:/clone3-workaround ubuntu:21.10 /clone3-workaround bash
root@490fd2f29a88:/# apt-get update
Get:1 http://security.ubuntu.com/ubuntu impish-security InRelease [90.7 kB]
...
Fetched 19.4 MB in 6s (2996 kB/s)
Reading package lists... Done

root@490fd2f29a88:/# apt-get install -y hello
Reading package lists... Done
...
Unpacking hello (2.10-2ubuntu3) ...
Setting up hello (2.10-2ubuntu3) ...
```

```console
$ docker run -it --rm -v $(pwd)/clone3-workaround:/clone3-workaround fedora:35 /clone3-workaround bash
[root@c699df1e7bd4 /]# dnf install -y hello
Fedora 35 - x86_64                                                                                                                                                                                 6.5 MB/s |  61 MB     00:09
...
Installed:
  hello-2.10-6.fc35.x86_64                                                                                          info-6.8-2.fc35.x86_64

Complete!
```
