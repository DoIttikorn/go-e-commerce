package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func convertToInt(value, name string) int {
	result, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error convert to int name: %s err: %s", name, err)
	}
	return result
}

// func convertToBool(value string) bool {
// 	reslut, err := strconv.ParseBool(value)
// 	if err != nil {
// 		log.Fatal("Error convert to bool")
// 	}
// 	return reslut
// }

func convertToDuration(value, nameEnv string) time.Duration {
	t, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("Error convert to duration name: %s err: %s", nameEnv, err)
	}
	return time.Duration(int64(t) * int64(math.Pow10(10)))
}

func LoadConfig(path string) ConfigImpl {

	envMap, err := godotenv.Read(path)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &config{
		app: &app{
			host:         envMap["APP_HOST"],
			port:         convertToInt(envMap["APP_PORT"], "APP_PORT"),
			name:         envMap["APP_NAME"],
			version:      envMap["APP_VERSION"],
			readTimeout:  convertToDuration(envMap["APP_READ_TIMEOUT"], "APP_READ_TIMEOUT"),
			writeTimeout: convertToDuration(envMap["APP_WRITE_TIMEOUT"], "APP_WRITE_TIMEOUT"),
			bodyLimit:    convertToInt(envMap["APP_BODY_LIMIT"], "APP_BODY_LIMIT"),
			fileLimit:    convertToInt(envMap["APP_FILE_LIMIT"], "APP_FILE_LIMIT"),
			gcpBucket:    envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host:          envMap["DB_HOST"],
			port:          convertToInt(envMap["DB_PORT"], "DB_PORT"),
			protocal:      envMap["DB_PROTOCAL"],
			username:      envMap["DB_USERNAME"],
			password:      envMap["DB_PASSWORD"],
			database:      envMap["DB_DATABASE"],
			sslMode:       envMap["DB_SSL_MODE"],
			maxConnection: convertToInt(envMap["DB_MAX_CONNECTIONS"], "DB_MAX_CONNECTIONS"),
		},
		jwt: &jwt{
			adminKey:         envMap["JWT_ADMIN_KEY"],
			secertKey:        envMap["JWT_SECERT_KEY"],
			apiKey:           envMap["JWT_API_KEY"],
			accessExpiresAt:  convertToInt(envMap["JWT_ACCESS_EXPIRES"], "JWT_ACCESS_EXPIRES"),
			refreshExpiresAt: convertToInt(envMap["JWT_REFRESH_EXPIRES"], "JWT_REFRESH_EXPIRES"),
		},
	}
}

type ConfigImpl interface {
	App() AppConfigImpl
	DB() DBConfigImpl
	JWT() JWTConfigImpl
}

type config struct {
	app *app
	db  *db
	jwt *jwt
}

func (c *config) App() AppConfigImpl {
	return c.app
}
func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }
func (a *app) GCPBucket() string           { return a.gcpBucket }
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }

type AppConfigImpl interface {
	Url() string // host:port
	Name() string
	Version() string
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
	GCPBucket() string
	Host() string
	Port() int
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int
	fileLimit    int
	gcpBucket    string
}
type DBConfigImpl interface {
	Url() string
	MaxOpenConnection() int
}

type db struct {
	host          string
	port          int
	protocal      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

func (c *config) DB() DBConfigImpl {
	return c.db
}
func (db *db) Url() string {
	// return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", db.protocal, db.username, db.password, db.host, db.port, db.database, db.sslMode)
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.host, db.port, db.username, db.password, db.database, db.sslMode,
	)
}
func (db *db) MaxOpenConnection() int { return db.maxConnection }

type JWTConfigImpl interface {
	SecretKey() []byte
	AdminKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
	SetJwtAccessExpires(int)
	SetJwtRefreshExpires(int)
}

type jwt struct {
	adminKey         string
	secertKey        string
	apiKey           string
	accessExpiresAt  int //seconds
	refreshExpiresAt int //seconds
}

func (c *config) JWT() JWTConfigImpl {
	return c.jwt
}

func (j *jwt) SecretKey() []byte          { return []byte(j.secertKey) }
func (j *jwt) AdminKey() []byte           { return []byte(j.adminKey) }
func (j *jwt) ApiKey() []byte             { return []byte(j.apiKey) }
func (j *jwt) AccessExpiresAt() int       { return j.accessExpiresAt }
func (j *jwt) RefreshExpiresAt() int      { return j.refreshExpiresAt }
func (j *jwt) SetJwtAccessExpires(t int)  { j.accessExpiresAt = t }
func (j *jwt) SetJwtRefreshExpires(t int) { j.refreshExpiresAt = t }
