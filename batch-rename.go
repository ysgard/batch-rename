/*
batch-rename.go

A small program to handle batch renames and copies.

Author: Ysgard (Jan Van Uytven) 2013
*/

package main

import (
	"flag"
	"fmt"
	"os"
	//"ioutil" // For ReadDir
	//"path/filepath"
	"regexp"
)

// Why do I have to use "MustCompile"?  Why does a simple Compile raise an error?
var defaultFileMatch = regexp.MustCompile(`^([a-zA-Z0-9\s\._-]+)$`)

var usageString = `
  batch-rename [opts]

    batch-rename will construct a list of all files that match a given regex,
    or all files in the directory, and rename/copy them to a matching file 
    that is modified according to the specified prefix or suffix.

    For example, 'batch-rename' --prefix to_sort_ --regex "/.png$/"' will 
    rename all files matching .png in the current  directory to 
    'to_sort_<oldname>.png'.

    Arguments:
      
      --regex|-x <regex>      
        A regular expression for matching files.  You can use "/<regex>/" or 
        "<regex>", but the double-quotes are necessary.
      
      --prefix|-p <prefix>    
        Renames matching files to have the specified prefix.
      
      --suffix|-s <suffix>    
        Renames matching files to have the specified suffix.
      
      --enumerate|-e <name>   
        Rename matching files to <name>_<num>, where <num> is 
        incremented from 000.

      --target-dir|-t <path>  
        The directory within which we rename/copy.  Default is the current 
        working directory.

      --copy|-c               
        Copy instead of rename.

      --recurse|-r            
        Search for matching files in subdirectories.

      --lowercase|-l          
        Lowercase the final rename.

      --uppercase|-u          
        Uppercase the final rename.

      --dry-run|-n            
        List files, but don't copy/rename
`

var regexArg string
var prefixArg string
var suffixArg string
var enumerateArg string
var targetArg string
var copyArg bool
var recurseArg bool
var lowerArg bool
var upperArg bool
var dryrunArg bool

// Just return the concatenation of the prefix and the filename.
func prefixName(name, prefix string) string {
	return prefix + name
}

// Add a suffix to a filename, being careful to remove and re-add the extension
// on it (if it exists).
// func suffixName(name, suffix string) string {

// }

// Call before parsing flags
func flagInit() {
	// redefine flag's Usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageString)
	}

	const (
		regex_default     = ""
		prefix_default    = ""
		suffix_default    = ""
		enumerate_default = ""
		target_default    = ""
		copy_default      = false
		recurse_default   = false
		lowercase_default = false
		uppercase_default = false
		dryrun_default    = false
	)

	flag.StringVar(&regexArg, "regex", regex_default, "")
	flag.StringVar(&regexArg, "x", regex_default, "")
	flag.StringVar(&prefixArg, "prefix", prefix_default, "")
	flag.StringVar(&prefixArg, "p", prefix_default, "")
	flag.StringVar(&suffixArg, "suffix", suffix_default, "")
	flag.StringVar(&suffixArg, "s", suffix_default, "")
	flag.StringVar(&enumerateArg, "enumerate", enumerate_default, "")
	flag.StringVar(&enumerateArg, "e", enumerate_default, "")
	flag.StringVar(&targetArg, "target-dir", target_default, "")
	flag.StringVar(&targetArg, "t", target_default, "")
	flag.BoolVar(&copyArg, "copy", copy_default, "")
	flag.BoolVar(&copyArg, "c", copy_default, "")
	flag.BoolVar(&recurseArg, "recurse", recurse_default, "")
	flag.BoolVar(&recurseArg, "r", recurse_default, "")
	flag.BoolVar(&lowerArg, "lowercase", lowercase_default, "")
	flag.BoolVar(&lowerArg, "l", lowercase_default, "")
	flag.BoolVar(&upperArg, "uppercase", uppercase_default, "")
	flag.BoolVar(&upperArg, "u", uppercase_default, "")
	flag.BoolVar(&dryrunArg, "dry-run", dryrun_default, "")
	flag.BoolVar(&dryrunArg, "n", dryrun_default, "")

	flag.Parse()

}

func main() {

	// Get the command-line arguments
	flagInit()

	// Print out the values of each argument
	fmt.Fprintf(os.Stdout, "regexArg: "+regexArg+"\n")
	fmt.Fprintf(os.Stdout, "prefixArg: "+prefixArg+"\n")
	fmt.Fprintf(os.Stdout, "suffixArg: "+suffixArg+"\n")
	fmt.Fprintf(os.Stdout, "enumerateArg: "+enumerateArg+"\n")
	fmt.Fprintf(os.Stdout, "targetArg: "+targetArg+"\n")
	fmt.Fprintf(os.Stdout, "copyArg: %b\n", copyArg)
	fmt.Fprintf(os.Stdout, "recurseArg: %b\n", recurseArg)
	fmt.Fprintf(os.Stdout, "lowerArg: %b\n", lowerArg)
	fmt.Fprintf(os.Stdout, "upperArg: %b\n", upperArg)
	fmt.Fprintf(os.Stdout, "dryrunArg: %b\n", dryrunArg)
}
