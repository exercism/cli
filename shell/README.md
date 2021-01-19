## Executable
Unpack the archive relevant to your machine and place in $PATH

## Shell Completion Scripts

### Bash

    mkdir -p ~/.config/exercism
    mv ../shell/exercism_completion.bash ~/.config/exercism/exercism_completion.bash

Load the completion in your `.bashrc`, `.bash_profile` or `.profile` by
adding the following snippet:

    if [ -f ~/.config/exercism/exercism_completion.bash ]; then
      source ~/.config/exercism/exercism_completion.bash
    fi

### Zsh

Load up the completion by placing the `exercism_completion.zsh` somewhere on
your `$fpath` as `_exercism`. For example:

    mkdir -p ~/.zsh/functions
    mv ../shell/exercism_completion.zsh ~/.zsh/functions/_exercism

and then add the directory to your `$fpath` in your `.zshrc`, `.zsh_profile` or
`.profile` before running `compinit`:

    export fpath=(~/.zsh/functions $fpath)
    autoload -U compinit && compinit


#### Oh My Zsh

If you are using the popular [Oh My Zsh][oh-my-zsh] framework to manage your
zsh plugins, you need to move the file `exercism_completion.zsh` to a new
custom plugin:

[oh-my-zsh]: https://github.com/ohmyzsh/ohmyzsh

    mkdir -p $ZSH_CUSTOM/plugins/exercism
    cp exercism_completion.zsh $ZSH_CUSTOM/plugins/exercism/_exercism

Then edit the file `~/.zshrc` to include `exercism` in the list of plugins.
Completions will be activated the next time you open a new shell. If the
completions do not work, you should update Oh My Zsh to the latest version with
`omz update`. Oh My Zsh now checks whether the plugin list has changed (more
accurately, `$fpath`) and resets the `zcompdump` file.

### Fish

Completions must go in the user defined `$fish_complete_path`. By default, this is `~/.config/fish/completions`

    mv ../shell/exercism.fish ~/.config/fish/exercism.fish
