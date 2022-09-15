package main

import (
	fotos "github.com/DerZombiiie/fotos/fotos"

	_ "github.com/DerZombiiie/fotos/cloud_sync"
	_ "github.com/DerZombiiie/fotos/input_console"
	_ "github.com/DerZombiiie/fotos/local_webinterface"
	_ "github.com/DerZombiiie/fotos/ui_console"
	_ "github.com/DerZombiiie/fotos/ui_webview"
	_ "github.com/DerZombiiie/fotos/webcam"
)

func main() {
	fotos.Run()
}
