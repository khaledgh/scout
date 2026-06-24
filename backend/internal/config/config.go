package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       DBConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Upload   UploadConfig
	Geo      GeoConfig
	Gamify   GamifyConfig
	VAPID    VAPIDConfig
	SMTP     SMTPConfig
	Jobs     JobsConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     string
	URL      string
	Timezone string
}

type DBConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	Params          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	BcryptCost int
}

type CORSConfig struct {
	AllowedOrigins string
}

type UploadConfig struct {
	Dir          string
	MaxSizeMB    int64
	AllowedTypes string
	PublicPath   string
}

type GeoConfig struct {
	GeofenceRadiusMeters float64
	QRSigningSecret      string
}

type GamifyConfig struct {
	XPPerAttendance int
	XPPerQuizPass   int
	XPPerLeadership int
	LevelBaseXP     int
}

type VAPIDConfig struct {
	PublicKey  string
	PrivateKey string
	Subject    string
}

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

type JobsConfig struct {
	ReminderLookaheadHours int
	MonthlyReportDay       int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.App = AppConfig{
		Name:     requireEnv("APP_NAME"),
		Env:      getEnv("APP_ENV", "development"),
		Port:     getEnv("APP_PORT", "8080"),
		URL:      getEnv("APP_URL", "http://localhost:8080"),
		Timezone: getEnv("APP_TIMEZONE", "UTC"),
	}

	cfg.DB = DBConfig{
		Host:            requireEnv("DB_HOST"),
		Port:            getEnv("DB_PORT", "3306"),
		Name:            requireEnv("DB_NAME"),
		User:            requireEnv("DB_USER"),
		Password:        getEnv("DB_PASSWORD", ""),
		Params:          getEnv("DB_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
		MaxOpenConns:    parseInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    parseInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: time.Duration(parseInt("DB_CONN_MAX_LIFETIME", 300)) * time.Second,
	}

	cfg.JWT = JWTConfig{
		Secret:     requireEnv("JWT_SECRET"),
		AccessTTL:  time.Duration(parseInt("JWT_ACCESS_TTL", 900)) * time.Second,
		RefreshTTL: time.Duration(parseInt("JWT_REFRESH_TTL", 604800)) * time.Second,
		BcryptCost: parseInt("BCRYPT_COST", 12),
	}

	cfg.CORS = CORSConfig{
		AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
	}

	cfg.Upload = UploadConfig{
		Dir:          getEnv("UPLOAD_DIR", "./uploads"),
		MaxSizeMB:    int64(parseInt("UPLOAD_MAX_SIZE_MB", 15)),
		AllowedTypes: getEnv("UPLOAD_ALLOWED_TYPES", "image/jpeg,image/png,image/webp,video/mp4"),
		PublicPath:   getEnv("PUBLIC_UPLOAD_PATH", "/uploads"),
	}

	cfg.Geo = GeoConfig{
		GeofenceRadiusMeters: parseFloat("GEOFENCE_RADIUS_METERS", 150),
		QRSigningSecret:      requireEnv("QR_SIGNING_SECRET"),
	}

	cfg.Gamify = GamifyConfig{
		XPPerAttendance: parseInt("XP_PER_ATTENDANCE", 10),
		XPPerQuizPass:   parseInt("XP_PER_QUIZ_PASS", 20),
		XPPerLeadership: parseInt("XP_PER_LEADERSHIP", 50),
		LevelBaseXP:     parseInt("LEVEL_BASE_XP", 100),
	}

	cfg.VAPID = VAPIDConfig{
		PublicKey:  getEnv("VAPID_PUBLIC_KEY", ""),
		PrivateKey: getEnv("VAPID_PRIVATE_KEY", ""),
		Subject:    getEnv("VAPID_SUBJECT", ""),
	}

	cfg.SMTP = SMTPConfig{
		Host:     getEnv("SMTP_HOST", ""),
		Port:     getEnv("SMTP_PORT", "587"),
		User:     getEnv("SMTP_USER", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "no-reply@example.com"),
	}

	cfg.Jobs = JobsConfig{
		ReminderLookaheadHours: parseInt("REMINDER_LOOKAHEAD_HOURS", 24),
		MonthlyReportDay:       parseInt("MONTHLY_REPORT_DAY", 1),
	}

	return cfg, nil
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Params)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %s is not set", key))
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func parseFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return fallback
}
