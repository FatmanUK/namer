package main

import (
	"strings"
	"encoding/json"
	"fmt"
)

// TODO: oh yes, the 'h' thing. Makes 'h' a special bunny. Ignore it in matching (eg. match "cho" not "ho").
// TODO: add options '--yvowel' and '--wvowel' to analyse function?
// note, need to encode these settings in the chapter so there's no surprise when we find one or the other

type Stats struct {
	Strings map[string]string
	Chunks map[string]int
	Patterns map[string]int
	Prefixes map[string]int
	Suffixes map[string]int
	Level1 map[string]int
	Level2 map[string]int
	Level3 map[string]int
}

func (stats *Stats) Init() {
	stats.Strings = make(map[string]string)
	stats.Chunks = make(map[string]int)
	stats.Patterns = make(map[string]int)
	stats.Prefixes = make(map[string]int)
	stats.Suffixes = make(map[string]int)
	stats.Level1 = make(map[string]int)
	stats.Level2 = make(map[string]int)
	stats.Level3 = make(map[string]int)
}

func (stats *Stats) SetString(index string, name string) {
	stats.Strings[index] = name
}

func (stats *Stats) Add(re *Control, chunks []Chunk, pattern string, prefix string, suffix string, l1 []string, l2 []string, l3 []string) {
	re.view.log(LL_DEBUG, "Adding...")
	for _, chunk := range chunks {
		_, has_chunk := stats.Chunks[chunk.String()]
		if has_chunk {
			stats.Chunks[chunk.String()]++
		} else {
			stats.Chunks[chunk.String()] = 1
		}
	}
	_, has_pattern := stats.Patterns[pattern]
	if has_pattern {
		stats.Patterns[pattern]++
	} else {
		stats.Patterns[pattern] = 1
	}
	_, has_prefix := stats.Prefixes[prefix]
	if has_prefix {
		stats.Prefixes[prefix]++
	} else {
		stats.Prefixes[prefix] = 1
	}
	_, has_suffix := stats.Suffixes[suffix]
	if has_suffix {
		stats.Suffixes[suffix]++
	} else {
		stats.Suffixes[suffix] = 1
	}

	for _, chunk := range l1 {
		_, has_chunk := stats.Level1[chunk]
		if has_chunk {
			stats.Level1[chunk]++
		} else {
			stats.Level1[chunk] = 1
		}
	}
	for _, chunk := range l2 {
		_, has_chunk := stats.Level2[chunk]
		if has_chunk {
			stats.Level2[chunk]++
		} else {
			stats.Level2[chunk] = 1
		}
	}
	for _, chunk := range l3 {
		_, has_chunk := stats.Level3[chunk]
		if has_chunk {
			stats.Level3[chunk]++
		} else {
			stats.Level3[chunk] = 1
		}
	}
}

func (re *Stats) String() string {
	bytes, err := json.Marshal(re)
	if err != nil {
		fmt.Println(err)
	}
	return string(bytes)
}

type Chunk struct {
	Content string
}

type Stringify interface {
	String() string
}

func (re *Chunk) String() string {
	return re.Content
}

func StringFromChunkArray(chunks []Chunk) string {
	rv := ""
	for _, v := range chunks {
		trim_rv := strings.Trim(rv, " ")
		rv = strings.Join([]string{trim_rv, v.String()}, " ")
	}
	return "[" + rv + "]"
}

func IsVowel(r rune) bool {
	// is this all of them? aren't there any more? really?
	// note some Cyrillic ones are missing
	var vowels string
	is_welsh := false
	is_cyrillic := false
	is_greek := false
	is_german := true
	is_y_vowel := false
	is_extra := true
	if true {
		vowels += "aeiou"
	}
	if is_y_vowel {
		vowels += "yyẙỳỵỷỹŷÿý"
	}
	if is_welsh {
		vowels += "w"
	}
	if is_cyrillic {
		vowels += "ɯeɜи"
	}
	if is_greek {
		vowels += "ɤɛἐἐιήη"
	}
	if is_german {
		vowels += "äöü"
	}
	if is_extra {
		vowels += "ɨʉuɪʏʊøɘɵoəœɞʌɔæɐaɶɑɒaàáâãåāăą"
		vowels += "ǻȁȃạảấầẩẫậắằẳẵặḁæǽ"
		vowels += "eȅȇḕḗḙḛḝẹẻẽếềểễệēĕėęěèéêëе"
		vowels += "iȉȋḭḯỉịĩīĭįi̇ìíîïĳ"
		vowels += "oоœøǿȍȏṍṏṑṓọỏốồổỗộớờởỡợōòóŏőôõ"
		vowels += "uũūŭůűųùúûȕȗṳṵṷṹṻụủứừửữự"
	}
	return strings.ContainsRune(vowels, r)
}

func get_letter_kind(r rune) rune {
	if IsVowel(r) {
		return 'V'
	} else {
		return 'C'
	}
}

// is level 1 where we take only the nearest char of the cons chunk?
// might be 1 and 3
func make_level1(chunks []Chunk) []string {
	first_iter := true
	var rv []string
	var oldch Chunk
	for _, chunk := range chunks {
		if !first_iter {
			rv = append(rv, oldch.Content + chunk.Content)
		} else {
			first_iter = false
		}
		oldch = chunk
	}
	return rv
}

func add_match(rv []string, arr [2]Chunk, want_vowels bool) []string {
	// using range breaks the string by rune, not by byte
	// naive implementation (retrieve by string index) gets the byte :/
	var firstch rune
	for _, ch := range arr[0].Content {
		firstch = ch
		break
	}
	if want_vowels == IsVowel(firstch) {
		rv = append(rv, strings.Join([]string{arr[0].Content, arr[1].Content}, "*"))
	}
	return rv
}

func make_level2_and_level3(want_vowels bool, chunks []Chunk) []string {
	first_iter := true
	second_iter := false
	var rv []string
	var oldch Chunk
	var grandoldch Chunk
	for _, chunk := range chunks {
		if !first_iter && !second_iter {
			arr := [2]Chunk{grandoldch, chunk}
			rv = add_match(rv, arr, want_vowels)
		} else {
			if !first_iter { // second
				second_iter = false
			} else { // first
				first_iter = false
				second_iter = true
			}
		}
		grandoldch = oldch
		oldch = chunk
	}
	return rv
}

func make_level2(chunks []Chunk) []string {
	return make_level2_and_level3(true, chunks)
}

// or is it level 3? I think level 3.
func make_level3(chunks []Chunk) []string {
	// in fact I'm so sure it's level 3 that I'm implementing it
	l3 := make_level2_and_level3(false, chunks)
	// using range breaks the strings by rune, not by byte
	// naive algorithm gets the byte which yields wrong behaviour :/
	for key, fit := range l3 {
		arr := strings.Split(fit, "*")
		if len(arr[0]) > 1 {
			var lastch rune
			for _, ch := range arr[0] {
				lastch = ch
			}
			arr[0] = string(lastch)
		}
		if len(arr[1]) > 1 {
			var firstch rune
			for _, ch := range arr[1] {
				firstch = ch
				break
			}
			arr[1] = string(firstch)
		}
		l3[key] = strings.Join(arr, "*")
	}
	return l3
}

func make_prefix(chunks []Chunk) string {
	prefix := chunks[0].Content
	if (len(chunks) > 1) {
		prefix += chunks[1].Content
	}
	return prefix
}

func make_suffix(chunks []Chunk) string {
	l := len(chunks)
	suffix := chunks[l - 1].Content
	if (l > 1) {
		suffix = chunks[l - 2].Content + suffix
	}
	return suffix
}

func make_chunks(name string) []Chunk {
	first_iter := true
	var vc bool
	var ovc bool
	acc := ""
	var rv []Chunk
	for _, ch := range name {
		ovc = vc
		vc = IsVowel(ch)
		if !first_iter {
			// if vc and ovc differ, we must have VC or CV
			if vc != ovc {
				rv = append(rv, Chunk{acc})
				acc = string(ch)
			} else { // else, VV or CC
				//
				acc += string(ch)
			}
		} else {
			first_iter = false
			acc = string(ch)
		}
	}
	rv = append(rv, Chunk{acc})
	return rv
}

func make_pattern(name string) string {
	first_iter := true
	var vc bool
	var ovc bool
	acc := ""
	for _, ch := range name {
		ovc = vc
		vc = IsVowel(ch)
		if !first_iter {
			if vc != ovc {
				acc += string(get_letter_kind(ch))
			}
		} else {
			first_iter = false
			acc += string(get_letter_kind(ch))
		}
	}
	return acc
}
