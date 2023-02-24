package ui

import "runtime"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

func GetWorkerColor(team, number int) string {
	if team == 1 && number == 1 {
		return Blue
	}
	if team == 1 && number == 2 {
		return Cyan
	}
	if team == 2 && number == 1 {
		return Red
	}
	if team == 2 && number == 2 {
		return Yellow
	}
	if team == 3 && number == 1 {
		return Purple
	}
	if team == 3 && number == 2 {
		return Green
	}
	return White
}
