# bookit

A CLI tool for grabbing books off IRC

## Installation

- Download the platform-specific tarball
- `sudo tar xvfz bookit.tar.gz -C /usr/local/bin`
- `sudo chmod +x /usr/local/bin/bookit`

## Development

### Prerequisites

- Go
- golint - `go get -u github.com/golang/lint/golint`
- gox - `go get -u github.com/mitchellh/gox`
- glide - go dependency manager

### Setup

Create package directory

```
mkdir -p $GOPATH/src/gitub.com/mdzhang
```

Clone bookit

```
git clone git@github.com:mdzhang/bookit.git
```

Change directory

```
cd $GOPATH/src/gitub.com/mdzhang/bookit
```

Install dependencies:

```
glide install
```

### Tasks

Run tests

```
make test
```

Lint code:

```
make lint
```

Compile and generate binary for current platform/architecture

```
make compile
```

Cross-platform packaging

```
make package
```
