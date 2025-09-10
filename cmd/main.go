// package main - entrypoint for client rpc
package main

import (
	"log"

	"github.com/fatih/color"
)

// main выводит в лог сообщение "Hello, World!" зелёного цвета.
// Отображение цвета зависит от поддержки терминала.
func main() {
	log.Println(color.GreenString("Hello, World!"))
}
