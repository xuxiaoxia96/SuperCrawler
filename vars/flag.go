package vars

import "flag"

var (
	Mode             = flag.String("m", "update", "Update / All")
	Target			 = flag.String("t", "", "target to crawl, list of register func, use ',' to split, like 'aa,bb,cc'")
	Version     	 = flag.Bool("v", false, "show version and exit.")
)