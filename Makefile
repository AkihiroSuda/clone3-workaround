GO ?= go

clone3-workaround: main.go go.mod go.sum
	CGO_ENABLED=1 $(GO) build -ldflags="-extldflags -static" -o $@ .

.PHONY: clean
clean:
	rm -rf clone3-workaround _artifacts

# `make artifacts` requires Ubuntu
.PHONY: artifacts
artifacts: clean
	mkdir _artifacts
	make
	strip clone3-workaround
	mv clone3-workaround ./_artifacts/clone3-workaround.x86_64
	# TODO: cross-compile clone3-workaround.aarch64
	echo "The binary is linked with the following version of libseccomp (LGPL-2.1):" > ./_artifacts/libseccomp-version.txt
	echo "-----" >> ./_artifacts/libseccomp-version.txt
	dpkg-query -s libseccomp-dev >> ./_artifacts/libseccomp-version.txt
	echo "-----" >> ./_artifacts/libseccomp-version.txt
	echo "The source code is available at https://packages.ubuntu.com/search?keywords=libseccomp-dev" >> _artifacts/libseccomp-version.txt
