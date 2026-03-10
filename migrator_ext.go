package postgres

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (m Migrator) DropConstraint(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		constraint, _ := m.GuessConstraintInterfaceAndTable(stmt, name)
		if constraint != nil {
			name = constraint.GetName()
		}
		// 42704
		err := m.DB.Exec("ALTER TABLE ? DROP CONSTRAINT ?", m.CurrentTable(stmt), clause.Column{Name: name}).Error
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Code == "42704" {
					return nil
				}
			}
			return err
		}
		return nil
	})
}

func GetTableName(schema, tableName string) string {
	return fmt.Sprintf("%s.%s", schema, tableName)
}
