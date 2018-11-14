# Configure
complete -f -c exercism -n "__fish_use_subcommand" -a "configure" -d "Writes config values to a JSON file."
complete -f -c exercism -n "__fish_seen_subcommand_from configure" -s t -l token -d "Set token"
complete -f -c exercism -n "__fish_seen_subcommand_from configure" -s w -l workspace -d "Set workspace"
complete -f -c exercism -n "__fish_seen_subcommand_from configure" -s a -l api -d "set API base url"
complete -f -c exercism -n "__fish_seen_subcommand_from configure" -s s -l show -d "show settings"

# Download
complete -f -c exercism -n "__fish_use_subcommand" -a "download" -d "Downloads and saves a specified submission into the local system"
complete -f -c exercism -n "__fish_seen_subcommand_from download" -s e -l exercise -d "the exercise slug"
complete -f -c exercism -n "__fish_seen_subcommand_from download" -s h -l help -d "help for download"
complete -f -c exercism -n "__fish_seen_subcommand_from download" -s T -l team -d "the team slug"
complete -f -c exercism -n "__fish_seen_subcommand_from download" -s t -l track -d "the track ID"
complete -f -c exercism -n "__fish_seen_subcommand_from download" -s u -l uuid -d "the solution UUID"

# Help
complete -f -c exercism -n "__fish_use_subcommand" -a "help" -d "Shows a list of commands or help for one command"
complete -f -c exercism -n "__fish_seen_subcommand_from help" -a "configure download help open submit troubleshoot upgrade version workspace"

# Open
complete -f -c exercism -n "__fish_use_subcommand" -a "open" -d "Opens a browser to exercism.io for the specified submission."
complete -f -c exercism -n "__fish_seen_subcommand_from open" -s h -l help -d "help for open"

# Submit
complete -f -c exercism -n "__fish_use_subcommand" -a "submit" -d "Submits a new iteration to a problem on exercism.io."
complete -f -c exercism -n "__fish_seen_subcommand_from submit" -s h -l help -d "help for submit"

# Troubleshoot
complete -f -c exercism -n "__fish_use_subcommand" -a "troubleshoot" -d "Outputs useful debug information."
complete -f -c exercism -n "__fish_seen_subcommand_from troubleshoot" -s f -l full-api-key -d "display full API key (censored by default)"
complete -f -c exercism -n "__fish_seen_subcommand_from troubleshoot" -s h -l help -d "help for troubleshoot"

# Upgrade
complete -f -c exercism -n "__fish_use_subcommand" -a "upgrade" -d "Upgrades to the latest available version."
complete -f -c exercism -n "__fish_seen_subcommand_from help" -s h -l help -d "help for help"

# Version
complete -f -c exercism -n "__fish_use_subcommand" -a "version" -d "Outputs version information."
complete -f -c exercism -n "__fish_seen_subcommand_from version" -s l -l latest -d "check latest available version"
complete -f -c exercism -n "__fish_seen_subcommand_from version" -s h -l help -d "help for version"

# Workspace
complete -f -c exercism -n "__fish_use_subcommand" -a "workspace" -d "Outputs the root directory for Exercism exercises."
complete -f -c exercism -n "__fish_seen_subcommand_from workspace" -s h -l help -d "help for workspace"

# Options
complete -f -c exercism -s h -l help -d "show help"
complete -f -c exercism -l timeout -a "(seq 0 100 100000)" -d "override default HTTP timeout"
complete -f -c exercism -s v -l verbose -d "turn on verbose logging"
