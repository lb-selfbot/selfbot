package variables

import "regexp"

// variable represents a variable in a string.
type variable struct {
	name string
	re   *regexp.Regexp
}

var variables = []variable{
	{name: "ver", re: regexp.MustCompile(`\$VERSION`)},
	{name: "date", re: regexp.MustCompile(`\$DATE`)},
	{name: "time", re: regexp.MustCompile(`\$TIME`)},
	{name: "tz", re: regexp.MustCompile(`\$TIME\((.*?)\)`)},
	{name: "uptime", re: regexp.MustCompile(`\$UPTIME`)},
	{name: "activewindow", re: regexp.MustCompile(`\$ACTIVEWINDOW`)},
	{name: "usedcmds", re: regexp.MustCompile(`\$USEDCOMMANDS`)},
	{name: "totalcmds", re: regexp.MustCompile(`\$TOTALCOMMANDS`)},
	{name: "exts", re: regexp.MustCompile(`\$EXTENSIONS`)},
	{name: "servers", re: regexp.MustCompile(`\$SERVERS`)},
	{name: "friends", re: regexp.MustCompile(`\$FRIENDS`)},
	{name: "av", re: regexp.MustCompile(`\$AVATAR`)},
	{name: "os", re: regexp.MustCompile(`\$OS`)},
	{name: "arch", re: regexp.MustCompile(`\$ARCH`)},
	{name: "cpu%", re: regexp.MustCompile(`\$CPUPERCENT`)},
	{name: "cpu", re: regexp.MustCompile(`\$CPU`)},
	{name: "ram", re: regexp.MustCompile(`\$RAMUSAGE`)},
	{name: "ram%", re: regexp.MustCompile(`\$RAMPERCENT`)},
	{name: "disk", re: regexp.MustCompile(`\$DISKUSAGE`)},
	{name: "disk%", re: regexp.MustCompile(`\$DISKPERCENT`)},
	{name: "sname", re: regexp.MustCompile(`\$SONGNAME`)},
	{name: "sartist", re: regexp.MustCompile(`\$SONGARTIST`)},
	{name: "salbum", re: regexp.MustCompile(`\$SONGALBUM`)},
	{name: "sduration", re: regexp.MustCompile(`\$SONGDURATION`)},
	{name: "scover", re: regexp.MustCompile(`\$SONGCOVER`)},
	{name: "slink", re: regexp.MustCompile(`\$SONGLINK`)},
	{name: "weather", re: regexp.MustCompile(`\$WEATHER\((.*?)\)`)},
	{name: "tempc", re: regexp.MustCompile(`\$TEMPC\((.*?)\)`)},
	{name: "tempf", re: regexp.MustCompile(`\$TEMPF\((.*?)\)`)},
}
