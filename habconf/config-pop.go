package habconf

import (
	"github.com/gobuffalo/pop/v5"
)

// [ISSUE #7] Here obviously - we'll want a name or whatever
// [ISSUE #7] This willneed to be made into an array
type tomlDB pop.ConnectionDetails


const defaultSqliteURL = "sqlite3://.habitat/db.sqlite3"


func (c * Config) processDB() error  { 


	if c.toml.Database.Dialect == "" { 
		c.toml.Database.Dialect = "sqlite"
		c.toml.Database.URL = defaultSqliteURL
	}

	return nil 
}



// [ISSUE #7] Parameter here
func (c Config) PopConnectionDetails() * pop.ConnectionDetails { 
	return (*pop.ConnectionDetails)(&c.toml.Database)
}
func (c Config) NewConnection() (*pop.Connection, error) { 
	return pop.NewConnection( c.PopConnectionDetails() )
}