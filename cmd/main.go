package main

import (
	"go-concurrent/concurrent"
)

func main() {
	callReadFileCountWord()
}

func callReadFileCountWord() {
	concurrent.ReadFileCountWord()
}
