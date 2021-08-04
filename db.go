package gomigrate

import (
	"fmt"
	"strings"
)

type Migratable interface {
	SelectMigrationTableSql() string
	CreateMigrationTableSql() string
	GetMigrationSql() string
	MigrationLogInsertSql() string
	MigrationLogDeleteSql() string
	GetMigrationCommands(string) []string
	SetMigrationTableName(string) error
}

// POSTGRES

type Postgres struct {
	TableName string
}

func (p *Postgres) SelectMigrationTableSql() string {
	fmt.Println(p.TableName)
	return "SELECT tablename FROM pg_catalog.pg_tables WHERE tablename = '" + p.TableName + "'"
}

func (p *Postgres) CreateMigrationTableSql() string {
	return `CREATE TABLE ` + p.TableName + ` (
                  id           SERIAL       PRIMARY KEY,
                  migration_id BIGINT       UNIQUE NOT NULL,
				  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
				  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
                )`
}

func (p *Postgres) GetMigrationSql() string {
	return `SELECT migration_id FROM ` + p.TableName + ` WHERE migration_id = $1`
}

func (p *Postgres) MigrationLogInsertSql() string {
	return "INSERT INTO " + p.TableName + " (migration_id) values ($1)"
}

func (p *Postgres) MigrationLogDeleteSql() string {
	return "DELETE FROM " + p.TableName + " WHERE migration_id = $1"
}

func (p *Postgres) GetMigrationCommands(sql string) []string {
	return []string{sql}
}

func (p *Postgres) SetMigrationTableName(migrateTableName string) error {
	p.TableName = migrateTableName
	return nil
}

// MYSQL

type Mysql struct{ TableName string }

func (m *Mysql) SelectMigrationTableSql() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_name = " + m.TableName + " AND table_schema = (SELECT DATABASE())"
}

func (m *Mysql) CreateMigrationTableSql() string {
	return `CREATE TABLE ` + m.TableName + ` (
                  id           INT          NOT NULL AUTO_INCREMENT,
                  migration_id BIGINT       NOT NULL UNIQUE,
				  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
				  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
                  PRIMARY KEY (id)
                )`
}

func (m *Mysql) GetMigrationSql() string {
	return `SELECT migration_id FROM ` + m.TableName + ` WHERE migration_id = ?`
}

func (m *Mysql) MigrationLogInsertSql() string {
	return "INSERT INTO " + m.TableName + " (migration_id) values (?)"
}

func (m *Mysql) MigrationLogDeleteSql() string {
	return "DELETE FROM " + m.TableName + " WHERE migration_id = ?"
}

func (m *Mysql) GetMigrationCommands(sql string) []string {
	count := strings.Count(sql, ";")
	commands := strings.SplitN(string(sql), ";", count)
	return commands
}

func (m *Mysql) SetMigrationTableName(migrateTableName string) error {
	m.TableName = migrateTableName
	return nil
}

// MARIADB

type Mariadb struct {
	TableName string
	Mysql
}

// SQLITE3

type Sqlite3 struct{ TableName string }

func (s *Sqlite3) SelectMigrationTableSql() string {
	return "SELECT name FROM sqlite_master WHERE type = 'table' AND name = " + s.TableName
}

func (s *Sqlite3) CreateMigrationTableSql() string {
	return `CREATE TABLE ` + s.TableName + ` (
  id INTEGER PRIMARY KEY,
  migration_id INTEGER NOT NULL UNIQUE
)`
}

func (s *Sqlite3) GetMigrationSql() string {
	return "SELECT migration_id FROM " + s.TableName + " WHERE migration_id = ?"
}

func (s *Sqlite3) MigrationLogInsertSql() string {
	return "INSERT INTO " + s.TableName + " (migration_id) values (?)"
}

func (s *Sqlite3) MigrationLogDeleteSql() string {
	return "DELETE FROM " + s.TableName + " WHERE migration_id = ?"
}

func (s *Sqlite3) GetMigrationCommands(sql string) []string {
	return []string{sql}
}

func (s *Sqlite3) SetMigrationTableName(migrateTableName string) error {
	s.TableName = migrateTableName
	return nil
}

// SqlServer

type SqlServer struct{ TableName string }

func (s *SqlServer) SelectMigrationTableSql() string {
	return "SELECT name FROM  sys.objects  WHERE object_id = object_id(" + s.TableName + ")"
}

func (s *SqlServer) CreateMigrationTableSql() string {
	return `CREATE TABLE ` + s.TableName + ` (
                  id           INT     IDENTITY(1,1)     PRIMARY KEY,
                  migration_id BIGINT  NOT NULL
                )`
}

func (s *SqlServer) GetMigrationSql() string {
	return `SELECT migration_id FROM ` + s.TableName + ` WHERE migration_id = ?`
}

func (s *SqlServer) MigrationLogInsertSql() string {
	return "INSERT INTO " + s.TableName + " (migration_id) values (?)"
}

func (s *SqlServer) MigrationLogDeleteSql() string {
	return "DELETE FROM " + s.TableName + " WHERE migration_id = ?"
}

func (s *SqlServer) GetMigrationCommands(sql string) []string {
	return []string{sql}
}
func (s *SqlServer) SetMigrationTableName(migrateTableName string) error {
	s.TableName = migrateTableName
	return nil
}
