package repository

import (
	"bioskuy/api/v1/studio/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type studioRepository struct {
}

func NewStudioRepository() StudioRepository {
	return &studioRepository{}
}

func (r *studioRepository) Save(ctx context.Context, tx *sql.Tx, studio entity.Studio, c *gin.Context) (entity.Studio, error){
	query := "INSERT INTO studios (name, capacity) VALUES ($1, $2) RETURNING id"

	err := tx.QueryRowContext(ctx, query, studio.Name, studio.Capacity).Scan(&studio.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return studio, err
	}

	return studio, nil
}

func (r *studioRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.Studio, error){

	fmt.Println(id)

	query := `SELECT id, name, capacity FROM studios WHERE id = $1`
	
	studio := entity.Studio{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  studio, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(&studio.ID, &studio.Name, &studio.Capacity)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return  studio, err
		}

		return studio, nil
	}else{
		return studio, errors.New("studio not found")
	}
}

func (r *studioRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.Studio, error){

	query := `SELECT id, name, capacity FROM studios`

	studios := []entity.Studio{}
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  studios, err
	}
	defer rows.Close()

	for rows.Next() {
		studio := entity.Studio{}
		if err := rows.Scan(&studio.ID, &studio.Name, &studio.Capacity); err != nil {
			return nil, err
		}
		studios = append(studios, studio)
	}
	return studios, nil
}

func (r *studioRepository) Update(ctx context.Context, tx *sql.Tx, studio entity.Studio, c *gin.Context) (entity.Studio, error){

	query := `UPDATE studios SET name = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, studio.Name, studio.ID)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  studio, err
	}

	return studio, nil
}

func (r *studioRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error{
	query := `DELETE FROM studios WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  err
	}

	return nil
}
