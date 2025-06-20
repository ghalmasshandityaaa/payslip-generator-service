package config

type Config struct {
	App      appConfig
	Security securityConfig
	Logger   loggerConfig
	Postgres postgresConfig
}

type appConfig struct {
	Name         string
	Version      string
	Domain       string
	Host         string
	Port         int
	Env          string
	Debug        bool
	ReadTimeout  int
	WriteTimeout int
	SSL          bool
	Prefork      bool
}

type securityConfig struct {
	Cors      corsConfig
	Csrf      csrfConfig
	Jwt       jwtConfig
	RateLimit rateLimitConfig
	Crypto    cryptoConfig
	Cookie    cookieConfig
}

type corsConfig struct {
	AllowedOrigins   string
	AllowedMethods   string
	AllowCredentials bool
}

type csrfConfig struct {
	Enabled    bool
	CookieName string
	HeaderName string
}

type jwtConfig struct {
	Issuer               string
	Audience             string
	Subject              string
	SigningMethod        string
	AccessTokenLifetime  int
	AccessTokenSecret    string
	RefreshTokenLifetime int
	RefreshTokenSecret   string
}

type rateLimitConfig struct {
	Duration    int
	MaxRequests int
}

type cryptoConfig struct {
	Key string
}

type loggerConfig struct {
	Level  int
	Pretty bool
}

type postgresConfig struct {
	ConnMaxIdleTime int
	ConnMaxLifetime int
	MaxIdleCons     int
	MaxOpenCons     int
	User            string
	Password        string
	Host            string
	Port            int
	Dbname          string
	Driver          string
	SSLMode         string
	DryRun          bool
}

type cookieConfig struct {
	Name     string
	Secure   bool
	HttpOnly bool
	SameSite string
	Domain   string
	MaxAge   int
	Key      string
}
