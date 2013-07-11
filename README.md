*WARNING: In Development!*

=batch-rename=

A small program to handle batch renames and copies.

Yes, I know that [GNU program here] exists and does this already.  This was more
an excuse to stretch my Go programming muscles and make something that I'm more
likely to use.


Syntax:

	batch-rename [opts]

    batch-rename will construct a list of all files that match a given regex,
    or all files in the directory, and rename/copy them to a matching file that is
    modified according to the specified prefix or suffix.

    For example, 'batch-rename' -prefix to_sort_ -regex "/.png$/"' will rename all
    files matching .png in the current  directory to 'to_sort_<oldname>.png'.

    Arguments:

      -regex|-x <regex>      A regular expression for matching files.  You can use
                             "/<regex>/" or "<regex>", but the double-quotes are
                             necessary.
      -prefix|-p <prefix>    Renames matching files to have the specified prefix
      -suffix|-s <suffix>    Renames matching files to have the specified suffix
      -enumerate|-e <name>   Rename matching files to <name>_<num>, where <num>
                             is incremented from 000.
      -target-dir|-t <path>  The directory within which we rename/copy.  Default
                             is the current working directory.
      -copy|-c               Copy instead of rename.
      -recurse|-r            Search for matching files in subdirectories
      -lowercase|-l          lowercase the final rename
      -uppercase|-u          uppercase the final rename
      -dry-run|-n            List files, but don't copy/rename
      -force|-f              Override files, if they exist.
