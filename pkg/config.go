package pkg

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HydrangeaGardenURL    string
	HydrangeaRALogURL     string
	HydrangeaSensorLogURL string
	HydrangeaRAURL        string
	GardenId              string
}

func LoadConfigFromEnv() Config {
	return Config{
		HydrangeaGardenURL:    os.Getenv("HYDRANGEA_GARDEN_URL"),
		HydrangeaRALogURL:     os.Getenv("HYDRANGEA_RALOG_URL"),
		HydrangeaSensorLogURL: os.Getenv("HYDRANGEA_SENSORLOG_URL"),
		HydrangeaRAURL:        os.Getenv("HYDRANGEA_RA_URL"),
		GardenId:              os.Getenv("GARDEN_ID"),
	}
}
