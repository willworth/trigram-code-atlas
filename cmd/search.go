package cmd

import "flag"

func Search(fs *flag.FlagSet, args []string) {
	fs.Parse(args)
}
