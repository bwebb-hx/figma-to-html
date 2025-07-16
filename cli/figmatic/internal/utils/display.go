package utils

import (
	"fmt"

	"github.com/fatih/color"
)

type colors struct {
	Lowkey             func(s string) string
	LowkeyPrint        func(s string)
	Bold               func(s string) string
	BoldPrint          func(s string)
	Success            func(s string) string
	SuccessPrint       func(s string)
	SuccessLowkeyPrint func(s string)
	Error              func(s string) string
	ErrorPrint         func(s string)
}

var Colors = colors{
	Lowkey: func(s string) string {
		return lowkeyString(s)
	},
	LowkeyPrint: func(s string) {
		fmt.Println(lowkeyString(s))
	},
	Bold: func(s string) string {
		return boldString(s)
	},
	BoldPrint: func(s string) {
		fmt.Println(boldString(s))
	},
	Success: func(s string) string {
		return successString(s)
	},
	SuccessPrint: func(s string) {
		fmt.Println(successString(s))
	},
	SuccessLowkeyPrint: func(s string) {
		fmt.Println(lowkeySuccessString(s))
	},
	Error: func(s string) string {
		return errorString(s)
	},
	ErrorPrint: func(s string) {
		fmt.Println(errorString(s))
	},
}

func lowkeyString(s string) string {
	return color.New(color.FgHiBlack).SprintFunc()(s)
}

func boldString(s string) string {
	return color.New(color.Bold).SprintFunc()(s)
}

func successString(s string) string {
	return color.New(color.FgHiGreen).SprintFunc()(s)
}

func errorString(s string) string {
	return color.New(color.FgHiRed).SprintFunc()(s)
}

func lowkeySuccessString(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}
