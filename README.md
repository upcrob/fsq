# fsq

[![Build Status](https://travis-ci.org/upcrob/fsq.png)](https://travis-ci.org/upcrob/fsq)

The `fsq` ('file system query' - pronounced, 'fisk') utility is a tool for doing ad-hoc queries against a file system using a SQL-like expression language.  This is useful for finding files that match certain criteria without writing a one-off script to do so.

## Installation

Download the binary for your platform and add it to your command line path.

[![Download](https://api.bintray.com/packages/upcrob/generic/fsq/images/download.svg)](https://bintray.com/upcrob/generic/fsq/_latestVersion#files)

## Usage

### Query Structure

`fsq` takes a single argument: the expression.  This expression is composed of the following parts:

	<attribute list> in <locations> where <conditions>

### Example Queries

To recursively find all files under the '/data' directory that start with the characters 'hello' and are larger than 5 mb, the following query could be used:

	fsq "name in '/data' where name startswith 'hello' and size > 5m"

If the location (in the above case, '/data') is omitted, `fsq` will default to the current directory:

	fsq "name where name startswith 'hello' and size > 5m"

Multiple locations can be specified as well:

	fsq "name in '/opt', '/media' where size > 5m"

The attribute list specifies which attributes are printed to standard out by `fsq`.  In the above case, this is just the filename ('name').  The following example will print both the path to the file and the size (in bytes):

	fsq "path,size in '/opt' where size > 5m"

### Supported Attributes

* `name`
* `path`
* `size`
* `fsize` (can be used in the attribute list, but cannot be queried)
* `content` (content can be queried, but cannot be added to the attribute list for printing)
* `modified` (format: 'MM/DD/YYYY' or 'MM/DD/YYYY hh:mm:ss')
* `stats` (can be used in the attribute list, but cannot be queried)

### Supported Conditional Operators

* `<`
* `<=`
* `>`
* `>=`
* `=`
* `!=`
* `startswith`
* `endswith`
* `isdir` (this operator does not take any arguments)
* `isfile` (this operator does not take any arguments)
* `contains`
* `ignorecase` (must be followed by '=', '!=', 'startswith', 'endswith', or 'contains')
* `matches` (regular expression matching)

### Logic Operators

Parentheses as well as the logical operators *or*, *and*, and *not* can be used to group conditions.  For example:

	fsq "name in '.' where name startswith 'hello' or (isdir and not name startswith 'world')"

### Size Qualifiers

The following size qualifiers can be appended to integer values to indicate non-default units.  These are especially useful when specifying file sizes in expressions.  If no size qualifier is appended to an integer, `fsq` compares the value in bytes.

* k - Kilobytes
* m - Megabytes
* g - Gigabytes

For example, to find all files greater than 10 kilobytes and less than 1 megabyte:

	fsq "path where size > 10k and size < 1m"

## Building

The `go` compiler is required to build `fsq`.  If you have `make` installed, `fsq` can be installed with:

	make install

Otherwise, the following commands will need to be run while in the `fsq` directory:

	go get golang.org/x/tools/cmd/goyacc
	go install golang.org/x/tools/cmd/goyacc
	goyacc parser.y
	go install
