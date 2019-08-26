package main

import (
	"github.com/fatih/color"
)


func infoMsg(s string, newline bool)  {
	info := color.New(color.FgHiMagenta)
	if newline {
		info.Println(s)
	} else {
		info.Print(s)
	}
}

func dangerMsg(s string, newline bool)  {
	info := color.New(color.FgHiRed)
	if newline {
		info.Println(s)
	} else {
		info.Print(s)
	}
}

func successMsg(s string, newline bool) {
	info := color.New(color.FgHiYellow)
	if newline {
		info.Println(s)
	} else {
		info.Print(s)
	}
}

func warningMsg(s string, newline bool)  {
	info := color.New(color.FgHiCyan)
	if newline {
		info.Println(s)
	} else {
		info.Print(s)
	}
}