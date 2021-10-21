package repository

import (
	"context"
	"fmt"
	"github.com/portnyagin/practicum_project/internal/app/database/ddl"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
)

func InitDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, ddl.CreateDatabaseStructure)
	if err != nil {
		return err
	}
	fmt.Println("Database structure created successfully")
	return nil
}

func ClearDatabase(ctx context.Context, h basedbhandler.DBHandler) error {
	err := h.Execute(ctx, ddl.ClearDatabaseStructure)
	if err != nil {
		return err
	}
	return nil
}
