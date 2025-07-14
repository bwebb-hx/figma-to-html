package utils

import (
	"fmt"

	"github.com/fatih/color"
)

type colors struct {
	Lowkey       func(s string) string
	LowkeyPrint  func(s string)
	Bold         func(s string) string
	BoldPrint    func(s string)
	Success      func(s string) string
	SuccessPrint func(s string)
	Error        func(s string) string
	ErrorPrint   func(s string)
}

var Colors = colors{
	Lowkey: func(s string) string {
		return lowkey(s)
	},
	LowkeyPrint: func(s string) {
		fmt.Println(lowkey(s))
	},
	Bold: func(s string) string {
		return bold(s)
	},
	BoldPrint: func(s string) {
		fmt.Println(bold(s))
	},
	Success: func(s string) string {
		return success(s)
	},
	SuccessPrint: func(s string) {
		fmt.Println(success(s))
	},
	Error: func(s string) string {
		return error(s)
	},
	ErrorPrint: func(s string) {
		fmt.Println(error(s))
	},
}

func lowkey(s string) string {
	return color.New(color.FgHiBlack).SprintFunc()(s)
}

func bold(s string) string {
	return color.New(color.Bold).SprintFunc()(s)
}

func success(s string) string {
	return color.New(color.FgHiGreen).SprintFunc()(s)
}

func error(s string) string {
	return color.New(color.FgHiRed).SprintFunc()(s)
}
