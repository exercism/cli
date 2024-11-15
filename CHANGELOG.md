# Change log

The exercism CLI follows [semantic versioning](http://semver.org/).

---

## Next Release

- **Your contribution here**

## v3.5.4 (2024-11-15)

- [#1183](https://github.com/exercism/cli/pull/1183) Add support for Uiua track to `exercism test` - [@vaeng]

## v3.5.3 (2024-11-03)

- [#1178](https://github.com/exercism/cli/pull/1178) Add arm64-assembly test configuration [@keiravillekode]
- [#1177](https://github.com/exercism/cli/pull/1177) refactored exercism.io links to exercism.org [@ladokp]
- [#1165](https://github.com/exercism/cli/pull/1165) Add support for the YAMLScript language [@ingydotnet]

## v3.5.2 (2024-10-09)

- [#1174](https://github.com/exercism/cli/pull/1174) Fix an issue with `exercism completion bash` where the command name is not present in the completion output. - [@petrem]
- [#1172](https://github.com/exercism/cli/pull/1172) Fix `exercism test` command for Batch track - [@bnandras]

## v3.5.1 (2024-08-28)

- [#1162](https://github.com/exercism/cli/pull/1162) Add support for Roc to `exercism test` - [@ageron]

## v3.5.0 (2024-08-22)

- [#1157](https://github.com/exercism/cli/pull/1157) Add support for Batch to `exercism test` - [@GroophyLifefor]
- [#1159](https://github.com/exercism/cli/pull/1159) Fix duplicated `t` alias -
  [@muzimuzhi]

## v3.4.2 (2024-08-20)

- [#1156](https://github.com/exercism/cli/pull/1156) Add `test` command to Shell completions -
  [@muzimuzhi]

## v3.4.1 (2024-08-15)

- [#1152](https://github.com/exercism/cli/pull/1152) Add support for Idris to `exercism test` -
  [@isberg]
- [#1151](https://github.com/exercism/cli/pull/1151) Add support for Cairo to `exercism test` - [@isberg]
- [#1147](https://github.com/exercism/cli/pull/1147) Add support for Arturo to `exercism test` - [@erikschierboom]

## v3.4.0 (2024-05-09)

- [#1126](https://github.com/exercism/cli/pull/1126) Update `exercism test` to use Gradle wrapper to test Java exercise - [@sanderploegsma]
- [#1139](https://github.com/exercism/cli/pull/1139) Add support for Pyret to `exercism test`
- [#1136](https://github.com/exercism/cli/pull/1136) Add support for J to `exercism test` - [@enascimento178]
- [#1070](https://github.com/exercism/cli/pull/1070) `exercism open` does not require specifying the directory (defaults to current directory) - [@halfdan]
- [#1122](https://github.com/exercism/cli/pull/1122) Troubleshoot command suggest to open forum post instead of GitHub issue - [@glennj]
- [#1065](https://github.com/exercism/cli/pull/1065) Update help text for `exercism submit` to indicate specifying files is optional - [@andrerfcsantos]
- [#1140](https://github.com/exercism/cli/pull/1140) Fix release notes link

## v3.3.0 (2024-02-15)

- [#1128](https://github.com/exercism/cli/pull/1128) Fix `exercism test` command not working for the `8th` and `emacs-lisp` tracks - [@glennj]
- [#1125](https://github.com/exercism/cli/pull/1125) Simplify root command description
- [#1124](https://github.com/exercism/cli/pull/1124) Use correct domain for FAQ link [@tomasnorre]

## v3.2.0 (2023-07-28)

- [#1092](https://github.com/exercism/cli/pull/1092) Add `exercism test` command to run the unit tests for nearly any track (inspired by [universal-test-runner](https://github.com/xavdid/universal-test-runner)) - [@xavdid]
- [#1073](https://github.com/exercism/cli/pull/1073) Add `arm64` build for each OS

## v3.1.0 (2022-10-04)

- [#979](https://github.com/exercism/cli/pull/979) Protect existing solutions from being overwritten by 'download' - [@harugo]
- [#981](https://github.com/exercism/cli/pull/981) Check if authorisation header is set before attempting to extract token - [@harugo]
- [#1044](https://github.com/exercism/cli/pull/1044) Submit without specifying files - [@andrerfcsantos]

## v3.0.13 (2019-10-23)

- [#866](https://github.com/exercism/cli/pull/866) The API token outputted during verbose will now be masked by default - [@Jrank2013]
- [#873](https://github.com/exercism/cli/pull/873) Make all errors in cmd package checked - [@avegner]
- [#871](https://github.com/exercism/cli/pull/871) Error message returned if the track is locked - [@Jrank2013]
- [#886](https://github.com/exercism/cli/pull/886) Added GoReleaser config, updated docs, made archive naming adjustments - [@ekingery]

## v3.0.12 (2019-07-07)

- [#770](https://github.com/exercism/cli/pull/770) Print API error messages in submit command - [@Smarticles101]
- [#763](https://github.com/exercism/cli/pull/763) Add Fish shell tab completions - [@John-Goff]
- [#806](https://github.com/exercism/cli/pull/806) Make Zsh shell tab completions work on $fpath - [@QuLogic]
- [#797](https://github.com/exercism/cli/pull/797) Fix panic when submit command is not given args - [@jdsutherland]
- [#828](https://github.com/exercism/cli/pull/828) Remove duplicate files before submitting - [@larson004]
- [#793](https://github.com/exercism/cli/pull/793) Submit handles non 2xx responses - [@jdsutherland]

## v3.0.11 (2018-11-18)

- [#752](https://github.com/exercism/cli/pull/752) Improve error message on upgrade command - [@farisj]
- [#759](https://github.com/exercism/cli/pull/759) Update shell tab completion for bash and zsh - [@nywilken]
- [#762](https://github.com/exercism/cli/pull/762) Improve usage documentation - [@Smarticles101]
- [#766](https://github.com/exercism/cli/pull/766) Tweak messaging to work for teams edition - [@kytrinyx]

## v3.0.10 (2018-10-03)

- official release of v3.0.10-alpha.1 - [@nywilken]

## v3.0.10-alpha.1 (2018-09-21)

- [#739](https://github.com/exercism/cli/pull/739) update maxFileSize error to include filename - [@nywilken]
- [#736](https://github.com/exercism/cli/pull/736) Metadata file .solution.json renamed to metadata.json - [@jdsutherland]
- [#738](https://github.com/exercism/cli/pull/738) Add missing contributor URLs to CHANGELOG - [@nywilken]
- [#737](https://github.com/exercism/cli/pull/737) Remove unused solutions type - [@jdsutherland]
- [#729](https://github.com/exercism/cli/pull/729) Update Oh My Zsh instructions - [@katrinleinweber]
- [#725](https://github.com/exercism/cli/pull/725) Do not allow submission of enormous files - [@sfairchild]
- [#724](https://github.com/exercism/cli/pull/724) Update submit error message when submitting a directory - [@sfairchild]
- [#723](https://github.com/exercism/cli/pull/720) Move .solution.json to hidden subdirectory - [@jdsutherland]

## v3.0.9 (2018-08-29)

- [#720](https://github.com/exercism/cli/pull/720) Make the timeout configurable globally - [@kytrinyx]
- [#721](https://github.com/exercism/cli/pull/721) Handle windows filepaths that accidentally got submitted to the server - [@kytrinyx]
- [#722](https://github.com/exercism/cli/pull/722) Handle exercise directories with numeric suffixes - [@kytrinyx]

## v3.0.8 (2018-08-22)

- [#713](https://github.com/exercism/cli/pull/713) Fix broken support for uuid flag on download command - [@nywilken]

## v3.0.7 (2018-08-21)

- [#705](https://github.com/exercism/cli/pull/705) Fix confusion about path and filepath - [@kytrinyx]
- [#650](https://github.com/exercism/cli/pull/650) Fix encoding problem in filenames - [@williandrade]

## v3.0.6 (2018-07-17)

- [#652](https://github.com/exercism/cli/pull/652) Add support for teams feature - [@kytrinyx]
- [#683](https://github.com/exercism/cli/pull/683) Fix typo in welcome message - [@glebedel]
- [#675](https://github.com/exercism/cli/pull/675) Improve output of troubleshoot command when CLI is unconfigured - [@kytrinyx]
- [#679](https://github.com/exercism/cli/pull/679) Improve error message for failed /ping on configure - [@kytrinyx]
- [#669](https://github.com/exercism/cli/pull/669) Add debug as alias for troubleshoot - [@kytrinyx]
- [#647](https://github.com/exercism/cli/pull/647) Ensure welcome message has full link to settings page - [@kytrinyx]
- [#667](https://github.com/exercism/cli/pull/667) Improve bash completion script - [@cookrn]

## v3.0.5 (2018-07-17)

- [#646](https://github.com/exercism/cli/pull/646) Fix issue with upgrading on Windows - [@nywilken]

## v3.0.4 (2018-07-15)

- [#644](https://github.com/exercism/cli/pull/644) Add better error messages when solution metadata is missing - [@kytrinyx]

## v3.0.3 (2018-07-14)

- [#642](https://github.com/exercism/cli/pull/642) Add better error messages when configuration is needed before download - [@kytrinyx]
- [#641](https://github.com/exercism/cli/pull/641) Fix broken download for uuid flag - [@kytrinyx]
- [#618](https://github.com/exercism/cli/pull/618) Fix broken test in Windows build for relative paths - [@nywilken]
- [#631](https://github.com/exercism/cli/pull/631) Stop accepting token flag on download command - [@kytrinyx]
- [#616](https://github.com/exercism/cli/pull/616) Add shell completion scripts to build artifacts - [@jdsutherland]
- [#624](https://github.com/exercism/cli/pull/624) Tweak command documentation to reflect reality - [@kytrinyx]
- [#625](https://github.com/exercism/cli/pull/625) Fix wildly excessive whitespace in error messages - [@kytrinyx]

## v3.0.2 (2018-07-13)

- [#622](https://github.com/exercism/cli/pull/622) Fix bug with multi-file submission - [@kytrinyx]

## v3.0.1 (2018-07-13)

- [#619](https://github.com/exercism/cli/pull/619) Improve error message for successful configuration - [@kytrinyx]

## v3.0.0 (2018-07-13)

This is a complete rewrite from the ground up to work against the new https://exercism.org site.

## v2.4.1 (2017-07-01)

- [#385](https://github.com/exercism/cli/pull/385) Fix broken upgrades for Windows - [@Tonkpils]

## v2.4.0 (2017-03-24)

- [#344](https://github.com/exercism/cli/pull/344) Make the CLI config paths more XDG friendly - [@narqo]
- [#346](https://github.com/exercism/cli/pull/346) Fallback to UTF-8 if encoding is uncertain - [@petertseng]
- [#350](https://github.com/exercism/cli/pull/350) Add ARMv8 binaries to CLI releases - [@Tonkpils]
- [#352](https://github.com/exercism/cli/pull/352) Fix case sensitivity on slug and track ID - [@Tonkpils]
- [#353](https://github.com/exercism/cli/pull/353) Print confirmation when fetching --all - [@neslom]
- [#356](https://github.com/exercism/cli/pull/356) Resolve symlinks before attempting to read files - [@lcowell]
- [#358](https://github.com/exercism/cli/pull/358) Redact API key from debug output - [@Tonkpils]
- [#359](https://github.com/exercism/cli/pull/359) Add short flag `-m` for submit comment flag - [@jgsqware]
- [#366](https://github.com/exercism/cli/pull/366) Allow obfuscation on configure command - [@dmmulroy]
- [#367](https://github.com/exercism/cli/pull/367) Use supplied confirmation text from API on submit - [@nilbus]

## v2.3.0 (2016-08-07)

- [#339](https://github.com/exercism/cli/pull/339) Don't run status command if API key is missing - [@ests]
- [#336](https://github.com/exercism/cli/pull/336) Add '--all' flag to fetch command - [@neslom]
- [#333](https://github.com/exercism/cli/pull/333) Update references of codegangsta/cli -> urfave/cli - [@manusajith], [@blackerby]
- [#331](https://github.com/exercism/cli/pull/331) Improve usage/help text of submit command - [@manusajith]

## v2.2.6 (2016-05-30)

- [#306](https://github.com/exercism/cli/pull/306) Don't use Fatal to print usage - [@broady]
- [#307](https://github.com/exercism/cli/pull/307) Pass API key when fetching individual exercises - [@kytrinyx]
- [#312](https://github.com/exercism/cli/pull/312) Add missing newline on usage string - [@jppunnett]
- [#318](https://github.com/exercism/cli/pull/318) Show activity stream URL after submitting - [@lcowell]
- [4710640](https://github.com/exercism/cli/commit/4710640751c7a01deb1b5bf8a9a65b611b078c05) - [@lcowell]
- Update codegangsta/cli dependency - [@manusajith], [@lcowell]
- [#320](https://github.com/exercism/cli/pull/320) Add missing newlines to usage strings - [@hjljo]
- [#328](https://github.com/exercism/cli/pull/328) Append solution URL path consistently - [@Tonkpils]

## v2.2.5 (2016-04-02)

- [#284](https://github.com/exercism/cli/pull/284) Update release instructions - [@kytrinyx]
- [#285](https://github.com/exercism/cli/pull/285) Create a copy/pastable release text - [@kytrinyx]
- [#289](https://github.com/exercism/cli/pull/289) Fix a typo in the usage statement - [@AlexWheeler]
- [#290](https://github.com/exercism/cli/pull/290) Fix upgrade command for Linux systems - [@jbaiter]
- [#292](https://github.com/exercism/cli/pull/292) Vendor dependencies - [@Tonkpils]
- [#293](https://github.com/exercism/cli/pull/293) Remove extraneous/distracting details from README - [@Tonkpils]
- [#294](https://github.com/exercism/cli/pull/294) Improve usage statement: alphabetize commands - [@beanieboi]
- [#297](https://github.com/exercism/cli/pull/297) Improve debug output when API key is unconfigured - [@mrageh]
- [#299](https://github.com/exercism/cli/pull/299) List output uses track ID and problem from list - [@Tonkpils]
- [#301](https://github.com/exercism/cli/pull/301) Return error message for unknown track status - [@neslom]
- [#302](https://github.com/exercism/cli/pull/302) Add helpful error message when user tries to submit a directory - [@alebaffa]

## v2.2.4 (2016-01-28)

- [#270](https://github.com/exercism/cli/pull/270) Allow commenting on submission with --comment - [@Tonkpils]
- [#271](https://github.com/exercism/cli/pull/271) Increase timeout to 20 seconds - [@Tonkpils]
- [#273](https://github.com/exercism/cli/pull/273) Guard against submitting spec files and README - [@daveyarwood]
- [#278](https://github.com/exercism/cli/pull/278) Create files with 0644 mode, create missing directories for downloaded solutions - [@petertseng]
- [#281](https://github.com/exercism/cli/pull/281) Create missing directories for downloaded problems - [@petertseng]
- [#282](https://github.com/exercism/cli/pull/282) Remove random encouragement after submitting - [@kytrinyx]
- [#283](https://github.com/exercism/cli/pull/283) Print current configuration after calling configure command - [@kytrinyx]

## v2.2.3 (2015-12-27)

- [#264](https://github.com/exercism/cli/pull/264) Fix version flag to use --version and --v - [@Tonkpils]

## v2.2.2 (2015-12-26)

- [#212](https://github.com/exercism/cli/pull/212) extract path related code from config - [@lcowell]
- [#215](https://github.com/exercism/cli/pull/215) use $XDG_CONFIG_HOME if available - [@lcowell]
- [#248](https://github.com/exercism/cli/pull/248) [#253](https://github.com/exercism/cli/pull/253) add debugging output - [@lcowell]
- [#256](https://github.com/exercism/cli/pull/256) clean up build scripts - [@lcowell]
- [#258](https://github.com/exercism/cli/pull/258) reduce filesystem noise on fetch [@devonestes]
- [#261](https://github.com/exercism/cli/pull/261) improve error message when track and exercise can't be identified on submit - [@anxiousmodernman]
- [#262](https://github.com/exercism/cli/pull/262) encourage iterating to improve after first submission on an exercise - [@eToThePiIPower]

## v2.2.1 (2015-08-11)

- [#200](https://github.com/exercism/cli/pull/200): Add guard to unsubmit command - [@kytrinyx]
- [#204](https://github.com/exercism/cli/pull/204): Improve upgrade failure messages and increase timeout - [@Tonkpils]
- [#206](https://github.com/exercism/cli/pull/207): Fix verbose flag and removed short `-v` - [@zabawaba99]
- [#208](https://github.com/exercism/cli/pull/208): avoid ambiguous or unresolvable exercism paths - [@lcowell]

## v2.2.0 (2015-06-27)

- [b3c3d6f](https://github.com/exercism/cli/commit/b3c3d6fe54c622fc0ee07fdd221c8e8e5b73c8cd): Improve error message on Internal Server Error - [@Tonkpils]
- [#196](https://github.com/exercism/cli/pull/196): Add upgrade command - [@Tonkpils]
- [#194](https://github.com/exercism/cli/pull/194): Fix home expansion on configure update - [@Tonkpils]
- [523c5bd](https://github.com/exercism/cli/commit/523c5bdec5ef857f07b39de738a764589660cd5a): Document release process - [@kytrinyx]

## v2.1.1 (2015-05-13)

- [#192](https://github.com/exercism/cli/pull/192): Loosen up restrictions on --test flag for submissions - [@Tonkpils]
- [#190](https://github.com/exercism/cli/pull/190): Fix bug in home directory expansion for Windows - [@Tonkpils]

## v2.1.0 (2015-05-08)

- [1a2fd1b](https://github.com/exercism/cli/commit/1a2fd1bfb2dba358611a7c3266f935cccaf924b5): Handle config as either directory or file - [@lcowell]
- [#177](https://github.com/exercism/cli/pull/177): Improve JSON error handling and reporting - [@Tonkpils]
- [#178](https://github.com/exercism/cli/pull/178): Add support for $XDG_CONFIG_HOME - [@lcowell]
- [#184](https://github.com/exercism/cli/pull/184): Handle different file encodings in submissions - [@ambroff]
- [#179](https://github.com/exercism/cli/pull/179): Pretty print the JSON config - [@Tonkpils]
- [#181](https://github.com/exercism/cli/pull/181): Fix path issue when downloading problems - [@Tonkpils]
- [#186](https://github.com/exercism/cli/pull/186): Allow people to specify a target directory for the demo - [@Tonkpils]
- [#189](https://github.com/exercism/cli/pull/189): Implement `--test` flag to allow submitting a test file in the solution - [@pminten]

## v2.0.2 (2015-04-01)

- [#174](https://github.com/exercism/cli/issues/174): Fix panic during fetch - [@kytrinyx]
- Refactor handling of ENV vars - [@lcowell]

## v2.0.1 (2015-03-25)

- [#167](https://github.com/exercism/cli/pull/167): Fixes misspelling of exercism list command - [@queuebit]
- Tweak output from `fetch` so that languages are scannable.
- [#35](https://github.com/exercism/cli/issues/35): Add support for submitting multiple-file solutions
- [#171](https://github.com/exercism/cli/pull/171): Implement `skip` command to bypass individual exercises - [@Tonkpils]

## v2.0.0 (2015-03-05)

Added:

- [#154](https://github.com/exercism/cli/pull/154): Add 'list' command to list available exercises for a language - [@lcowell]
- [3551884](https://github.com/exercism/cli/commit/3551884e9f38d6e563b99dae7b28a18d4525455d): Add host connectivity status to debug output. - [@lcowell]
- [#162](https://github.com/exercism/cli/pull/162): Allow users to open the browser from the terminal. - [@zabawaba99]

Removed:

- Stop supporting legacy config files (`~/.exercism.go`)
- Deleted deprecated login/logout commands
- Deleted deprecated key names in config

Fixed:

- [#151](https://github.com/exercism/cli/pull/151): Expand '~' in config path to home directory - [@lcowell]
- [#155](https://github.com/exercism/cli/pull/155): Display problems not yet submitted on fetch API - [@Tonkpils]
- [f999e69](https://github.com/exercism/cli/commit/f999e69e5290cec6c5c9933aecc6fddfad8cf019): Disambiguate debug and verbose flags. - [@lcowell]
- Report 'new' at the bottom after fetching, it's going to be more relevant than 'unchanged', which includes all the languages they don't care about.

Tweaked:

- Set environment variable in build script
- [#153](https://github.com/exercism/cli/pull/153): Refactored configuration package - [@kytrinyx]
- [#157](https://github.com/exercism/cli/pull/157): Refactored API package - [@Tonkpils]

## v1.9.2 (2015-01-11)

- [exercism#2155](https://github.com/exercism/exercism/issues/2155): Fixed problem with passed in config file being ignored.
- Added first version of changelog

## v1.9.1 (2015-01-10)

- [#147](https://github.com/exercism/cli/pull/147): added `--api` option to exercism configure - [@morphatic]

## v1.9.0 (2014-11-27)

- [#143](https://github.com/exercism/cli/pull/143): added command for downloading a specific solution - [@harimp]
- [#142](https://github.com/exercism/cli/pull/142): fixed command name to be `exercism` rather than `cli` on `go get` - [@Tonkpils]

## v1.8.2 (2014-10-24)

- [9cbd069](https://github.com/exercism/cli/commit/9cbd06916cc05bbb165e8c2cb00d5e03cb4dbb99): Made path comparison case insensitive

## v1.8.1 (2014-10-23)

- [0ccc7a4](https://github.com/exercism/cli/commit/0ccc7a479940d2d7bb5e12eab41c91105519f135): Implemented debug flag on submit command

## v1.8.0 (2014-10-15)

- [#138](https://github.com/exercism/cli/pull/138): Added conversion to line endings for submissions on Windows - [@rprouse]
- [#116](https://github.com/exercism/cli/issues/116): Added support for setting name of config file in an environment variable
- [47d6fd4](https://github.com/exercism/cli/commit/47d6fd407fd0410f5c81d60172e01e8624608f53): Added a `track` command to list the problems in a given language
- [#126](https://github.com/exercism/cli/issues/126): Added explanation in `submit` response about fetching the next problems
- [#133](https://github.com/exercism/cli/pull/133): Changed config command to create the exercism directory, rather than waiting until the first time problems are fetched - [@Tonkpils]

## v1.7.5 (2014-10-5)

- [88cf1a1fbc884545dfc10e98535f667e4a43e693](https://github.com/exercism/cli/commit/88cf1a1fbc884545dfc10e98535f667e4a43e693): Added ARMv6 to build
- [12672c4](https://github.com/exercism/cli/commit/12672c4f695cfe3891f96467619a3615e6d57c34): Added an error message when people submit a file that is not within the exercism directory tree
- [#128](https://github.com/exercism/cli/pull/128): Made paths os-agnostic in tests - [@ccnp123]

## v1.7.4 (2014-09-27)

- [4ca3e97](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Fixed incorrect HOME directory on Windows
- [8bd1a25](https://github.com/exercism/cli/commit/4ca3e9743f6d421903c91dfa27f4747fb1081392): Added ARMv5 to build
- [#117](https://github.com/exercism/cli/pull/117): Archive windows binaries using zip rather than tar and gzip - [@LegalizeAdulthood]

## v1.7.3 (2014-09-26)

- [8bec393](https://github.com/exercism/cli/commit/8bec39387094680990af7cf438ada1780cf87129): Fixed submit so it can handle symlinks

## v1.7.2 (2014-09-24)

- [#111](https://github.com/exercism/cli/pull/111): Don't clobber existing config values when adding more - [@jish]

## v1.7.1 (2014-09-19)

- Completely reorganized the code, separating each command into a separate handler
- [17fc164](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Fixed paths to be platform-independent
- [8b174e2](https://github.com/exercism/cli/commit/17fc1644e9fc9ee5aa4e136de11556e65a7b6036): Made the output of demo command more helpful
- [8b174e2](https://github.com/exercism/cli/commit/8b174e2fd8c7a545ea5c47c998ac10c5a7ab371f): Deleted the 'current' command

## v1.7.0 (2014-08-28)

- [ac6dbfd](https://github.com/exercism/cli/commit/ac6dbfd81a86e7a9a5a9b68521b0226c40d8e813): Added os and architecture to the user agent
- [5d58fd1](https://github.com/exercism/cli/commit/5d58fd14b9db84fb752b3bf6112123cd6f04c532): Fixed bug in detecting user's home directory
- [#100](https://github.com/exercism/cli/pull/100): Added 'debug' command, which supersedes the 'info' command - [@Tonkpils]
- Extracted a couple of commands into separate handlers
- [6ec5876](https://github.com/exercism/cli/commit/6ec5876bde0b02206cacbe685bb8aedcbdba25d4): Added a hack to rename old config files to the new default name
- [bb7d0d6](https://github.com/exercism/cli/commit/bb7d0d6151a950c92590dc771ec3ff5fdd1c83b0): Rename 'home' command to 'info'
- [#95](https://github.com/exercism/cli/issues/95): Added 'home' command
- Deprecate login/logout commands
- [1a39134](https://github.com/exercism/cli/commit/1a391342da93aa32ae398f1500a3981aa65b9f41): Changed demo to write exercises to the default exercism problems directory
- [07cc334](https://github.com/exercism/cli/commit/07cc334739465b21d6eb5d973e16e1c88f67758e): Deleted the whoami command, we weren't using github usernames for anything
- [#97](https://github.com/exercism/cli/pull/97): Changed default exercism directory to ~/exercism - [@lcowell]
- [#94](https://github.com/exercism/cli/pull/94): Updated language detection to handle C++ - [@LegalizeAdulthood]
- [#92](https://github.com/exercism/cli/pull/92): Renamed config json file to .exercism.json instead of .exercism.go - [@lcowell]
- [f55653f](https://github.com/exercism/cli/commit/f55653f35863914086a54375afb0898e142c1638): Deleted go vet from travis build temporarily until the codebase can be cleaned up
- [#91](https://github.com/exercism/cli/pull/91): Replaced temp file usage with encode/decode - [@lcowell]
- [#90](https://github.com/exercism/cli/pull/90): Added sanitization to config values to trim whitespace before writing it - [@lcowell]
- Did a fair amount of cleanup to make code a bit more idiomatic
- [#86](https://github.com/exercism/cli/pull/86): Triggered interactive login command for commands that require auth - [@Tonkpils]

## v1.6.2 (2014-06-02)

- [a5b7a55](https://github.com/exercism/cli/commit/a5b7a55f52c23ac5ce2c6bd1826ea7767aea38c4): Update login prompt

## v1.6.1 (2014-05-16)

- [#84](https://github.com/exercism/cli/pull/84): Change hard-coded filepath so that it will work on any platform - [@simonjefford]

## v1.6.0 (2014-05-10)

- [#82](https://github.com/exercism/cli/pull/82): Fixed typo in tests - [@srt32]
- [aa7446d](https://github.com/exercism/cli/commit/aa7446d598fc894ef329756555c48ef358baf676): Clarified output to user after they fetch
- [#79](https://github.com/exercism/cli/pull/79): Updated development instructions to fix permissions problem - [@andrewsardone]
- [#78](https://github.com/exercism/cli/pull/78): Deleted deprecated action `peek` - [@djquan]
- [#74](https://github.com/exercism/cli/pull/74): Implemented new option on `fetch` to get a single language - [@Tonkpils]
- [#75](https://github.com/exercism/cli/pull/75): Improved feedback to user after logging in - [@Tonkpils]
- [#72](https://github.com/exercism/cli/pull/72): Optimized use of temp file - [@Dparker1990]
- [#70](https://github.com/exercism/cli/pull/70): Fixed a panic - [@Tonkpils]
- [#68](https://github.com/exercism/cli/pull/68): Fixed how user input is read so that it doesn't stop at the first space - [@Tonkpils]

## v1.5.1 (2014-03-14)

- [5b672ee](https://github.com/exercism/cli/commit/5b672ee7bf26859c41de9eed83396b7454286063): Provided a visual mark next to new problems that get fetched

## v1.5.0 (2014-02-28)

- [#63](https://github.com/exercism/cli/pull/63): Implemeted `fetch` for a single language - [@Tonkpils]
- [#62](https://github.com/exercism/cli/pull/62): Expose error message from API to user on `fetch` - [@Tonkpils]
- [#59](https://github.com/exercism/cli/pull/59): Added global flag to pass the path to the config file instead of relying on default - [@isbadawi]
- [#57](https://github.com/exercism/cli/pull/57): Added description to the restore command - [@rcode5]
- [#56](https://github.com/exercism/cli/pull/56): Updated developer instructions in README based on real-life experience - [@rcode5]

## v1.4.0 (2014-01-13)

- [#47](https://github.com/exercism/cli/pull/47): Added 'restore' command to download all of a user's existing solutions with their corresponding problems - [@ebautistabar]
- Numerous small fixes and cleanup to code and documentation - [@dpritchett], [@TrevorBramble], [@elimisteve]

## v1.3.2 (2013-12-14)

- [f8dd974](https://github.com/exercism/cli/commit/f8dd9748078b1b191629eae385aaeda8af94305b): Fixed content-type header when posting to API
- Fixed user-agent string

## v1.3.1 (2013-12-01)

- [exercism#1039](https://github.com/exercism/exercism/issues/1039): Stopped clobbering existing files on fetch

## v1.3.0 (2013-11-16)

- [7f39ee4](https://github.com/exercism/cli/commit/7f39ee4802752925466bc2715790dc965026b09d): Allow users to specify a particular problem when fetching.

## v1.2.3 (2013-11-13)

- [exercism#998](https://github.com/exercism/exercism/issues/998): Fix problem with writing an empty config file under certain circumstances.

## v1.2.2 (2013-11-12)

- [#28](https://github.com/exercism/cli/issues/28): Create exercism directory immediately upon logging in.
- Upgrade to newer version of [codegangsta/cli](https://github.com/codegansta/cli) library, which returns an error from the main Run() function.

## v1.2.1 (2013-11-09)

- [371521f](https://github.com/exercism/cli/commit/371521fd97460aa92269831f10dadd467cb06592): Add support for nested directories under the language track directory allowing us to create idiomatic scala, clojure, and other exercises.

## v1.2.0 (2013-11-07)

- [371521f](https://github.com/exercism/cli/commit/371521fd97460aa92269831f10dadd467cb06592): Consume the new hash of filename => content that the problem API returns.

## v1.1.1 (2013-10-20)

- [371521f](https://github.com/exercism/cli/commit/371521fd97460aa92269831f10dadd467cb06592): Add output when fetching to tell the user where the files where created.

## v1.1.0 (2013-10-24)

- Refactor to extract config package
- Delete stray binary **TODO** we might rewrite history on this one, see [#102](https://github.com/exercism/xgo/issues/102).
- [#22](https://github.com/exercism/cli/pull/22): Display submission url after submitting solution - [@Tonkpils]
- [#21](https://github.com/exercism/cli/pull/21): Add unsubmit command - [@Tonkpils]
- [#20](https://github.com/exercism/cli/pull/20): Add current command - [@Tonkpils]
- Inline refactoring experiment, various cleanup

## v1.0.1 (2013-09-27)

- [#11](https://github.com/exercism/cli/pull/11): Don't require authentication for demo - [@msgehard]
- [#14](https://github.com/exercism/cli/pull/14): Print out fetched assignments - [@Tonkpils]
- [#16](https://github.com/exercism/cli/pull/16): Fix broken submit for relative path names - [@nf]
- Create a separate demo directory if there's no configured exercism directory

## v1.0.0 (2013-09-22)

- [#7](https://github.com/exercism/cli/pull/7): Recognize haskell test files
- [#5](https://github.com/exercism/cli/pull/5): Fix typo - [@simonjefford]
- [#1](https://github.com/exercism/cli/pull/1): Output the location of the config file - [@msgehard]
- Recognize more language test files - [@msgehard]

## v0.0.27.beta (2013-08-25)

All changes by [@msgehard]

- Clean up homedir
- Add dev instructions to README

## v0.0.26.beta (2013-08-24)

All changes by [@msgehard]

- Ensure that ruby gem's config file doesn't get clobbered
- Add cross-compilation
- Set proper User-Agent so server doesn't blow up.
- Implement `submit`
- Implement `demo`
- Implement `peek`
- Expand ~ in config
- Implement `fetch`
- Implement `current`
- Implement `whoami`
- Implement login and logout
- Build on Travis

[@AlexWheeler]: https://github.com/AlexWheeler
[@andrerfcsantos]: https://github.com/andrerfcsantos
[@avegner]: https://github.com/avegner
[@Dparker1990]: https://github.com/Dparker1990
[@John-Goff]: https://github.com/John-Goff
[@LegalizeAdulthood]: https://github.com/LegalizeAdulthood
[@QuLogic]: https://github.com/QuLogic
[@Smarticles101]: https://github.com/Smarticles101
[@Tonkpils]: https://github.com/Tonkpils
[@TrevorBramble]: https://github.com/TrevorBramble
[@alebaffa]: https://github.com/alebaffa
[@ambroff]: https://github.com/ambroff
[@andrewsardone]: https://github.com/andrewsardone
[@anxiousmodernman]: https://github.com/anxiousmodernman
[@beanieboi]: https://github.com/beanieboi
[@blackerby]: https://github.com/blackerby
[@broady]: https://github.com/broady
[@ccnp123]: https://github.com/ccnp123
[@cookrn]: https://github.com/cookrn
[@daveyarwood]: https://github.com/daveyarwood
[@devonestes]: https://github.com/devonestes
[@djquan]: https://github.com/djquan
[@dmmulroy]: https://github.com/dmmulroy
[@dpritchett]: https://github.com/dpritchett
[@eToThePiIPower]: https://github.com/eToThePiIPower
[@ebautistabar]: https://github.com/ebautistabar
[@ekingery]: https://github.com/ekingery
[@elimisteve]: https://github.com/elimisteve
[@ests]: https://github.com/ests
[@farisj]: https://github.com/farisj
[@glebedel]: https://github.com/glebedel
[@harimp]: https://github.com/harimp
[@harugo]: https://github.com/harugo
[@hjljo]: https://github.com/hjljo
[@isbadawi]: https://github.com/isbadawi
[@jbaiter]: https://github.com/jbaiter
[@jdsutherland]: https://github.com/jdsutherland
[@jgsqware]: https://github.com/jgsqware
[@jish]: https://github.com/jish
[@Jrank2013]: https://github.com/Jrank2013
[@jppunnett]: https://github.com/jppunnett
[@katrinleinweber]: https://github.com/katrinleinweber
[@kytrinyx]: https://github.com/kytrinyx
[@larson004]: https://github.com/larson004
[@lcowell]: https://github.com/lcowell
[@manusajith]: https://github.com/manusajith
[@morphatic]: https://github.com/morphatic
[@mrageh]: https://github.com/mrageh
[@msgehard]: https://github.com/msgehard
[@narqo]: https://github.com/narqo
[@neslom]: https://github.com/neslom
[@nf]: https://github.com/nf
[@nilbus]: https://github.com/nilbus
[@nywilken]: https://github.com/nywilken
[@petertseng]: https://github.com/petertseng
[@pminten]: https://github.com/pminten
[@queuebit]: https://github.com/queuebit
[@rcode5]: https://github.com/rcode5
[@rprouse]: https://github.com/rprouse
[@sfairchild]: https://github.com/sfairchild
[@simonjefford]: https://github.com/simonjefford
[@srt32]: https://github.com/srt32
[@xavdid]: https://github.com/xavdid
[@williandrade]: https://github.com/williandrade
[@zabawaba99]: https://github.com/zabawaba99
[@GroophyLifefor]: https://github.com/GroophyLifefor
[@muzimuzhi]: https://github.com/muzimuzhi
[@isberg]: https://github.com/isberg
[@erikschierboom]: https://github.com/erikschierboom
[@sanderploegsma]: https://github.com/sanderploegsma
[@enascimento178]: https://github.com/enascimento178
[@halfdan]: https://github.com/halfdan
[@glennj]: https://github.com/glennj
[@tomasnorre]: https://github.com/tomasnorre
[@ageron]: https://github.com/ageron
[@petrem]: https://github.com/petrem
[@bnandras]: https://github.com/bnandras
[@vaeng]: https://github.com/vaeng
