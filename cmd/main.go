// package main - entrypoint for client rpc
package main

import (
	"log"

	"github.com/fatih/color"
)

// main — точка входа программы. Выводит зелёно окрашенную строку "Hello, World!" в стандартный лог.
func main() {
	log.Println(color.GreenString("Hello, World!"))
}
