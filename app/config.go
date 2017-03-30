package app

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
)

import (
	"github.com/jessevdk/go-flags"
	"github.com/kelseyhightower/envconfig"
)

// AuthPlzConfig configuration structure
type AuthPlzConfig struct {
	Name                  string `short:"n" long:"name" description:"User friendly service name"`
	Address               string `short:"a" long:"address" description:"Set server address"`
	Port                  string `short:"p" long:"port" description:"Set server port"`
	Database              string `short:"d" long:"database" description:"Database connection string"`
	CookieSecret          string `long:"cookie-secret" description:"32-byte base64 encoded secret for cookie / session storage" default-mask:"-"`
	TokenSecret           string `long:"token-secret" description:"32-byte base64 encoded secret for token use" default-mask:"-"`
	TLSCert               string `short:"c" long:"tls-cert" description:"TLS Certificate file"`
	TLSKey                string `short:"k" long:"tls-key" description:"TLS Key File"`
	NoTLS                 bool   `long:"disable-tls" description:"Disable TLS for testing or reverse proxying"`
	StaticDir             string `short:"s" long:"static-dir" description:"Directory to load static assets from"`
	TemplateDir           string `short:"t" long:"template-dir" description:"Directory to load templates from"`
	MinimumPasswordLength int
}

// GenerateSecret Helper to generate a default secret to use
func GenerateSecret(len int) (string, error) {
	data := make([]byte, len)
	n, err := rand.Read(data)
	if err != nil {
		return "", err
	}
	if n != len {
		return "", errors.New("Config: RNG failed")
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

// DefaultConfig Generate default configuration
func DefaultConfig() (*AuthPlzConfig, error) {
	var c AuthPlzConfig

	c.Name = "AuthPlz"
	c.Address = "localhost"
	c.Port = "9000"
	c.Database = "host=localhost user=postgres dbname=postgres sslmode=disable password=postgres"

	// Certificate files in environment
	c.TLSCert = "server.pem"
	c.TLSKey = "server.key"
	c.NoTLS = false
	c.StaticDir = "../authplz-ui"
	c.TemplateDir = "./templates"

	c.MinimumPasswordLength = 12

	var err error

	c.CookieSecret, err = GenerateSecret(32)
	if err != nil {
		return nil, err
	}
	c.TokenSecret, err = GenerateSecret(32)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetConfig fetch the server configuration
// This parses environmental variables, command line flags, and in future
// will handle file based loading of configurations.
func GetConfig() *AuthPlzConfig {
	// Fetch default configuration
	c, err := DefaultConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Parse config structure through environment
	err = envconfig.Process("authplz", c)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Override environment with command line args
	_, err = flags.Parse(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Decode secrets to strings
	tokenSecret, err := base64.URLEncoding.DecodeString(c.TokenSecret)
	if err != nil {
		log.Println(err)
		log.Panic("Error decoding token secret")
	}

	cookieSecret, err := base64.URLEncoding.DecodeString(c.CookieSecret)
	if err != nil {
		log.Println(err)
		log.Panic("Error decoding cookie secret")
	}

	c.TokenSecret = string(tokenSecret)
	c.CookieSecret = string(cookieSecret)

	return c
}