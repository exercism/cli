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


#### Oh my Zsh

If you are using the popular [oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh) framework to manage your zsh plugins, move the file `exercism_completion.zsh` into `~/.oh-my-zsh/custom`.

### Fish

Completions must go in the user defined `$fish_complete_path`. By default, this is `~/.config/fish/completions`

    mv ../shell/exercism.fish ~/.config/fish/exercism.fish
