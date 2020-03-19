package cmd

// ServerCmd runs the web-server
//
// Connection string should be provided in format: dbdriver://uname:password@address:port/dbname?[param1=][&param2=][...etc]
// Examples: postgres://dbuser:dbpass@localhost:8080/dbname?sslmode=verify-ca
type ServerCmd struct {
	ImgStoragePath      string `long:"imgstoragepath" env:"IMGSTORAGEPATH" required:"true" description:"local storage path for images"`
	JWTSecret           string `long:"jwtsecret" env:"JWTSECRET" required:"true" description:"jwt secret for hashing"`
	ServiceURL          string `long:"serviceurl" env:"SERVICEURL" required:"true" description:"url to this web-server"`
	ImageProcServiceURL string `long:"imgprocurl" env:"IMGPROCURL" required:"true" description:"image proc service url"`

	DefaultAdmin struct {
		Email    string `long:"email" env:"EMAIL" required:"true" description:"email of default admin"`
		Password string `long:"password" env:"PWD" required:"true" description:"password of default admin"`
	} `group:"admin" namespace:"admin" env-namespace:"ADMIN"`

	Database struct {
		ConnStr    string `long:"connstr" env:"CONN_STR" required:"true" description:"connection string to database"`
		BCryptCost string `long:"bcryptcost" env:"BCRYPT_COST" required:"true" description:"cost of bcrypt hashing algo"`
	} `group:"db" namespace:"db" env-namespace:"DB"`
}

// Execute starts the listening and serves http
func (s *ServerCmd) Execute(_ []string) error {
	panic("implement me")
}
