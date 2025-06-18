package main

import (
	"crypto/rand"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// TemplateData holds data for rendering the index template
type TemplateData struct {
	PoWChallenge   string
	FullURL        string
	Path           string
	ErrorMessage   string
	SuccessMessage string
}

var challengeStore ChallengeStorage

func main() {
	// Initialize challenge storage
	challengeStore = NewChallengeStorage()
	defer challengeStore.Close()

	e := echo.New()
	e.GET("/", serveHome)
	e.POST("/shorten.html", handleShorten)
	e.GET("/shorten.html", serveHome)
	e.GET("/*", handleRedirectOrStatic)
	e.Start(":8080")
}

// generateRandomString generates a random string of specified length using [a-zA-Z0-9] characters
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomByte[0]%byte(len(charset))]
	}

	return string(result), nil
}

// generateRandomPath generates a random path string of specified length using [a-z0-9] characters
func generateRandomPath(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)

	for i := range result {
		randomByte := make([]byte, 1)
		_, err := rand.Read(randomByte)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomByte[0]%byte(len(charset))]
	}

	return string(result), nil
}

// renderIndexWithData renders the index.html template with the provided data
func renderIndexWithData(c echo.Context, data TemplateData) error {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))
	c.Response().Header().Set("Content-Type", "text/html")
	return tmpl.Execute(c.Response().Writer, data)
}

// generateNewChallenge generates a new unique challenge and stores it
func generateNewChallenge() (string, error) {
	var challenge string
	var err error
	maxTries := 1000
	for i := 0; i < maxTries; i++ {
		challenge, err = generateRandomString(200)
		if err != nil {
			return "", err
		}
		if _, exists, err := challengeStore.Get(challenge); err != nil {
			return "", err
		} else if exists {
			continue // If challenge already exists, generate a new one
		}
		break // Exit loop if a unique challenge is generated
	}

	// Store the challenge as unsolved initially
	err = challengeStore.Store(challenge, false)
	if err != nil {
		return "", err
	}

	return challenge, nil
}

func serveHome(c echo.Context) error {
	challenge, err := generateNewChallenge()
	if err != nil {
		log.Printf("Failed to generate challenge: %v", err)
		return c.String(http.StatusInternalServerError, "error generating challenge")
	}

	acceptHeader := c.Request().Header.Get("Accept")

	// Check if the client accepts WAP content
	if strings.Contains(acceptHeader, "text/vnd.wap.wml") {
		// Wap device detected, redirect to Bevelgacom WAP site
		return c.Redirect(http.StatusMovedPermanently, "https://wap.bevelgacom.be")
	}

	data := TemplateData{
		PoWChallenge:   challenge,
		FullURL:        "",
		Path:           "",
		ErrorMessage:   "",
		SuccessMessage: "",
	}

	return renderIndexWithData(c, data)
}

func verifyChallenge(c echo.Context) (bool, string, error) {
	challenge := c.FormValue("pow_challenge")
	solution := c.FormValue("pow_solution")

	if challenge == "" || solution == "" {
		return false, "challenge and solution are required", nil
	}

	// Convert solution to integer
	solutionInt, err := strconv.Atoi(solution)
	if err != nil {
		return false, "invalid solution format", nil
	}

	// check if the challenge exists and is unsolved
	solved, exists, err := challengeStore.Get(challenge)
	if err != nil {
		log.Printf("Failed to retrieve challenge: %v", err)
		return false, "", err // Return actual error for internal server errors
	}
	if !exists {
		return false, "challenge not found", nil
	}
	if solved {
		return false, "challenge already solved", nil
	}

	// Verify the proof of work
	if !VerifyProofOfWork(challenge, solutionInt, 4) {
		return false, "invalid proof of work", nil
	}

	// Mark the challenge as solved
	err = challengeStore.Store(challenge, true)
	if err != nil {
		log.Printf("Failed to mark challenge as solved: %v", err)
		return false, "", err // Return actual error for internal server errors
	}

	return true, "", nil
}

// serve404 serves the appropriate 404 page based on the Accept header
func serve404(c echo.Context) error {
	acceptHeader := c.Request().Header.Get("Accept")

	// Check if the client accepts WAP content
	if strings.Contains(acceptHeader, "text/vnd.wap.wml") {
		// Serve WAP 404 page
		c.Response().Header().Set("Content-Type", "text/vnd.wap.wml")
		return c.File("templates/404.wml")
	}

	// Serve regular 404 response
	return c.String(http.StatusNotFound, "404 - Not Found")
}

func handleShorten(c echo.Context) error {
	fullURL := c.FormValue("fullURL")
	path := c.FormValue("path")

	// Helper function to render error with form values preserved
	renderError := func(errorMsg string) error {
		challenge, err := generateNewChallenge()
		if err != nil {
			log.Printf("Failed to generate new challenge: %v", err)
			return c.String(http.StatusInternalServerError, "error generating new challenge")
		}

		data := TemplateData{
			PoWChallenge:   challenge,
			FullURL:        fullURL,
			Path:           path,
			ErrorMessage:   errorMsg,
			SuccessMessage: "",
		}

		return renderIndexWithData(c, data)
	}

	ok, errorMsg, err := verifyChallenge(c)
	if err != nil {
		// Internal server error
		return c.String(http.StatusInternalServerError, "internal server error")
	}
	if !ok {
		return renderError(errorMsg)
	}

	// If the challenge is verified, proceed with URL shortening
	if path == "" {
		// Generate a random 5-character path with [a-z0-9]
		var err error
		maxTries := 1000
		for i := 0; i < maxTries; i++ {
			path, err = generateRandomPath(5)
			if err != nil {
				log.Printf("Failed to generate random path: %v", err)
				return c.String(http.StatusInternalServerError, "error generating random path")
			}

			// Check if the generated path already exists
			_, exists, err := challengeStore.GetURL(path)
			if err != nil {
				log.Printf("Failed to check if path exists: %v", err)
				return c.String(http.StatusInternalServerError, "error checking path availability")
			}
			if !exists {
				break // Path is available, use it
			}
			// If path exists, generate a new one
		}
	}

	// check if the path is [a-zA-Z0-9_-] and not too long
	// it may not conain any characters other than [a-zA-Z0-9_-]
	if len(path) < 1 || len(path) > 50 {
		return renderError("invalid path length, must be between 1 and 50 characters")
	}
	// check if path contains only valid characters with regexp

	regexpPath := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !regexpPath.MatchString(path) {
		return renderError("invalid path format, must contain only [a-zA-Z0-9_-]")
	}

	if fullURL == "" {
		return renderError("full URL is required")
	}
	if len(fullURL) > 200 {
		return renderError("full URL is too long, must be less than 200 characters")
	}

	// check if the path already exists
	_, exists, err := challengeStore.GetURL(path)
	if err != nil {
		log.Printf("Failed to retrieve URL mapping: %v", err)
		return c.String(http.StatusInternalServerError, "error retrieving URL mapping")
	}
	if exists {
		return renderError("path already exists")
	}

	// Check if the full URL is valid, if http:// is not provided, add it
	if !isValidURL(fullURL) {
		if !isValidURL("http://" + fullURL) {
			return renderError("invalid full URL format")
		}
		fullURL = "http://" + fullURL
	}

	// Store the URL mapping
	err = challengeStore.StoreURL(path, fullURL)
	if err != nil {
		log.Printf("Failed to store URL mapping: %v", err)
		return c.String(http.StatusInternalServerError, "error storing URL mapping")
	}

	// Generate a new challenge for the success page
	challenge, err := generateNewChallenge()
	if err != nil {
		log.Printf("Failed to generate new challenge: %v", err)
		return c.String(http.StatusInternalServerError, "error generating new challenge")
	}

	// Render success page with shortened URL
	data := TemplateData{
		PoWChallenge:   challenge,
		FullURL:        "",
		Path:           "",
		ErrorMessage:   "",
		SuccessMessage: "URL shortened successfully! Your short URL is: wap.fyi/" + path,
	}

	return renderIndexWithData(c, data)
}

// handleRedirectOrStatic handles requests that could be shortened URLs or static files
func handleRedirectOrStatic(c echo.Context) error {
	path := c.Param("*")

	// Remove leading slash if present
	path = strings.TrimPrefix(path, "/")

	// If path is empty, serve index.html
	if path == "" {
		return serveHome(c)
	}

	// Check if path matches shortened URL pattern [a-zA-Z0-9_-]
	regexpPath := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if regexpPath.MatchString(path) && len(path) >= 1 && len(path) <= 20 {
		// Try to get the full URL from storage
		fullURL, exists, err := challengeStore.GetURL(path)
		if err != nil {
			log.Printf("Failed to retrieve URL mapping for %s: %v", path, err)
			// Fall through to static file serving
		} else if exists {
			// Redirect to the full URL
			return c.Redirect(http.StatusMovedPermanently, fullURL)
		}
	}

	// Security: Protect against path traversal attacks
	// Clean the path and ensure it doesn't contain path traversal sequences
	cleanPath := filepath.Clean(path)
	if strings.Contains(cleanPath, "..") || strings.HasPrefix(cleanPath, "/") {
		log.Printf("Path traversal attempt detected: %s", path)
		return c.String(http.StatusNotFound, "404 - Not Found")
	}

	// Construct the full file path within the templates directory
	fullPath := filepath.Join("templates", cleanPath)

	// Ensure the resolved path is still within the templates directory
	absTemplatesDir, err := filepath.Abs("templates")
	if err != nil {
		log.Printf("Failed to get absolute path for templates directory: %v", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		log.Printf("Failed to get absolute path for %s: %v", fullPath, err)
		return c.String(http.StatusNotFound, "404 - Not Found")
	}

	// Check if the resolved absolute path is within the templates directory
	if !strings.HasPrefix(absFullPath, absTemplatesDir+string(filepath.Separator)) && absFullPath != absTemplatesDir {
		log.Printf("Path traversal attempt detected (absolute path check): %s -> %s", path, absFullPath)
		return c.String(http.StatusNotFound, "404 - Not Found")
	}

	// Check if file exists before serving
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("File not found: %s", fullPath)
		return serve404(c)
	} else if err != nil {
		log.Printf("Error checking file %s: %v", fullPath, err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// Serve the static file
	return c.File(fullPath)
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Check if scheme and host are present
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}

	// Only allow http and https schemes
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Basic host validation - must contain at least one dot for domain
	if !strings.Contains(parsedURL.Host, ".") {
		return false
	}

	return true
}
