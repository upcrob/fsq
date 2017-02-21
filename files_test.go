package main

import (
    "testing"
    "os"
)

func TestUpdateSearchCaseSensitive(t *testing.T) {
    path := "testdata/sub1/Test3.Txt"
    searchStrings := make([]SearchString, 0, 6)
    searchStrings = append(searchStrings,
        SearchString{"So", true},
        SearchString{"SO", true},
        SearchString{"xt", true},
        SearchString{"XT", true},
        SearchString{"me", true},
        SearchString{"ME", true})
    fs := newFileSearch(searchStrings, path)
    info, _ := os.Stat(path)
    updateSearch(fs, info)

    if len(fs.contains) != 3 {
        t.Error("expected 3 entries in fs.contains, found", len(fs.contains))
    }

    if !searchStringExists(fs.contains, SearchString{"So", true}) {
        t.Error("'So' expected in contains")
    }
    if !searchStringExists(fs.contains, SearchString{"me", true}) {
        t.Error("'me' expected in contains")
    }
    if !searchStringExists(fs.contains, SearchString{"xt", true}) {
        t.Error("'xt' expected in contains")
    }
}

func TestUpdateSearchCaseInsensitive(t *testing.T) {
    path := "testdata/sub1/Test3.Txt"
    searchStrings := make([]SearchString, 0, 3)
    searchStrings = append(searchStrings,
        SearchString{"SO", false},
        SearchString{"so", false},
        SearchString{"xT", false},
        SearchString{"XT", false},
        SearchString{"mE", false},
        SearchString{"ME", false})
    fs := newFileSearch(searchStrings, path)
    info, _ := os.Stat(path)
    updateSearch(fs, info)

    if len(fs.contains) != 6 {
        t.Error("expected 6 entries in fs.contains, found", len(fs.contains))
    }

    if !searchStringExists(fs.contains, SearchString{"so", false}) {
        t.Error("'so' expected in contains")
    }
    if !searchStringExists(fs.contains, SearchString{"mE", false}) {
        t.Error("'mE' expected in contains")
    }
    if !searchStringExists(fs.contains, SearchString{"ME", false}) {
        t.Error("'ME' expected in contains")
    }
    if !searchStringExists(fs.contains, SearchString{"xT", false}) {
        t.Error("'xT' expected in contains")
    }
}

func TestUpdateSearchAtBlockBorder(t *testing.T) {
    path := "testdata/sub4/overlap1.txt"
    searchStrings := make([]SearchString, 0, 1)
    searchStrings = append(searchStrings, SearchString{"helloworldthisisatest", true})
    fs := newFileSearch(searchStrings, path)
    info, _ := os.Stat(path)

    if !fileContainsString(fs, info, "helloworldthisisatest", true) {
        t.Error("'helloworldthisisatest' expected in contains")
    }

    if info.Size() != 1307 {
        t.Error("warning: expected test file to be 1213 characters long, update this test if necessary")
    }
}

func TestUpdateSearchAtMultipleBlockBorders(t *testing.T) {
    path := "testdata/sub4/overlap2.txt"
    searchStrings := make([]SearchString, 0, 1)
    searchStrings = append(searchStrings,
        SearchString{"helloworldthisisatest", true},
        SearchString{"thisisoverlapnumbertwo", true},
        SearchString{"thisisoverlapnumberthree", true})
    fs := newFileSearch(searchStrings, path)
    info, _ := os.Stat(path)

    if !fileContainsString(fs, info, "helloworldthisisatest", true) {
        t.Error("'helloworldthisisatest' expected in contains")
    }
    if !fileContainsString(fs, info, "thisisoverlapnumbertwo", true) {
        t.Error("'thisisoverlapnumbertwo' expected in contains")
    }
    if !fileContainsString(fs, info, "thisisoverlapnumberthree", true) {
        t.Error("'thisisoverlapnumberthree' expected in contains")
    }

    if info.Size() != 3179 {
        t.Error("warning: expected test file to be 1213 characters long, update this test if necessary")
    }
}
