package main

import (
	"embed"
	_ "embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed assets
var iconAssets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "wails-react-ts-template",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},

		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
			ActivationPolicy: application.ActivationPolicyAccessory,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	mainWindow := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Name:      "main",
		Title:     "Window 1",
		Frameless: true,

		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		URL: "/",
	})

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	setupSystemTray(app, mainWindow)

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}

// setupSystemTray creates and configures the system tray for the application
func setupSystemTray(app *application.App, mainWindow *application.WebviewWindow) {
	// Create a new system tray
	systray := app.NewSystemTray()

	// Read icon data
	iconBytes, err := iconAssets.ReadFile("assets/icon.png")
	if err != nil {
		log.Printf("Failed to load system tray icon: %v", err)
		return
	}

	// Set icon and label
	systray.SetIcon(iconBytes)
	//systray.SetLabel("Wails App")

	// Create system tray menu
	menu := app.NewMenu()

	menu.Add("Show Window").OnClick(func(ctx *application.Context) {
		mainWindow.Show()
		mainWindow.Focus()
	})

	menu.AddSeparator()

	menu.Add("Quit").OnClick(func(ctx *application.Context) {
		app.Quit()
	})

	systray.SetMenu(menu)

	// Attach the window to the system tray
	systray.AttachWindow(mainWindow)

	// Set window offset and debounce time
	systray.WindowOffset(10)
	systray.WindowDebounce(200 * time.Millisecond)

}
