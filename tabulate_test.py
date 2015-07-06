#!/usr/bin/python2.6
#
# Copyright 2010 Kate Ward. All Rights Reserved.

"""Unit tests for tabulate.
"""

__author__ = 'kate.ward@forestent.com (Kate Ward)'

import unittest

import tabulate


DATA = ['a bb ccc', 'dd eee ffff', 'ggg hhhh iiiii', 'jj kk ll']

class ColumnReaderTest(unittest.TestCase):
  def testRead(self):
    reader = tabulate.ColumnReader(DATA)
    reader.read()
    self.assertEquals(reader.widths, [3, 4, 5])


class ColumnOutputTest(unittest.TestCase):
  def testJustify(self):
    output = tabulate.ColumnOutput([], [1, 2, 3])
    self.assertEquals(output.justify(['a', 'b', 'c']), ['a', 'b ', 'c  '])



if __name__ == '__main__':
  unittest.main()
