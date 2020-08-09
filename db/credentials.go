package db

import "strconv"

type Credentials struct {
	username   string
	password   string
	port       int
	opts       string
	database   string
	collection string
}

func NewCredentials(user string, pass string, port int, opts string, database string, collection string) Credentials {
	return Credentials{
		user,
		pass,
		port,
		opts,
		database,
		collection,
	}
}

func (c *Credentials) GetURI() string {
	uri := "mongodb://" + c.username + ":" + c.password + "@localhost:" + strconv.Itoa(c.port) + "/" + c.opts
	return uri
}

var DefaultCredentials Credentials = NewCredentials(
	"productListUser",
	"productListPassword",
	27017,
	"?authSource=admin&readPreference=primary&ssl=false",
	"promotions",
	"products",
)
