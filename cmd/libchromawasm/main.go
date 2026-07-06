//go:build wasm

// Package main is an experimental WASM library intended for TinyGO.
package main

import (
	"syscall/js"

	"github.com/alecthomas/chroma/v3/internal/playground"
)

func main() {
	// Register the highlight function with the JavaScript global object
	js.Global().Set("highlight", js.FuncOf(highlight))

	// Keep the program running
	select {}
}

// Highlight source code using Chroma.
//
// Equivalent to the JS function:
//
//	function highlight(source, lexer, styleName, classes)
//
// If the "lexer" is unknown, this will attempt to autodetect the language type.
func highlight(this js.Value, args []js.Value) any {
	source := args[0].String()
	lexer := args[1].String()
	styleName := args[2].String()
	classes := args[3].Bool()

	res, err := playground.Highlight(source, lexer, styleName, classes)
	if err != nil {
		panic(err)
	}
	return js.ValueOf(map[string]any{
		"html":       res.HTML,
		"language":   res.Language,
		"background": res.Background,
	})
}
