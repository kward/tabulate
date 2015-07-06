#! /bin/sh
# $Id: tabulate 4490 2009-04-30 04:25:33Z kward $
# vim:et:ft=sh:sts=2:sw=2
#
# Copyright 2008 Kate Ward. All Rights Reserved.
# Author: kate.ward@forestent.com (Kate Ward)
#
# Script to retabulate a file.
#
# This script will take a space (or other ifs) separated file as input, and
# will return the file with all columns nicely lined up.
#
# Notes:
# - this script requires a temporary files to be created. _cleanup() must be
#   called to clean up afterwards.
#
# Sample usages:
# $ ./tabulate.sh some_file
# $ cat some_file |./tabulate.sh

# treat unset variables as an error
set -u

ARGV0=`basename "$0"`

# load libraries
[ -z "${SHLIB_DIR:-}" ] && SHLIB_DIR='/home/kward/lib/sh'
. ${SHFLAGS_LIB:-${SHLIB_DIR}/shflags}

# define flags
DEFINE_boolean 'debug' false 'enable debug mode' 'd'

DEFINE_string 'justify' 'left' 'justification [left|right]' 'j'
DEFINE_boolean 'include_comments' false 'tabulate the comments too' 'c'

DEFINE_string 'ifs' ' ' 'input field separator' 'i'
DEFINE_string 'ofs' ' ' 'output field separator' 'o'
FLAGS_HELP="USAGE: ${ARGV0} [flags] <file(s)>"

# declare global variables
__tmpDir=''

#------------------------------------------------------------------------------
# functions
#

# This function creates a private temporary directory.
#
# Output:
#   string: the name of a temporary directory
_mktempDir()
{
  # try the standard mktemp function
  ( exec mktemp -dqt tabulate.XXXXXX 2>/dev/null ) && return

  # the standard mktemp didn't work.  doing our own.
  if [ -r '/dev/urandom' ]; then
    _tabulate_random_=`od -vAn -N4 -tx4 </dev/urandom |sed 's/^[^0-9a-f]*//'`
  elif [ -n "${RANDOM:-}" ]; then
    # $RANDOM works
    _tabulate_random_=${RANDOM}${RANDOM}${RANDOM}$$
  else
    # $RANDOM doesn't work
    _tabulate_date_=`date '+%Y%m%d%H%M%S'`
    _tabulate_random_=`expr ${_tabulate_date_} / $$`
  fi

  _tabulate_tmpDir_="${TMPDIR:-/tmp}/tabulate.${_tabulate_random_}"
  ( umask 077 && mkdir "${_tabulate_tmpDir_}" ) || {
    echo 'fatal: could not create temporary directory! exiting' >&2
    exit ${FALSE}
  }

  echo ${_tabulate_tmpDir_}
  unset _tabulate_date_ _tabulate_random_ _tabulate_tmpDir_
}

# This function creates a directory for temporary files.
_setup()
{
  __tmpDir=`_mktempDir`
}

# This function cleans up all temporary files.
_cleanup()
{
  rm -fr ${__tmpDir}
}

# This function exits the script, optionally printing a message
#
# Args:
#   message: string: an error message to be output (optional)
die()
{
  [ $# -ne 0 ] && echo "$@" >&2
  flags_help
  _cleanup
  exit 1
}

# Determines the maximum size of each column in a file.
#
# Args:
#   STDIN: a data file (optional)
#   file_: string: filename of a data file (optional)
# Output:
#   string: a space separated list of the column sizes
getColumnSizes()
{
  if [ $# -eq 0 ]; then file='-'; else file=$1; fi

  awkFile="${__tmpDir}/count_columns.awk"

  cat <<EOF >${awkFile}
BEGIN { FS="${FLAGS_ifs}"; nf=0 }
NF>nf { nf=NF }
{
  for(i=1; i<=NF; i++)
    if (length(\$i)>len[i]) { len[i]=length(\$i) }
}
END {
  OFS=""; ORS=" "
  for(i=1; i<=nf; i++) print len[i]
}
EOF
  catCmd="cat ${file}"
  [ ${FLAGS_include_comments} -eq ${FLAGS_FALSE} ] && \
      catCmd="${catCmd} |sed 's/#.*$//'"
  eval ${catCmd} |awk -f ${awkFile} |sed 's/ *$//'

  unset awkFile file tmpDir
}

# This function retabulates a file.
#
# Args:
#   STDIN: a data file (optional)
#   dataFile_: string: filename of a data file (optional)
# Output:
#   string: the retabulated data file
tabulate()
{
  if [ $# -eq 0 ]; then
    dataFile="${__tmpDir}/datafile.dat"
    cat - >${dataFile}
  else
    dataFile=$1
  fi

  awkFile="${__tmpDir}/tabulate.awk"
  col=1 colstr='' fmtstr=''
  case "${FLAGS_justify}" in
    left) justify='-' ;;
    right) justify='' ;;
  esac
  sizes=`getColumnSizes ${dataFile}`

  if [ ${FLAGS_include_comments} -eq ${FLAGS_FALSE} ]; then
    if [ -n "${sizes}" ]; then
      for size in ${sizes}; do
        colstr="${colstr:+${colstr}, }cols[${col}]"
        fmtstr="${fmtstr:+${fmtstr}${FLAGS_ofs}}%${justify}${size}s"
        col=`expr ${col} + 1`
      done
    else
      # handle empty input correctly
      colstr='cols[0]'
    fi
    cat <<EOF >${awkFile}
{
  comment_idx = index(\$0, "#")  # check for comment
  if (comment_idx) {
    split(\$0, line, "#")
    text = line[1]
  } else {
    text = \$0
  }
  split(text, cols, "${FLAGS_ifs}")
  output = ""
  if (length(text))
    output = sprintf("${fmtstr}", ${colstr})
  if (comment_idx) {  # handle the comment
    comment = substr(\$0, index(\$0, "#"))
    if (length(cols) > 0)  # line is text + comment
      output = sprintf("%s %s", output, comment)
    else  # entire line is a comment
      output = comment
  }
  print output
}
EOF
  else
    for size in ${sizes}; do
      colstr="${colstr:+${colstr}, }\$${col}"
      fmtstr="${fmtstr:+${fmtstr}${FLAGS_ofs}}%${justify}${size}s"
      col=`expr ${col} + 1`
    done
    cat <<EOF >${awkFile}
BEGIN{ FS="${FLAGS_ifs}"; }
{ printf("${fmtstr}\\n", ${colstr}) }
EOF
  fi
  awk -f ${awkFile} ${dataFile} |sed 's/ *$//'

  unset awkFile col colstr dataFile file fmtstr justify sizes
}

#------------------------------------------------------------------------------
# main
#

main()
{
  # sanity checking
  case "${FLAGS_justify}" in
    left|right) ;;
    *) die 'invalid justification' ;;
  esac

  _setup

  # tabulate files
  if [ $# -eq 0 ]; then
    tabulate
  else
    for f in $@; do
      tabulate "${f}"
    done
  fi

  _cleanup
}

# execute main() if this is run in standalone mode (i.e. not in a unit test)
argv0=`echo "${ARGV0}" |sed 's/_test$//;s/_test\.sh$//'`
if [ "${ARGV0}" = "${argv0}" ]; then
  FLAGS "$@" || exit $?
  eval set -- "${FLAGS_ARGV}"
  main "$@"
fi
