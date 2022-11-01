#!/bin/bash

# Script to configure your Bash shell

__main() {
  if ! [[ "$PATH" =~ "%ROOT_DIRECTORY%/bin" ]]; then
    export PATH="%ROOT_DIRECTORY%/bin:$PATH"
  fi
  if ! [[ "$MANPATH" =~ "%ROOT_DIRECTORY%" ]]; then
    export MANPATH="%ROOT_DIRECTORY%/man:$MANPATH"
  fi
  if ! [[ "$LB_LIBRARY_PATH" =~ "%ROOT_DIRECTORY%" ]]; then
    export LB_LIBRARY_PATH="%ROOT_DIRECTORY%/lib:$LB_LIBRARY_PATH"
  fi

  for f in "%ROOT_DIRECTORY%/completions.d/bash/"*; do
    source "$f"
  done
}

__main
unset -f __main
