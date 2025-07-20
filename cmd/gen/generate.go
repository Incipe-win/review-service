package main

import (
	"errors"
	"flag"
	"fmt"
	"review-service/internal/conf"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func connectDB(cfg *conf.Data_Database) *gorm.DB {
	if cfg == nil {
		panic(errors.New("GEN: Database configuration is required"))
	}
	switch strings.ToLower(cfg.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	}
	panic(errors.New("gen: unsupported database driver"))
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/data/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(connectDB(bc.Data.Database))

	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
