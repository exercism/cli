# Next Release

----------------

* [#151](https://github.com/exercism/cli/pull/151): Expand '~' in config path to home directory - [@lcowell](https://github.com/lcowell)
* Stop supporting legacy config files (~/.exercism.go)
* Deleted deprecated login/logout commands
* Your contribution here

## v1.9.2 (Jan 11, 2015)

* [exercism.io#2155](https://github.com/exercism/exercism.io/issues/2155): Fixed problem with passed in config file being ignored.
* Added first version of changelog

## v1.9.1 (Jan 10, 2015)

* [#147](https://github.com/exercism/cli/pull/147): added `--api` option to exercism configure - [@morphatic](https://github.com/morphatic)

## v1.9.0 (Nov 27, 2014)

* [#143](https://github.com/exercism/cli/pull/143): added command for downloading a specific solution - [@harimp](https://github.com/harimp)
* [#142](https://github.com/exercism/cli/pull/142): fixed command name to be `exercism` rather than `cli` on `go get` - [@Tonkpils](https://github.com/Tonkpils)

## v1.8.2 (Oct 24, 2014)

* [9cbd069](https://github.com/exercism/cli/commit/9cbd06916cc05bbb165e8c2cb00d5e03cb4dbb99): Made path comparison case insensitive

## v1.8.1 (Oct 23, 2014)

* [0ccc7a4](https://github.com/exercism/cli/commit/0ccc7a479940d2d7bb5e12eab41c91105519f135): Implemented debug flag on submit command

## v1.8.0 (Oct 15, 2014)

* [#138](https://github.com/exercism/cli/pull/138): Added conversion to line endings for submissions on Windows - [@rprouse](https://github.com/rprouse)
* [#116](https://github.com/exercism/cli/issues/116): Added support for setting name of config file in an environment variable
* [47d6fd4](https://github.com/exercism/cli/commit/47d6fd407fd0410f5c81d60172e01e8624608f53): Added a `track` command to list the problems in a given language
* [#126](https://github.com/exercism/cli/issues/126): Added explanation in `submit` response about fetching the next problems
* [#133](https://github.com/exercism/cli/pull/133): Changed config command to create the exercism directory, rather than waiting until the first time problems are fetched - [@Tonkpils](https://github.com/Tonkpils)

## v1.7.5 (Oct 5, 2014)

* [88cf1a1fbc884545dfc10e98535f667e4a43e693](https://github.com/exercism/cli/commit/88cf1a1fbc884545dfc10e98535f667e4a43e693): Added ARMv6 to build
* [12672c4](https://github.com/exercism/cli/commit/12672c4f695cfe3891f96467619a3615e6d57c34): Added an error message when people submit a file that is not within the exercism directory tree
* [#128](https://github.com/exercism/cli/pull/128): Made paths os-agnostic in tests - [@ccnp123](https://github.com/ccnp123)

## v1.7.4 (Sep 27, 2014)

* [4ca3e97](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Fixed incorrect HOME directory on Windows
* [8bd1a25](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Added ARMv5 to build
* [#117](https://github.com/exercism/cli/pull/117): Archive windows binaries using zip rather than tar and gzip - [@LegalizeAdulthood](https://github.com/LegalizeAdulthood)

## v1.7.3 (Sep 26, 2014)

* [8bec393](https://github.com/exercism/cli/commit/8bec39387094680990af7cf438ada1780cf87129): Fixed submit so it can handle symlinks

## v1.7.2 (Sep 24, 2014)

* [#111](https://github.com/exercism/cli/pull/111): Don't clobber existing config values when adding more - [@jish](https://github.com/jish)

## v1.7.1 (Sep 19, 2014)

* Completely reorganized the code, separating each command into a separate handler
* [17fc164](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Fixed paths to be platform-independent
* [8b174e2](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Made the output of demo command more helpful
* [8b174e2](https://github.com/exercism/cli/commit/8b174e2fd8c7a545ea5c47c998ac10c5a7ab371f): Deleted the 'current' command

## v1.7.0 (Aug 28, 2014)

* [ac6dbfd](https://github.com/exercism/cli/commit/ac6dbfd81a86e7a9a5a9b68521b0226c40d8e813): Added os and architecture to the user agent
* [5d58fd1](https://github.com/exercism/cli/commit/5d58fd14b9db84fb752b3bf6112123cd6f04c532): Fixed bug in detecting user's home directory
* [#100](https://github.com/exercism/cli/pull/100): Added 'debug' command, which supersedes the 'info' command - [@Tonkpils](https://github.com/Tonkpils)
* Extracted a couple of commands into separate handlers
* [6ec5876](https://github.com/exercism/cli/commit/6ec5876bde0b02206cacbe685bb8aedcbdba25d4): Added a hack to rename old config files to the new default name
* [bb7d0d6](https://github.com/exercism/cli/commit/bb7d0d6151a950c92590dc771ec3ff5fdd1c83b0): Rename 'home' command to 'info'
* [#95](https://github.com/exercism/cli/issues/95): Added 'home' command
* Deprecate login/logout commands
* [1a39134](https://github.com/exercism/cli/commit/1a391342da93aa32ae398f1500a3981aa65b9f41): Changed demo to write exercises to the default exercism problems directory
* [07cc334](https://github.com/exercism/cli/commit/07cc334739465b21d6eb5d973e16e1c88f67758e): Deleted the whoami command, we weren't using github usernames for anything
* [#97](https://github.com/exercism/cli/pull/97): Changed default exercism directory to ~/exercism - [@lcowell](https://github.com/lcowell)
* [#94](https://github.com/exercism/cli/pull/94): Updated language detection to handle C++ - [@LegalizeAdulthood](https://github.com/LegalizeAdulthood)
* [#92](https://github.com/exercism/cli/pull/92): Renamed config json file to .exercism.json instead of .exercism.go - [@lcowell](https://github.com/lcowell)
* [f55653f](https://github.com/exercism/cli/commit/f55653f35863914086a54375afb0898e142c1638): Deleted go vet from travis build temporarily until the codebase can be cleaned up
* [#91](https://github.com/exercism/cli/pull/91): Replaced temp file usage with encode/decode - [@lcowell](https://github.com/lcowell)
* [#90](https://github.com/exercism/cli/pull/90): Added sanitization to config values to trim whitespace before writing it - [@lcowell](https://github.com/lcowell)
* Did a fair amount of cleanup to make code a bit more idiomatic
* [#86](https://github.com/exercism/cli/pull/86): Triggered interactive login command for commands that require auth - [@Tonkpils](https://github.com/Tonkpils)

## v1.6.2 (Jun 2, 2014)

* [a5b7a55](https://github.com/exercism/cli/commit/a5b7a55f52c23ac5ce2c6bd1826ea7767aea38c4): Update login prompt

## v1.6.1 (May 16, 2014)

* [#84](https://github.com/exercism/cli/pull/84): Change hard-coded filepath so that it will work on any platform - [@simonjefford](https://github.com/simonjefford)

## v1.6.0 (May 10, 2014)

* [#82](https://github.com/exercism/cli/pull/82): Fixed typo in tests - [@srt32](https://github.com/srt32)
* [aa7446d](https://github.com/exercism/cli/commit/aa7446d598fc894ef329756555c48ef358baf676): Clarified output to user after they fetch
* [#79](https://github.com/exercism/cli/pull/79): Updated development instructions to fix permissions problem - [@andrewsardone](https://github.com/andrewsardone)
* [#78](https://github.com/exercism/cli/pull/78): Deleted deprecated action `peek` - [@djquan](https://github.com/djquan)
* [#74](https://github.com/exercism/cli/pull/74): Implemented new option on `fetch` to get a single language - [@Tonkpils](https://github.com/Tonkpils)
* [#75](https://github.com/exercism/cli/pull/75): Improved feedback to user after logging in - [@Tonkpils](https://github.com/Tonkpils)
* [#72](https://github.com/exercism/cli/pull/72): Optimized use of temp file - [@Dparker1990](https://github.com/Dparker1990)
* [#70](https://github.com/exercism/cli/pull/70): Fixed a panic - [@Tonkpils](https://github.com/Tonkpils)
* [#68](https://github.com/exercism/cli/pull/68): Fixed how user input is read so that it doesn't stop at the first space - [@Tonkpils](https://github.com/Tonkpils)

-------------

**WIP** - I'm slowly building up the full changelog history for the project.

