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

    mkdir -p ~/.config/exercism
    mv ../shell/exercism_completion.zsh ~/.config/exercism/exercism_completion.zsh

Load up the completion in your `.zshrc`, `.zsh_profile` or `.profile` by adding
the following snippet

    if [ -f ~/.config/exercism/exercism_completion.zsh ]; then
      source ~/.config/exercism/exercism_completion.zsh
    fi

**Note:** If you are using the popular [oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh) framework to manage your zsh plugins, you don't need to add the above snippet, all you need to do is create a file `exercism_completion.zsh` inside the `~/.oh-my-zsh/custom`.
