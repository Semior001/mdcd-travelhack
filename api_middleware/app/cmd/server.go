package cmd

import (
	"github.com/Semior001/mdcd-travelhack/app/rest"
	"github.com/Semior001/mdcd-travelhack/app/rest/private"
	"github.com/Semior001/mdcd-travelhack/app/store/image"
	"github.com/Semior001/mdcd-travelhack/app/store/user"
	"github.com/pkg/errors"
	"log"
	"time"
)

// ServeCommand to run the server
type ServeCommand struct {
	Database
	JWTSecret  string `long:"jwtsecret" env:"JWTSECRET" required:"true" description:"jwt secret for hashing"`
	ServiceURL string `long:"serviceurl" env:"SERVICEURL" required:"true" description:"url to this web-server"`
	Force      bool   `long:"force" env:"DBMIGRATEFORCE" required:"false" description:"force to migrate db"`

	Email    string `long:"email" env:"RU_EMAIL" required:"true" description:"email of registering user"`
	Password string `long:"password" env:"RU_PASSWORD" required:"true" description:"password of registering user"`

	ImageProcServiceURL string `long:"imgprocurl" env:"IMGPROCURL" required:"true" description:"image proc service url"`

	CommonOptions
}

// Execute runs web server
func (s *ServeCommand) Execute(_ []string) error {
	// initializing user service
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

	err = us.Migrate(s.Force)

	if err != nil {
		return errors.Wrapf(err, "failed to migrate user service")
	}

	// initializing images service
	im, err := image.NewService(image.ServiceOpts{
		Driver:           s.Database.Driver,
		User:             s.Database.User,
		Password:         s.Database.Password,
		Source:           s.Database.Source,
		LoggerFlags:      s.LoggerFlags,
		LocalStoragePath: s.LocalStoragePath,
	})

	if err != nil {
		return errors.Wrapf(err, "failed to create image service")
	}

	err = im.Migrate(s.Force)

	if err != nil {
		return errors.Wrapf(err, "failed to migrate image service")
	}

	// registering admin, if not exists
	log.Printf("[DEBUG] creating admin user %s", s.Email)
	id, err := us.PutUser(user.User{
		Email:    s.Email,
		Password: s.Password,
		Privileges: map[string]bool{
			user.PrivilegeAdmin: true,
		},
	})
	if err != nil {
		log.Printf("failed to create admin user %s: %+v", s.Email, err)
	} else {
		log.Printf("[INFO] user %s has been created successfully with id: %d", s.Email, id)
	}

	// initializing rest
	r := rest.Rest{
		Version:        s.Version,
		AppName:        s.AppName,
		AppAuthor:      s.AppAuthor,
		JWTSecret:      s.JWTSecret,
		ServiceURL:     s.ServiceURL,
		UserService:    *us,
		ImageService:   *im,
		UserController: private.UserController{ServiceUsr: *us},
		ImageController: private.ImageController{
			ServiceImg:          *im,
			ServiceUsr:          *us,
			ImageProcServiceURL: s.ImageProcServiceURL,
		},
		Auth: struct {
			TTL struct {
				JWT    time.Duration
				Cookie time.Duration
			}
		}{
			TTL: struct {
				JWT    time.Duration
				Cookie time.Duration
			}{
				JWT:    time.Hour * 5000,
				Cookie: time.Hour * 5000,
			},
		},
	}
	r.Run(8080)
	return nil
}
