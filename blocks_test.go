package main

import (
	"testing"
)

var names = []string{"lisa", "anna", "bob", "Боъ", "hannah", "mark", "paul", "testing", "alfræd", "ἐπιστήμη", "adam", "james", "hammandaturian", "пример"}
var patts = []string{"CVCV", "VCV", "CVC", "CVC", "CVCVC", "CVC", "CVC", "CVCVC", "VCVC", "VCVCVCV", "VCVC", "CVCVC", "CVCVCVCVCVC", "CVCVC"}
var prefixes = []string{"li", "ann", "bo", "Бо", "ha", "ma", "pau", "te", "alfr", "ἐπ", "ad", "ja", "ha", "при"}
var suffixes = []string{"sa", "nna", "ob", "оъ", "ah", "ark", "aul", "ing", "æd", "μη", "am", "es", "ian", "ер"}

func TestChunks(t *testing.T) {
	for _, name := range names {
		t.Log("Name: " + name)
		chunks := make_chunks(name)
		t.Log(chunks)
	}
}

func TestPatterns(t *testing.T) {
	for key, name := range names {
		patt := make_pattern(name)
		want_patt := patts[key]
		//t.Log(patt)
		if patt == want_patt {
			t.Logf("Pattern %s for '%s' is correct", patt, name)
		} else {
			t.Errorf("Pattern for %s is %s, expected %s", name, patt, want_patt)
		}
	}
}

func TestPrefixes(t *testing.T) {
	var prefix string
	for key, name := range names {
		chunks := make_chunks(name)
		prefix = make_prefix(chunks)
		want_prefix := prefixes[key]
		if prefix == want_prefix {
			t.Logf("Prefix %s for '%s' is correct", prefix, name)
		} else {
			t.Errorf("Prefix for %s is %s, expected %s", name, prefix, want_prefix)
		}
	}
}

func TestSuffixes(t *testing.T) {
	var suffix string
	for key, name := range names {
		chunks := make_chunks(name)
		suffix = make_suffix(chunks)
		want_suffix := suffixes[key]
		if suffix == want_suffix {
			t.Logf("Suffix %s for '%s' is correct", suffix, name)
		} else {
			t.Errorf("Suffix for %s is %s, expected %s", name, suffix, want_suffix)
		}
	}
}

func TestLevel1(t *testing.T) { // adjacent blocks
	for _, name := range names {
		chunks := make_chunks(name)
		t.Log(make_level1(chunks))
	}
}

func TestLevel2(t *testing.T) { // semi-adjacent vowels
	for _, name := range names {
		chunks := make_chunks(name)
		t.Log(make_level2(chunks))
	}
}

func TestLevel3(t *testing.T) { // semi-adjacent consonants
	for _, name := range names {
		chunks := make_chunks(name)
		t.Log(make_level3(chunks))
	}
}
