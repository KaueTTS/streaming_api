package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Port                           string
	AppEnv                         string
	AppName                        string
	OTLPExporterEndpoint           string
	JWTSecret                      string
	SQLiteDatabaseURL              string
	AuthTokenExpirationTimeInHours float64
	TMDBBaseURL                    string
	TMDBAccessToken                string
)

func Init() error {
	_ = godotenv.Load()

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}

	AppEnv = os.Getenv("APP_ENV")
	if AppEnv == "" {
		AppEnv = "development"
	}

	AppName = os.Getenv("APP_NAME")
	if AppName == "" {
		AppName = "streaming_api"
	}

	OTLPExporterEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if OTLPExporterEndpoint == "" {
		OTLPExporterEndpoint = "localhost:4318"
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		return fmt.Errorf("a variável de ambiente JWT_SECRET precisa ser informada")
	}

	SQLiteDatabaseURL = os.Getenv("SQLITE_DATABASE_URL")
	if SQLiteDatabaseURL == "" {
		return fmt.Errorf("a variável de ambiente SQLITE_DATABASE_URL precisa ser informada")
	}

	AuthTokenExpirationTimeInHours = 8
	if value := os.Getenv("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS"); value != "" {
		parsedValue, err := strconv.ParseFloat(value, 64)
		if err != nil || parsedValue <= 0 {
			return fmt.Errorf("AUTH_TOKEN_EXPIRATION_TIME_IN_HOURS precisa ser um número maior que zero")
		}

		AuthTokenExpirationTimeInHours = parsedValue
	}

	TMDBBaseURL = os.Getenv("TMDB_BASE_URL")
	if TMDBBaseURL == "" {
		return fmt.Errorf("a variável de ambiente TMDB_BASE_URL precisa ser informada")
	}

	TMDBAccessToken = os.Getenv("TMDB_ACCESS_TOKEN")
	if TMDBAccessToken == "" {
		return fmt.Errorf("a variável de ambiente TMDB_ACCESS_TOKEN precisa ser informada")
	}

	return nil
}
