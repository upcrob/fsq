# fsq

[![Build Status](https://travis-ci.org/upcrob/fsq.png)](https://travis-ci.org/upcrob/fsq)

The `fsq` ('file system query' - pronounced, 'fisk') utility is a tool for doing ad-hoc queries against a file system using a SQL-like expression language.  This is useful for finding files that match certain criteria without writing a one-off script to do so.

## Installation

Download the binary for your platform and add it to your command line path.

[![Download](https://api.bintray.com/packages/upcrob/generic/fsq/images/download.svg)](https://bintray.com/upcrob/generic/fsq/_latestVersion)

## Usage

### Example Query

`fsq` is designed to be quickly invoked from the command line.  For example, in order to find all files under the current directory that start with the characters 'hello' and are larger than 5 mb, the following query could be used:

	fsq "name in '.' where name startswith 'hello' and size > 5m"

## Query Structure

Notice that `fsq` takes a single argument: the expression.  This expression is composed of the following parts:

	<attribute list> in <root directories> where <conditions>

The attribute list specifies which attributes are printed to standard out by `fsq`.  In the above case, this is just the filename ('name').

The root directory tells `fsq` where to start searching in the file system.  Every directory under the root will be searched recursively for files matching the given conditions.  In the above case, it starts searching at the current directory ('.').  If you want to search multiple directories, separate them with commas.  Note that if the `in <location>` part of the expression is left out, it will default to the current directory.

The set of conditions tells `fsq` what files it should print out as matches.  In the above case, it looks for a name that *startswith* the string 'hello' and has a *size* on disk greater than 5 megabytes.

### Supported Attributes

* name
* path
* size
* content
* modified

### Supported Conditional Operators

* <
* <=
* >
* >=
* =
* !=
* startswith
* endswith
* isdir (this operator does not take any arguments)
* isfile (this operator does not take any arguments)
* contains
* ignorecase (must be followed by '=', '!=', 'startswith', 'endswith', or 'contains')

### Logic Operators

Parentheses as well as the logical operators *or*, *and*, and *not* can be used to group conditions.  For example:

	fsq "name in '.' where name startswith 'hello' or (isdir and name startswith 'world')"

### Size Qualifiers

The following size qualifiers can be appended to integer values to indicate non-default units.  These are especially useful when specifying file sizes in expressions.  If no size qualifier is appended to an integer, `fsq` compares the value in bytes.

* k - Kilo
* m - Mega
* g - Giga

## Building

The `go` compiler is required to build `fsq`.  If you have `make` installed, `fsq` can be installed with:

	make install

Otherwise, the following commands will need to be run while in the `fsq` directory:

	go tool yacc parser.y
	go install
