_exercism () {
  local cur prev

  COMPREPLY=()   # Array variable storing the possible completions.
  cur=${COMP_WORDS[COMP_CWORD]}
  prev=${COMP_WORDS[COMP_CWORD-1]}

  commands="configure debug download fetch list open
  restore skip status submit tracks unsubmit
  upgrade help"
  tracks="csharp cpp clojure coffeescript lisp crystal
  dlang ecmascript elixir elm elisp erlang
  fsharp go haskell java javascript kotlin
  lfe lua mips ocaml objective-c php
  plsql perl5 python racket ruby rust scala
  scheme swift typescript bash c ceylon
  coldfusion delphi factor groovy haxe
  idris julia nim perl6 pony prolog
  purescript r sml vbnet powershell"
  config_opts="--dir --host --key --api"
  submit_opts="--test --comment"

  if [ "${#COMP_WORDS[@]}" -eq 2 ]; then
    COMPREPLY=( $( compgen -W  "${commands}" "${cur}" ) )
    return 0
  fi

  if [ "${#COMP_WORDS[@]}" -eq 3 ]; then
    case "${prev}" in
      configure)
        COMPREPLY=( $( compgen -W "${config_opts}" -- "${cur}" ) )
        return 0
        ;;
      fetch)
        COMPREPLY=( $( compgen -W "${tracks}" "${cur}" ) )
        return 0
        ;;
      list)
        COMPREPLY=( $( compgen -W "${tracks}" "${cur}" ) )
        return 0
        ;;
      open)
        COMPREPLY=( $( compgen -W "${tracks}" "${cur}" ) )
        return 0
        ;;
      skip)
        COMPREPLY=( $( compgen -W "${tracks}" "${cur}" ) )
        return 0
        ;;
      status)
        COMPREPLY=( $( compgen -W "${tracks}" "${cur}" ) )
        return 0
        ;;
      submit)
        COMPREPLY=( $( compgen -W "${submit_opts}" -- "${cur}" ) )
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
