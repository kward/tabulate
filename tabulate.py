#!/usr/bin/python
#
# Copyright 2010 Kate Ward. All Rights Reserved.

"""Tool to columnify textual data.

This script will take a space separated file as input, and will return the file
with all columns nicely lined up.

Sample usage:
$ tabulate some_file
$ cat some_file |tabulate
"""

__author__ = 'kate.ward@forestent.com (Kate Ward)'

import csv
import sys

import gflags as flags

FLAGS = flags.FLAGS

DEFAULT_COMMENT = '#'
DEFAULT_DELIMITER = ' '

flags.DEFINE_string('filename', '', 'file to tabulate', short_name='f')
flags.DEFINE_string('comment', DEFAULT_COMMENT, 'comment character')
flags.DEFINE_string('delimiter', DEFAULT_DELIMITER, 'field separator',
                    short_name='d')
flags.DEFINE_enum('output_format', 'spaced', ['spaced', 'mediawiki', 'mysql'],
                  'output format', short_name='O')
flags.DEFINE_integer('spaces', 1, 'spaces in spaced output', short_name='s')


class ColumnReader(object):
  def __init__(self, fh, delimiter=DEFAULT_DELIMITER,
               comment_char=DEFAULT_COMMENT, **kwds):
    self.comment_char = comment_char
    self.reader = csv.reader(fh, delimiter=delimiter, **kwds)
    self.rows = []
    self.widths = []

  def Read(self):
    for row in self.reader:
      self.rows.append(row)

      if row and row[0].startswith(self.comment_char):
        continue

      row_widths = [len(x) for x in row]
      self.widths = map(lambda x, y: max(x, y), row_widths, self.widths)


class ColumnOutput(object):
  def __init__(self, rows, widths, delimiter=DEFAULT_DELIMITER,
               comment_char=DEFAULT_COMMENT):
    self.comment_char = comment_char
    self.delimiter = delimiter
    self.rows = rows
    self.widths = widths

    self.spaces = []
    for width in self.widths:
      self.spaces.append(' ' * width)

  def Justify(self, row):
    new_row = []
    for column in range(len(self.widths)):
      new_row.append(''.join(map(lambda x, y: x or y, row[column],
                                 self.spaces[column])))
    return new_row


class SpacedOutput(ColumnOutput):
  def Write(self):
    for row in self.rows:
      if row[0].startswith(self.comment_char):
        print self.delimiter.join(row)
        continue

      new_row = self.Justify(row)
      print ' '.join(new_row)


class MediaWikiOutput(ColumnOutput):
  pass


class MySQLOutput(ColumnOutput):
  pass


def main(argv):
  try:
    argv = FLAGS(argv)
  except flags.FlagsError, e:
    print '%s\nUsage: %s ARGS\n%s' % (e, sys.argv[0], FLAGS)
    sys.exit(1)

  output_formats = {'spaced': SpacedOutput,
                    'mediawiki': MediaWikiOutput,
                    'mysql': MySQLOutput}

  try:
    fh = open(FLAGS.filename, 'rU')
  except IOError:
    print 'error: unable to open %s.' % FLAGS.filename
    sys.exit()

  cr = ColumnReader(fh, FLAGS.delimiter, FLAGS.comment)
  cr.Read()
  fh.close()

  output = output_formats[FLAGS.output_format](cr.rows, cr.widths,
                                               FLAGS.delimiter, FLAGS.comment)
  output.Write()


if __name__ == '__main__':
  sys.exit(main(sys.argv))
