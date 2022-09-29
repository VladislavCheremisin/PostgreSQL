package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"sync"
	"time"
)

//type ContextKey int

//const ContextKeyDB ContextKey = iota + 1

var (
	db    *pgxpool.Pool
	dbMux = &sync.Mutex{}
)

type FoundName struct {
	FirstName string `json:"first_name"`
	Email     string `email:"email"`
}

type DB interface {
	GetNamesByEmailPrefix(ctx context.Context, prefix string) ([]*FoundName, error)
	Close()
}

type ConnString struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDB(connStr *ConnString) (DB, error) {
	pool, err := getConn(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get a connection pool: %w", err)
	}
	return &conn{
		db: pool,
	}, nil
	//gormDB, err := newGormDB(connStr)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to open a gorm connection: %w", err)
	//}
	//return gormDB, nil
}

type conn struct {
	db *pgxpool.Pool
}

func (c *conn) GetNamesByEmailPrefix(ctx context.Context, prefix string) ([]*FoundName, error) {
	rows, err := c.db.Query(
		context.Background(),
		`SELECT firstname, email
		FROM users
		WHERE email LIKE $1 || '%'`,
		prefix,
	)

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	names := make([]*FoundName, 0)
	for rows.Next() {
		p := &FoundName{}
		if err := rows.Scan(&p.FirstName, &p.Email); err != nil {
			return nil, fmt.Errorf("failed to scan a received phone: %w", err)
		}
		names = append(names, p)
	}
	return names, nil
}

func (c *conn) Close() {
	c.db.Close()
}

func getConn(connStr *ConnString) (*pgxpool.Pool, error) {
	dbMux.Lock()
	defer dbMux.Unlock()
	if db != nil {
		return db, nil
	}

	var err error
	db, err = initPGXPool(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a PGX pool: %w", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}
	return db, nil
}

func initPGXPool(c *ConnString) (*pgxpool.Pool, error) {
	connStr, err := composeConnectionString(c)
	if err != nil {
		return nil, fmt.Errorf("failed to compose the connection string: %w", err)
	}
	cfg, err := getPGXPoolConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the PGX pool config: %w", err)
	}
	db, err = pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the postgres DB using a PGX connection pool: %w", err)
	}
	return db, nil
}

func getPGXPoolConfig(connStr string) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create the PGX pool config from connection string: %w", err)
	}
	cfg.ConnConfig.ConnectTimeout = time.Second * 1
	cfg.ConnConfig.Logger = logrusadapter.NewLogger(
		&logrus.Logger{
			Out:          os.Stdout,
			Formatter:    new(logrus.JSONFormatter),
			Hooks:        make(logrus.LevelHooks),
			Level:        logrus.InfoLevel,
			ExitFunc:     os.Exit,
			ReportCaller: false,
		})
	return cfg, nil
}

func composeConnectionString(c *ConnString) (string, error) {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		url.QueryEscape(c.Host),
		url.QueryEscape(c.Port),
		url.QueryEscape(c.DBName),
	), nil
}
