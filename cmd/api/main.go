package main

import (
	"os"

	_ "github.com/arseniychernykh/test-backend-1-ChernykhITMO/docs/swagger"
)

// @title Room Booking Service API
// @version 1.0
// @description API for room booking service.
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	os.Exit(run())
}
