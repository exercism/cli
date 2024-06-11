# Exercism Command-line Interface (CLI)

[![CI](https://github.com/exercism/cli/actions/workflows/ci.yml/badge.svg)](https://github.com/exercism/cli/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/exercism/cli)](https://goreportcard.com/report/github.com/exercism/cli)

The CLI is the link between the [Exercism][exercism] website and your local work environment. It lets you download exercises and submit your solution to the site.

This CLI ships as a binary with no additional runtime requirements.

## Installing the CLI

Instructions can be found at [exercism/cli/releases](https://github.com/exercism/cli/releases)

- If you already setup your porject by following [Go Module](https://go.dev/doc/modules/layout)


```
git install github.com/exercism/cli/exercism@latest
```

## Tips: 

```
git mod tidy -go=1.20 -compat=1.20  # make sure the version change with this git mod tidy this will update the version of go if required.
```
## Contributing

If you wish to help improve the CLI, please see the [Contributing guide][contributing].

[exercism]: http://exercism.io
[contributing]: /CONTRIBUTING.md
