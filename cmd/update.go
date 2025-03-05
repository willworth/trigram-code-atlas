package cmd

import "flag"

func Update(fs *flag.FlagSet, args []string) {
	fs.Parse(args)
}
