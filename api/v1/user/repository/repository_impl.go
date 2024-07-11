package repository

import (
	"bioskuy/api/v1/user/entity"
	"bioskuy/exception"
	"context"
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
)

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Save(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error) {

	query := "INSERT INTO users (name, email, token) VALUES ($1, $2, $3) RETURNING id"

	err := tx.QueryRowContext(ctx, query, user.Name, user.Email, user.Token).Scan(&user.ID)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, tx *sql.Tx, email string, c *gin.Context) (entity.User, error){

	query := `SELECT id, name, email, token, role FROM users WHERE email = $1`
	
	user := entity.User{}
	rows, err := tx.QueryContext(ctx, query, email)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  user, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Role,)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return  user, err
		}

		return user, nil
	}else{
		return user, errors.New("user not found")
	}
}

func (r *userRepository) FindByID(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) (entity.User, error){

	query := `SELECT id, name, email, token, role FROM users WHERE id = $1`
	
	user := entity.User{}
	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  user, err
	}
	defer rows.Close()

	if rows.Next(){
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Token, &user.Role,)
		if err != nil {
			c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
			return  user, err
		}

		return user, nil
	}else{
		return user, errors.New("user not found")
	}
}

func (r *userRepository) FindAll(ctx context.Context, tx *sql.Tx, c *gin.Context) ([]entity.User, error){

	query := `SELECT id, name, email, token, role FROM users`

	users := []entity.User{}
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := entity.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Token,&user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) UpdateToken(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error){

	query := `UPDATE users SET token = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, user.Token, user.ID)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  user, err
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, tx *sql.Tx, user entity.User, c *gin.Context) (entity.User, error){

	query := `UPDATE users SET role = $1 WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, user.Role, user.ID)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  user, err
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, tx *sql.Tx, id string, c *gin.Context) error{
	query := `DELETE FROM users WHERE id = $1`

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		c.Error(exception.InternalServerError{Message: err.Error()}).SetType(gin.ErrorTypePublic)
		return  err
	}

	return nil
}
