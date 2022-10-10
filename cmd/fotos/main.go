package main

import (
	fotos "github.com/DerZombiiie/fotos/fotos"

	_ "github.com/DerZombiiie/fotos/cloud_sync"
	_ "github.com/DerZombiiie/fotos/input_console"
	_ "github.com/DerZombiiie/fotos/input_gpio"
	_ "github.com/DerZombiiie/fotos/local_webinterface"
	_ "github.com/DerZombiiie/fotos/ui_console"
	//	_ "github.com/DerZombiiie/fotos/ui_webview"
	//_ "github.com/DerZombiiie/fotos/libcamera_singleshot"
	_ "github.com/DerZombiiie/fotos/libcamera"
	// _ "github.com/DerZombiiie/fotos/ui_glfw"
	_ "github.com/DerZombiiie/fotos/ui_neopixel"
	// _ "github.com/DerZombiiie/fotos/ui_webview"
)

func main() {
	fotos.Run()
}
