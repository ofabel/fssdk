package app

import "fmt"

func (r *Runtime) Print(v ...any) {
	if !r.quiet {
		fmt.Print(v...)
	}
}

func (r *Runtime) Println(v ...any) {
	if !r.quiet {
		fmt.Println(v...)
	}
}

func (r *Runtime) Printf(format string, v ...any) {
	if !r.quiet {
		fmt.Printf(format, v...)
	}
}
