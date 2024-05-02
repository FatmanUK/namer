package main

import (
	"github.com/docopt/docopt-go"
	"strings"
	"os"
	"io"
)

type Control struct {
	args docopt.Opts
	loglevel uint16  // see LL_ constants
	must_continue bool
	view *View
}

func output_stats(re *Control, stats Stats, outfile string) {
	output := stats.String()
	if outfile == "-" {
		for _, line := range strings.Split(output, "\n") {
			re.view.output(line)
		}
	} else {
		re.view.log(LL_DEBUG, "Writing stats to " + outfile)
		f, err := os.Create(outfile)
		if err != nil {
			re.view.log(LL_ERROR, err.Error())
		}
		// defer is useful in case of exception; will always run
		defer f.Close()
		_, err = f.WriteString(output + "\n")
		if err != nil {
			re.view.log(LL_ERROR, err.Error())
		}
	}
}

func analyse_name(re *Control, stats Stats, name string) Stats {
	chunks := make_chunks(name)
	pattern := make_pattern(name)
	prefix := make_prefix(chunks)
	suffix := make_suffix(chunks)
	l1 := make_level1(chunks)
	l2 := make_level2(chunks)
	l3 := make_level3(chunks)
	re.view.log(LL_DEBUG, name)
	re.view.log(LL_DEBUG, StringFromChunkArray(chunks))
	re.view.log(LL_DEBUG, pattern)
	re.view.log(LL_DEBUG, prefix)
	re.view.log(LL_DEBUG, suffix)
	re.view.log(LL_DEBUG, StringFromStringArray(l1))
	re.view.log(LL_DEBUG, StringFromStringArray(l2))
	re.view.log(LL_DEBUG, StringFromStringArray(l3))
	stats.Add(re, chunks, pattern, prefix, suffix, l1, l2, l3)
	return stats
}

func analyse_files(re *Control, infile string, outfile string) Stats {
	re.view.log(LL_DEBUG, "Infile: " + infile)
	re.view.log(LL_DEBUG, "Outfile: " + outfile)
	// TODO: add size limit
	var fbuf []byte
	var err error
	if infile == "-" {
		fbuf, err = io.ReadAll(os.Stdin)
	} else {
		fbuf, err = os.ReadFile(infile)
	}
	if err != nil {
		re.view.log(LL_ERROR, err.Error())
		panic(err)
	}
	lines := strings.Split(string(fbuf), "\n")
	var stats Stats
	stats.Init()
	for key, name := range lines {
		if key == 0 {
			stats.SetString("0_NameAndVersion", name)
		}
		if key == 1 {
			stats.SetString("1_Byline", name)
		}
		if key == 2 {
			stats.SetString("2_Subtitle", name)
		}
		if key == 3 {
			stats.SetString("3_Description", name)
		}
		if len(name) == 0 || key < 4 || name[0] == '#' {
			continue
		}
		//re.view.log(LL_DEBUG, name)
		stats = analyse_name(re, stats, name)
	}
	return stats
}

func analyse_cli(re *Control, names []string) Stats {
	var stats Stats
	for _, name := range names {
		stats = analyse_name(re, stats, name)
	}
	return stats
}

func f_analyse(re *Control) {
	var stats Stats
	re.view.log(LL_DEBUG, "Analysing...")
	names := re.args["<name>"].([]string)
	var outfile string
	outfile = "-"
	if len(names) == 0 {
		infile := re.args["--in"].(string)
		outfile = re.args["--out"].(string)
		stats = analyse_files(re, infile, outfile)
	} else {
		stats = analyse_cli(re, names)
	}
	output_stats(re, stats, outfile)
}

func (re *Control) main_loop() {
	fn_ptrs := make(map[string]func(*Control))
	fn_ptrs["analyse"] = f_analyse
	for fn, _ := range fn_ptrs {
		if re.args[fn].(bool) {
			fn_ptrs[fn](re)
		}
	}
}
