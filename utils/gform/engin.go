package gform

import (
	"database/sql"
	"fmt"
)

// TAGNAME ...
var TAGNAME = "gform"

// IGNORE ...
var IGNORE = "-"

type cluster struct {
	master     []*sql.DB
	masterSize int
	slave      []*sql.DB
	slaveSize  int
}

// Engin ...
type Engin struct {
	config *ConfigCluster
	driver string
	prefix string
	dbs    *cluster
	logger ILogger
}

var _ IEngin = (*Engin)(nil)

// NewEngin : init Engin struct pointer
func NewEngin(conf ...interface{}) (e *Engin, err error) {
	engin := new(Engin)
	if len(conf) == 0 {
		return
	}

	engin.Use(DefaultLogger())

	switch conf[0].(type) {
	case *Config:
		err = engin.bootSingle(conf[0].(*Config))
	case *ConfigCluster:
		engin.config = conf[0].(*ConfigCluster)
		err = engin.bootCluster()
	default:
		panic(fmt.Sprint("Open() need *gform.Config or *gform.ConfigCluster param, also can empty for build sql string only, but ",
			conf, " given"))
	}

	return engin, err
}

// Use ...
func (c *Engin) Use(closers ...func(e *Engin)) {
	for _, closer := range closers {
		closer(c)
	}
}

// Ping ...
func (c *Engin) Ping() error {
	//for _,item := range c.dbs.master {
	return c.GetQueryDB().Ping()
}

func (c *Engin) TagName(arg string) {
	//c.tagName = arg
	TAGNAME = arg
}

func (c *Engin) IgnoreName(arg string) {
	//c.ignoreName = arg
	IGNORE = arg
}

func (c *Engin) SetPrefix(pre string) {
	c.prefix = pre
}

func (c *Engin) GetPrefix() string {
	return c.prefix
}

// GetDriver ...
func (c *Engin) GetDriver() string {
	return c.driver
}

func (c *Engin) GetQueryDB() *sql.DB {
	if c.dbs.slaveSize == 0 {
		return c.GetExecuteDB()
	}
	var randint = getRandomInt(c.dbs.slaveSize)
	return c.dbs.slave[randint]
}

func (c *Engin) GetExecuteDB() *sql.DB {
	if c.dbs.masterSize == 0 {
		return nil
	}
	var randint = getRandomInt(c.dbs.masterSize)
	return c.dbs.master[randint]
}

// GetLogger ...
func (c *Engin) GetLogger() ILogger {
	return c.logger
}

// SetLogger ...
func (c *Engin) SetLogger(lg ILogger) {
	c.logger = lg
}

func (c *Engin) bootSingle(conf *Config) error {
	var cc = new(ConfigCluster)
	cc.Master = append(cc.Master, *conf)
	c.config = cc
	return c.bootCluster()
}

func (c *Engin) bootCluster() error {
	//fmt.Println(len(c.config.Slave))
	if len(c.config.Slave) > 0 {
		for _, item := range c.config.Slave {
			if c.config.Driver != "" {
				item.Driver = c.config.Driver
			}
			if c.config.Prefix != "" {
				item.Prefix = c.config.Prefix
			}
			db, err := c.bootReal(item)
			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.slave = append(c.dbs.slave, db)
			c.dbs.slaveSize++
			c.driver = item.Driver
		}
	}
	var pre, dr string
	if len(c.config.Master) > 0 {
		for _, item := range c.config.Master {
			if c.config.Driver != "" {
				item.Driver = c.config.Driver
			}
			if c.config.Prefix != "" {
				item.Prefix = c.config.Prefix
			}
			db, err := c.bootReal(item)

			if err != nil {
				return err
			}
			if c.dbs == nil {
				c.dbs = new(cluster)
			}
			c.dbs.master = append(c.dbs.master, db)
			c.dbs.masterSize = c.dbs.masterSize + 1
			c.driver = item.Driver
			//fmt.Println(c.dbs.masterSize)
			if item.Prefix != "" {
				pre = item.Prefix
			}
			if item.Driver != "" {
				dr = item.Driver
			}
		}
	}
	if pre != "" && c.prefix == "" {
		c.prefix = pre
	}
	if dr != "" && c.driver == "" {
		c.driver = dr
	}

	return nil
}

// boot sql driver
func (c *Engin) bootReal(dbConf Config) (db *sql.DB, err error) {
	db, err = sql.Open(dbConf.Driver, dbConf.Dsn)
	if err != nil {
		return
	}

	err = db.Ping()
	if err != nil {
		return
	}

	if dbConf.SetMaxOpenConns > 0 {
		db.SetMaxOpenConns(dbConf.SetMaxOpenConns)
	}
	if dbConf.SetMaxIdleConns > 0 {
		db.SetMaxIdleConns(dbConf.SetMaxIdleConns)
	}

	return
}

func (c *Engin) NewSession() ISession {
	return NewSession(c)
}

func (c *Engin) NewOrm() IOrm {
	return NewOrm(c)
}
