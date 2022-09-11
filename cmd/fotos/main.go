package main

import (
	fotos "github.com/DerZombiiie/fotos/fotos"

	_ "github.com/DerZombiiie/fotos/input_console"
	_ "github.com/DerZombiiie/fotos/ui_console"
	_ "github.com/DerZombiiie/fotos/ui_webview"
	_ "github.com/DerZombiiie/fotos/webcam"
	_ "github.com/DerZombiiie/fotos/webinterface"
)

func main() {
	fotos.Run()
}
