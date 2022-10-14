package main

import (
	fotos "github.com/DerZombiiie/fotos/fotos"
	// configure your OWN version:
	//
	//_ "github.com/DerZombiiie/fotos/cloud_sync"
	// cloud sync syncronizes to a "github.com/DerZombiiie/fotos/cmd/server" - server
	// configure using:
	// - "cloud.addr" a the base url used (incl. protocol)
	//     e.g: "https://example.com/"
	// - "cloud.token" the upload token used to identify this instance
	//     (see cmd/server/help.md)
	//
	//_ "github.com/DerZombiiie/fotos/input_signal"
	// input signal listens on any SIGUSR1 signal send to this process
	// can be used togather with "github.com/DerZombiiie/fotos/cmd/fotos_signal"
	// (input_signal has to be executed as root if this process is running under root)
	//
	//_ "github.com/DerZombiiie/fotos/input_gpio"
	// input gpio listens on a gpio pin (pin: 26)
	// is broken currently as it breaks the camera
	// use a external keyboard & i3 & input_signal instead
	//
	//_ "github.com/DerZombiiie/fotos/local_webinterface"
	// local webinterface presents a local album instance
	// configure using:
	// - "http.album.addr" for a listening adress
	//     e.g: "localhost:8080" o. ":8080"
	// - "error.http.500" message presented on err 500
	// - "error.http.404" message presented on err 404
	//
	//_ "github.com/DerZombiiie/fotos/ui_console"
	// ui console provies a ui on the terminal you started fotos from
	//
	//_ "github.com/DerZombiiie/fotos/ui_webview"
	// provides a webview based interface when using webcam
	//
	//_ "github.com/DerZombiiie/fotos/libcamera"
	// libcamera communicates with any running libcamera-still
	// processes. use "DerZombiiie/cmd/fotos/libcamera.sh" to
	// start a correct instance (pwd = cmd/fotos/)
	//
	//_ "github.com/DerZombiiie/fotos/ui_glfw"
	// TODO: ui_glfw is broken and needs fixing
	//
	//_ "github.com/DerZombiiie/fotos/ui_neopixel"
	// ui neopixel provides the coundown on a neopixel stripe
	// attached to pin 33 // gpio 18 // PWM1
)

func main() {
	fotos.Run()
}
