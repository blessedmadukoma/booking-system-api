package util

import (
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application
type Config struct {
	GinMode             string        `mapstructure:"GIN_MODE"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	PORT                string        `mapstructure:"PORT"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	SendGridAPIKey      string        `mapstructure:"SENDGRID_API_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	Cloudinary          struct {
		CloudiName      string `mapstructure:"CLOUDI_NAME"`
		CloudiAPIKey    string `mapstructure:"CLOUDI_API_KEY"`
		CloudiAPISecret string `mapstructure:"CLOUDI_API_SECRET"`
		CloudiURL       string `mapstructure:"CLOUDINARY_URL"`
	}
	Limiter struct {
		RPS     float64
		BURST   int
		ENABLED bool
	}
}

// rateLimitValues retreives the values for the rate limiter from the env
func rateLimitValues() (int, int, bool) {

	rps, err := strconv.Atoi(os.Getenv("LIMITER_RPS"))
	if err != nil {
		log.Fatal("Error retrieving rps value:", err)
	}
	burst, err := strconv.Atoi(os.Getenv("LIMITER_BURST"))
	if err != nil {
		log.Fatal("Error retrieving burst value:", err)
	}
	enabled, err := strconv.ParseBool(os.Getenv("LIMITER_ENABLED"))
	if err != nil {
		log.Fatal("Error retrieving enabled value:", err)
	}

	return rps, burst, enabled
}

// LoadConfig reads configuration from file or env variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// viper.SetConfigName("app")
	// viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// LoadEnvConfig reads configuration from file or env variables
func LoadEnvConfig(path string) (config Config, err error) {
	config.GinMode = os.Getenv("GIN_MODE")
	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DBSource = os.Getenv("DB_SOURCE")
	config.PORT = os.Getenv("PORT")
	config.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	config.SendGridAPIKey = os.Getenv("SENDGRID_API_KEY")
	config.AccessTokenDuration, _ = time.ParseDuration(os.Getenv("ACCESS_TOKEN_DURATION"))
	config.Cloudinary.CloudiURL = os.Getenv("CLOUDINARY_URL")

	config.Cloudinary.CloudiName = os.Getenv("CLOUDI_NAME")
	config.Cloudinary.CloudiAPIKey = os.Getenv("CLOUDI_API_KEY")
	config.Cloudinary.CloudiAPISecret = os.Getenv("CLOUDI_API_SECRET")
	// config.RefreshTokenDuration, _ = time.ParseDuration(os.Getenv("REFRESH_TOKEN_DURATION"))

	// retrieve rate limit values
	rateRPS, rateBurst, rateEnabled := rateLimitValues()
	config.Limiter.RPS = float64(rateRPS)
	config.Limiter.BURST = rateBurst
	config.Limiter.ENABLED = rateEnabled

	return
}
