package database

import (
    log "github.com/cihub/seelog"
    "github.com/dpx-infinity/imaged/common"
    "github.com/dpx-infinity/imaged/config"
    "github.com/jmoiron/sqlx"
    _ "github.com/mattn/go-sqlite3"
    "os"
)

type Database struct {
    *sqlx.DB
}

func Initialize(conf *config.Config) (*Database, error) {
    return InitializeWithPath(conf.Paths.DatabaseFile())
}

func InitializeWithPath(databaseFile string) (*Database, error) {
    dbExists, err := checkDatabaseExists(databaseFile)
    if err != nil {
        return nil, common.NewError("Cannot access database file", err)
    }

    db, err := sqlx.Open("sqlite3", databaseFile)
    if err != nil {
        return nil, common.NewError("Cannot open database", err)
    }

    database := &Database{db}

    if !dbExists {
        log.Info("Creating tables")
        if err = database.createTables(); err != nil {
            return nil, common.NewError("Cannot create tables", err)
        }
    }

    return database, nil
}

func checkDatabaseExists(databaseFile string) (bool, error) {
    if _, err := os.Stat(databaseFile); err != nil {
        if os.IsNotExist(err) {
            log.Infof("Database file does not exist: %s", databaseFile)
            return false, nil
        } else {
            return false, err
        }
    }
    log.Infof("Found database file: %s", databaseFile)
    return true, nil
}

func (db *Database) createTables() error {
    if err := db.Tags().createTable(); err != nil {
        return err
    }

    if err := db.Files().createTable(); err != nil {
        return err
    }

    if err := db.Groups().createTable(); err != nil {
        return err
    }

    return db.establishKeys()
}
