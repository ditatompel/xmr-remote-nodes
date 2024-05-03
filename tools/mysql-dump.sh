#!/bin/sh
# Dump local dev database structure and required data from specific tables.
# ------------------------------------------------------------------------------
# shellcheck disable=SC2046 # Ignore SC2046: Quote this to prevent word splitting.

SD=$(dirname "$(readlink -f -- "$0")")
cd "$SD" || exit 1 && cd ".." || exit 1

## Structure only dump
mariadb-dump --no-data --skip-comments xmr_nodes | \
    sed 's/ AUTO_INCREMENT=[0-9]*//g' >            \
    "./tools/resources/database/structure.sql"

# vim: set ts=4 sw=4 tw=0 et ft=sh:
