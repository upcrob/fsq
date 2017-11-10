# Release Notes

## 1.7.1

* Updated build to use Go 1.9.

## 1.7.0

* Added "matches" keyword for handling regular expression matches.

## 1.6.0

* Added "fsize" keyword for printing the "friendly size" of files.
* Fixed a bug that allowed "stats" keyword to be used in the "where" clause.

## 1.5.0

* Fixed bug where upper case characters in "content contains" searches were not handled properly.
* Added "stats" keyword for aggregate query statistics.

## 1.4.0

* Changed build to use Go 1.8.  Go 1.8 removed the built-in yacc tool, so the external tool must now be downloaded prior to or as part of the build process.
* Optimized file content search.
* Fixed bug for path contains, startswith, and endswith searches on Windows.

## 1.3.0

* Added support for searching multiple root directories.

## 1.2.1

* Fixed file path normalization bug on Windows.

## 1.2.0

* Fixed slash character mixing in output.
* Improved performance due to parallel file evaluation.

## 1.1.0

* Added "ignorecase" keyword for case-insensitive string comparisons.
* Removed the root search directory from results.
* Expression validation fixes.

## 1.0.0

* Defaults to the current directory when the location clause is not present in the expression.
* Added "modified" attribute to display and compare file modification time.
* Fixed directory traversal bug.
* Added optimization for "contents contains" queries.
* Renamed "contents" keyword to "content".
* Updated "size" display and added units suffixes to integer values.
* Added support for "name" and "path" string equality checking.
* Added support for "content" to "startswith" and "endswith" queries.
* Added expression negation.

## 0.0.1

* Initial alpha release.
