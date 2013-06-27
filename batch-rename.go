/* batch-rename.go

A small program to handle batch renames and copies.

Syntax:

	batch-rename [opts]

    batch-rename will construct a list of all files that match a given regex,
    or all files in the directory, and rename/copy them to a matching file that is
    modified according to the specified prefix or suffix.

    For example, 'batch-rename' --prefix to_sort_ --regex "/.png$/"' will rename all
    files matching .png in the current  directory to 'to_sort_<oldname>.png'.

    Arguments:

      --regex|-x <regex>      A regular expression for matching files.  You can use
                             "/<regex>/" or "<regex>", but the double-quotes are
                             necessary.
      --prefix|-p <prefix>    Renames matching files to have the specified prefix
      --suffix|-s <suffix>    Renames matching files to have the specified suffix
      --enumerate|-e <name>   Rename matching files to <name>_<num>, where <num>
                             is incremented from 000.
      --target-dir|-t <path>  The directory within which we rename/copy.  Default
                             is the current working directory.
      --copy|-c               Copy instead of rename.
      --recurse|-r            Search for matching files in subdirectories
      --lowercase|-l          lowercase the final rename
      --uppercase|-u          uppercase the final rename
      --dry-run|-n            List files, but don't copy/rename

Author: Ysgard (Jan Van Uytven) 2013
*/

package main

import (
	"fmt"
	"os"
	//"ioutil" // For ReadDir
	"regexp"
)

var defaultFileMatch = regexp.MustCompile("^([a-zA-Z0-9\\s\\._-]+)$")

func prefixName(name, prefix string) string {
	return prefix + name
}

func main() {
	if defaultFileMatch == nil {
		fmt.Fprintf(os.Stderr, "Invalid regexp object!\n")
		os.Exit(1)
	}
	fname1 := "This is a valid filename.txt"
	fname2 := "This@@@@isnotvalid.txt"
	if defaultFileMatch.MatchString(fname1) {
		fmt.Fprintf(os.Stdout, "Properly matched!\n")
	} else {
		fmt.Fprintf(os.Stdout, "Should match, but doesn't!\n")
	}
	if !defaultFileMatch.MatchString(fname2) {
		fmt.Fprintf(os.Stdout, "Properly doesn't match!\n")
	} else {
		fmt.Fprintf(os.Stdout, "Matches, but shouldn't!\n")
	}

}
