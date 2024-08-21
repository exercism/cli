_exercism () {
  local cur prev

  COMPREPLY=()   # Array variable storing the possible completions.
  cur=${COMP_WORDS[COMP_CWORD]}
  prev=${COMP_WORDS[COMP_CWORD-1]}
  opts="--verbose --timeout"

  commands="configure download open
  submit test troubleshoot upgrade version workspace help"
  config_opts="--show"
  version_opts="--latest"

  if [ "${#COMP_WORDS[@]}" -eq 2 ]; then
    case "${cur}" in
      -*)
        COMPREPLY=( $( compgen -W  "${opts}" -- "${cur}" ) )
        return 0
        ;;
      *)
        COMPREPLY=( $( compgen -W  "${commands}" "${cur}" ) )
        return 0
        ;;
    esac
  fi

  if [ "${#COMP_WORDS[@]}" -eq 3 ]; then
    case "${prev}" in
      configure)
        COMPREPLY=( $( compgen -W "${config_opts}" -- "${cur}" ) )
        return 0
        ;;
      version)
        COMPREPLY=( $( compgen -W "${version_opts}" -- "${cur}" ) )
        return 0
        ;;
      help)
        COMPREPLY=( $( compgen -W "${commands}" "${cur}" ) )
        return 0
        ;;
      *)
        return 0
        ;;
    esac
  fi

  return 0
}

complete -o bashdefault -o default -o nospace -F _exercism exercism 2>/dev/null \
	|| complete -o default -o nospace -F _exercism exercism
