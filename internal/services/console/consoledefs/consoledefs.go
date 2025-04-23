package consoledefs

// Set contains the ANSI and HTML data for setting the console.
type Set struct {
	ANSI string
	HTML string
}

// Append contains the ANSI and HTML data for appending to the console.
type Append struct {
	ANSI string
	HTML string
}
