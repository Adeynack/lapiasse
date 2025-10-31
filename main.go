/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log/slog"

	"adeynack.net/lapiasse/cmd"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	cmd.Execute()
}
