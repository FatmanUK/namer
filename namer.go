package main

import (
	"os"
	"github.com/docopt/docopt-go"
	"fatgo/docopt_helpers"
)

// Default values may be overridden by build system
var build_mode = "Debug"
var app_name = "namer"
var app_version = "unbuilt version"

func main() {
	loglevel := LL_INFO
	if build_mode == "Debug" {
		loglevel = LL_DEBUG
	}
	view := &View{loglevel: loglevel}
	view.begin()

	uses := []string{}
	uses = append(uses, "analyse <name>...")
	uses = append(uses, "analyse [--in=<infile> [--out=<outfile>]]")

	opts := make(map[string]string)
	opts["--in=<infile>"] = "Input file [default: -]"
	opts["--out=<outfile>"] = "Output file [default: -]"

	args, err := docopt.ParseArgs(
		docopt_helpers.BuildUsageString(uses, opts),
		nil,
		(app_name + " " + app_version))

	if err != nil {
		view.log(LL_ERROR, err.Error())
		os.Exit(1)
	}

	app := Control{
		args: args,
		must_continue: true,
		view: view}
	app.main_loop()
	view.end()
}
