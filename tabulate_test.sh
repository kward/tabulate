#! /bin/sh
# $Id: erbringen 4489 2009-04-30 04:22:20Z kward $
# vim:et:ft=sh:sts=2:sw=2
#
# Copyright 2008 Kate Ward. All Rights Reserved.
# Author: kate.ward@forestent.com (Kate Ward)
#
# Unit tests for tabulate.

# treat unset variables as an error
set -u

test_getColumnSizes()
{
  rslt=`echo '1 2 3' |getColumnSizes`
  assertEquals '1 1 1' "${rslt}"

  rslt=`getColumnSizes ${test_4col}`
  assertEquals '1 2 3 4' "${rslt}"

  rslt=`FLAGS_include_comments=${FLAGS_TRUE} getColumnSizes ${test_4col}`
  assertEquals '1 4 3 4 4 4 4 9 4' "${rslt}"
}

test_tabulate()
{
  rslt=`echo '1 2 3' |tabulate`
  assertEquals '1 2 3' "${rslt}"

  (FLAGS_justify='left' tabulate ${test_4col} >${test_out})
  diff -u - ${test_out} <<EOF
# this is a test with four columns
1
1 22

1 2  333
1 2  3   4444

1 2  3        # this is a commented line
1 2
1
EOF
  assertTrue $?

  (FLAGS_justify='right' tabulate ${test_4col} >${test_out})
  diff -u - ${test_out} <<EOF
# this is a test with four columns
1
1 22

1  2 333
1  2   3 4444

1  2   3      # this is a commented line
1  2
1
EOF
  assertEquals ${SHUNIT_TRUE} $?
}

test_tabulate_IFS()
{
  # test pipe chars
  echo '1|2|3' |(FLAGS_ifs='|' tabulate >${test_out})
  diff - ${test_out} <<EOF
1 2 3
EOF
  assertEquals ${SHUNIT_TRUE} $?

  # test tab chars
  echo '1   2	3' |(FLAGS_ifs='	' tabulate>${test_out})
  diff - ${test_out} <<EOF
1   2 3
EOF
  assertEquals ${SHUNIT_TRUE} $?
}

test_tabulate_OFS()
{
  echo '1 2 3' |(FLAGS_ofs='|' tabulate ${test_4col} >${test_out})
  diff -u - ${test_out} <<EOF
# this is a test with four columns
1|  |   |
1|22|   |

1|2 |333|
1|2 |3  |4444

1|2 |3  |     # this is a commented line
1|2 |   |
1|  |   |
EOF
}

oneTimeSetUp()
{
  # load the script to test
  . ./tabulate

  # call the setup routine (normally done by main())
  _setup

  test_4col="${shunit_tmpDir}/4col.dat"
  cat >${test_4col} <<EOF
# this is a test with four columns
1
1 22

1 2 333
1 2 3 4444

1 2 3  # this is a commented line
1 2
1
EOF

  test_out="${shunit_tmpDir}/output.dat"
}

oneTimeTearDown()
{
  # call the cleanup routine (normally done by main())
  _cleanup
}

. ${SHLIB_DIR:-/home/kward/lib/sh}/shunit2
