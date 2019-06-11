package main

import (
	"echo-fb/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	projectID = getEnv("PROJECT_ID", "project-id")
	port      = getEnv("DERBY_PORT", "9090")
	keyFile   = getEnv("FIREBASE_CREDENTIALS_PATH", "")
)

type ApiResponse struct {
	Name   string `json:"name" `
	Result string `json:"result" `
	Code   int    `json:"code" `
}

type DerbyController struct {
	firebase *service.FirebaseService
}

func NewDerbyController(firebase *service.FirebaseService) *DerbyController {
	return &DerbyController{firebase: firebase}
}

func (dc DerbyController) getProfile(c echo.Context) error {
	name := c.QueryParam("name")
	log.Printf("Looking profile: %s", name)
	found := dc.firebase.GetProfile(name)
	r := &ApiResponse{Name: "profile", Result: fmt.Sprintf("%v", found)}
	if found {
		r.Code = http.StatusOK
		return c.JSON(http.StatusOK, r)
	}
	r.Code = http.StatusNotFound
	return c.JSON(http.StatusAccepted, r)
}

func main() {
	e := echo.New()
	firebaseService := service.NewFirebaseService(projectID, keyFile)

	fb := service.NewFirebaseMiddleWare(firebaseService)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(fb.Process)
	derbyController := NewDerbyController(firebaseService)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))

	e.GET("/profile", derbyController.getProfile)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
