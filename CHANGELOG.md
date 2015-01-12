# Next Release

----------------

* [exercism.io#2155](https://github.com/exercism/exercism.io/issues/2155): Fixed problem with passed in config file being ignored.
* Added first version of changelog
* Your contribution here

## v1.9.1 (Jan 10, 2015)

* [#147](https://github.com/exercism/cli/pull/147): added `--api` option to exercism configure - @morphatic

## v1.9.0 (Nov 27, 2014)

* [#143](https://github.com/exercism/cli/pull/143): added command for downloading a specific solution - @harimp
* [#142](https://github.com/exercism/cli/pull/142): fixed command name to be `exercism` rather than `cli` on `go get` - @Tonkpils

## v1.8.2 (Oct 24, 2014)

* [9cbd069](https://github.com/exercism/cli/commit/9cbd06916cc05bbb165e8c2cb00d5e03cb4dbb99): Made path comparison case insensitive

## v1.8.1 (Oct 23, 2014)

* [0ccc7a4](https://github.com/exercism/cli/commit/0ccc7a479940d2d7bb5e12eab41c91105519f135): Implemented debug flag on submit command

## v1.8.0 (Oct 15, 2014)

* [#138](https://github.com/exercism/cli/pull/138): Added conversion to line endings for submissions on Windows - @rprouse
* [#116](https://github.com/exercism/cli/issues/116): Added support for setting name of config file in an environment variable
* [47d6fd4](https://github.com/exercism/cli/commit/47d6fd407fd0410f5c81d60172e01e8624608f53): Added a `track` command to list the problems in a given language
* [#126](https://github.com/exercism/cli/issues/126): Added explanation in `submit` response about fetching the next problems
* [#133](https://github.com/exercism/cli/pull/133): Changed config command to create the exercism directory, rather than waiting until the first time problems are fetched - @Tonkpils

## v1.7.5 (Oct 5, 2014)

* [88cf1a1fbc884545dfc10e98535f667e4a43e693](https://github.com/exercism/cli/commit/88cf1a1fbc884545dfc10e98535f667e4a43e693): Added ARMv6 to build
* [12672c4](https://github.com/exercism/cli/commit/12672c4f695cfe3891f96467619a3615e6d57c34): Added an error message when people submit a file that is not within the exercism directory tree
* [#128](https://github.com/exercism/cli/pull/128): Made paths os-agnostic in tests - @ccnp123

## v1.7.4 (Sep 27, 2014)

* [4ca3e97](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Fixed incorrect HOME directory on Windows
* [8bd1a25](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Added ARMv5 to build
* [#117](https://github.com/exercism/cli/pull/117): Archive windows binaries using zip rather than tar and gzip - @LegalizeAdulthood

## v1.7.3 (Sep 26, 2014)

* [8bec393](https://github.com/exercism/cli/commit/8bec39387094680990af7cf438ada1780cf87129): Fixed submit so it can handle symlinks

## v1.7.2 (Sep 24, 2014)

* [#111](https://github.com/exercism/cli/pull/111): Don't clobber existing config values when adding more - @jish

## v1.7.1 (Sep 19, 2014)

* Completely reorganized the code, separating each command into a separate handler
* [17fc164](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Fixed paths to be platform-independent
* [8b174e2](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Made the output of demo command more helpful
* [8b174e2](https://github.com/exercism/cli/commit/8b174e2fd8c7a545ea5c47c998ac10c5a7ab371f): Deleted the 'current' command

-------------

**WIP** - I'm slowly building up the full changelog history for the project.


