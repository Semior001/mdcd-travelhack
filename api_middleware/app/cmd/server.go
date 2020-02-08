package cmd

import (
	"github.com/Semior001/mdcd-travelhack/app/rest"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/pkg/errors"
)

// ServeCommand to run the server
type ServeCommand struct {
	Database
	JWTSecret  string `long:"jwtsecret" env:"JWTSECRET" required:"true" description:"jwt secret for hashing"`
	MediaPath  string `long:"mediapath" env:"MEDIAPATH" required:"true" description:"path to local media"`
	ServiceURL string `long:"serviceurl" env:"SERVICEURL" required:"true" description:"url to this web-server"`
	CommonOptions
}

// Execute runs web server
func (s *ServeCommand) Execute(args []string) error {
	us, err := user.NewService(user.ServiceOpts{
		Driver:      s.Database.Driver,
		User:        s.Database.User,
		Password:    s.Database.Password,
		Source:      s.Database.Source,
		LoggerFlags: s.LoggerFlags,
		BcryptCost:  s.Hashing.BcryptCost,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create user service")
	}
	r := rest.Rest{
		Version:     s.Version,
		AppName:     s.AppName,
		AppAuthor:   s.AppAuthor,
		JWTSecret:   s.JWTSecret,
		ServiceURL:  s.ServiceURL,
		UserService: *us,
	}
	r.Run(8080)
	return nil
}
