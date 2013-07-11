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
	"path/filepath"
	"regexp"
	//"github.com/jteeuwen/go-pkg-optarg" optarg
	"errors"
	"io" // Needed for file copy
	"strings"
)

// Why do I have to use "MustCompile"?  Why does a simple Compile raise an error?
var defaultFileMatch = regexp.MustCompile(`^([a-zA-Z0-9\s\._-]+)$`)

var usageString = `
  batch-rename (-p <prefix>|-s <suffix>|-e <name>) [-x <regex>] [-t <target_dir>] -[crn] -[l|u]

    batch-rename will construct a list of all files that match a given regex,
    or all files in the directory, and rename/copy them to a matching file 
    that is modified according to the specified prefix or suffix.

    For example, 'batch-rename' -prefix to_sort_ -regex "/.png$/"' will 
    rename all files matching .png in the current  directory to 
    'to_sort_<oldname>.png'.

    Arguments:
      
      -regex|-x <regex>      
        A regular expression for matching files.  You can use "/<regex>/" or 
        "<regex>", but the double-quotes are necessary.
      
      -prefix|-p <prefix>    
        Renames matching files to have the specified prefix.
      
      -suffix|-s <suffix>    
        Renames matching files to have the specified suffix.
      
      -enumerate|-e <name>   
        Rename matching files to <name>_<num>, where <num> is 
        incremented from 000.

      -target-dir|-t <path>  
        The directory within which we rename/copy.  Default is the current 
        working directory.

      -copy|-c               
        Copy instead of rename.

      -recurse|-r            
        Search for matching files in subdirectories.

      -lowercase|-l          
        Lowercase the final rename. (Can't be used with '-u')

      -uppercase|-u          
        Uppercase the final rename. (Can't be used with '-l')

      -dry-run|-n            
        List files, but don't copy/rename

      -force|-for 
        If the file exists, overwrite.  Default is to not copy/rename
        if the target file already exists.
      }
`

var regexArg string
var fileRegex *regexp.Regexp
var prefixArg string
var suffixArg string
var enumerateArg string
var targetArg string
var copyArg bool
var recurseArg bool
var lowerArg bool
var upperArg bool
var dryrunArg bool
var forceArg bool

// Just return the concatenation of the prefix and the filename.
func prefixName(name, prefix string) string {
	return prefix + name
}

// Add a suffix to a filename, being careful to remove and re-add the extension
// on it (if it exists).
func suffixName(name, suffix string) string {
	ext := filepath.Ext(name)
	raw_base := strings.TrimSuffix(filepath.Base(name), ext)
	return raw_base + suffix + ext
}

// enumerated files take the form <name>_<dddd>.<ext>
func enumerateName(name, newname string, count int) string {
	ext := filepath.Ext(name)
	return fmt.Sprintf("%s_%04d.%s", newname, count, ext)
}

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
		force_default     = false
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
	flag.BoolVar(&forceArg, "f", force_default, "")
	flag.BoolVar(&forceArg, "force", force_default, "")

	flag.Parse()

}

func usage(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")
	fmt.Fprintf(os.Stderr, usageString)
	os.Exit(0)
}

// Build a list of all files that match the regex, and then walk
// them and rename them as we go.
func processFiles() (int, error) {

	var targetDir string

	if targetArg != "" {
		targetDir = targetArg // If target specified, use it
	} else {
		var err error
		targetDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not determine current directory!  Please specify target directory with -t.")
			return 0, err
		}
	}

	// Make sure the target directory is a valid one
	if dinfo, err := os.Lstat(targetDir); err != nil || dinfo.IsDir() == false {
		err = errors.New(fmt.Sprintf("Target directory %s is not a directory?\n", targetDir))
		fmt.Fprintf(os.Stderr, err.Error())
		return 0, err
	}

	fmt.Fprintf(os.Stderr, "Target directory is : %s\n", targetDir)

	// Process each file & return the results
	fCount, err := processDir(targetDir)
	return fCount, err

}

func processDir(dirname string) (int, error) {

	var count, counter int // count is how many files were successfully counted, counter is for enum
	var return_msg string
	var return_err error
	dp, open_err := os.Open(dirname)
	if open_err != nil {
		fmt.Fprintf(os.Stderr, open_err.Error())
		return 0, open_err
	}
	defer dp.Close()

	files, read_err := dp.Readdir(0)
	if read_err != nil {
		return_msg += read_err.Error() + "\n"
	}
	for _, f := range files {
		if f.IsDir() == true && recurseArg == true {
			dir_count, dir_err := processDir(filepath.Join(dirname, f.Name()))
			if dir_err != nil {
				return_msg += dir_err.Error() + "\n"
			}
			count += dir_count
		} else if f.IsDir() == true && recurseArg == false {
			continue
		} else {
			file_count, file_err := processFile(filepath.Join(dirname, f.Name()), counter)
			if file_err != nil {
				return_msg += file_err.Error() + "\n"
			}
			count += file_count
			counter++

		}
	}
	if return_msg == "" {
		return_err = nil
	} else {
		return_err = errors.New(return_msg)
	}
	return count, return_err
}

// Rename/Copy a file based on global arguments.
func processFile(name string, currentCount int) (int, error) {

	base := filepath.Base(name)
	path := filepath.Dir(name)
	if fileRegex != nil {
		// Check to see whether file matches, exit if it doesn't
		if fileRegex.MatchString(base) == false {
			return 0, nil
		}
	}

	// Order is - enumerate, suffix, prefix.
	var newName, outputMsg string
	ext := filepath.Ext(name)
	newName = strings.TrimSuffix(filepath.Base(name), ext)
	if enumerateArg != "" {
		newName = fmt.Sprintf("%s_%04d", enumerateArg, currentCount)
	}
	if suffixArg != "" {
		newName = newName + suffixArg
	}
	if prefixArg != "" {
		newName = prefixArg + newName
	}
	newName = filepath.Join(path, newName+ext)

	// Check to see if the renamed file exists.  If it does, what
	// we do next depends on whether or not the force flag is applied
	_, fileInfoErr := os.Lstat(newName)
	if fileInfoErr == nil && forceArg == false { // file does exist & force is not set
		fmt.Fprintf(os.Stderr, "File %s already exists, not copying/renaming. Use --force to override\n", newName)
		return 0, nil
	}

	// Copy/Rename the file
	if copyArg == true {
		outputMsg = fmt.Sprintf("Copying %s to %s...\n", name, newName)
	} else {
		outputMsg = fmt.Sprintf("Renaming %s to %s...\n", name, newName)
	}
	fmt.Fprintf(os.Stdout, outputMsg)

	if dryrunArg == true {
		return 0, nil
	} else {
		if copyArg == true {

			copyFile, openFileErr := os.Create(newName)
			if openFileErr != nil {
				return 0, openFileErr
			}
			srcFile, openFileErr := os.Open(name)
			if openFileErr != nil {
				return 0, openFileErr
			}
			io.Copy(copyFile, srcFile)
			return 1, nil

		} else {

			if rename_err := os.Rename(name, newName); rename_err != nil {
				return 0, rename_err
			}
			return 1, nil
		}
	}

}

func main() {

	// Get the command-line arguments
	flagInit()

	// Sanity check - no upper with lower
	if lowerArg == true && upperArg == true {
		usage("Cannot combine -u (uppercase) and -l (lowercase) flags")
	}

	// We require at least one of 'prefix', 'suffix' or 'enumerate', otherwise
	// we don't know how to rename.
	if prefixArg == "" && suffixArg == "" && enumerateArg == "" {
		usage("Specify one of -p <prefix>, -s <suffix> or -e <enumerate>")
	}

	// Some people like to bracket a regex with '/'.  Strip these out, if found
	regexArg = strings.TrimLeft(strings.TrimRight(regexArg, "/"), "/")

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

	// Compile the passed regex, if any
	if regexArg != "" {
		var err error
		fileRegex, err = regexp.Compile(regexArg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid regex: %s\n", regexArg)
			fmt.Fprintf(os.Stderr, "Cannot continue, quitting...")
			os.Exit(1)
		}
	}

	fCount, err := processFiles()

	if err == nil {
		fmt.Fprintf(os.Stdout, "\nOperation complete: %d files renamed/copied\n", fCount)
	} else {
		fmt.Fprintf(os.Stdout, "\nOperation could not be completed: %s\n", err.Error())
	}

}
