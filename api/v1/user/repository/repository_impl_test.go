package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"bioskuy/api/v1/user/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	mockSql sqlmock.Sqlmock
	db      *sql.DB
	repo    UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	var err error
	suite.db, suite.mockSql, err = sqlmock.New()
	suite.NoError(err)
	suite.repo = NewUserRepository()
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UserRepositoryTestSuite) TestSave_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	user := entity.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Token: "token123",
		Role:  "user",
	}

	query := "INSERT INTO users \\(name, email, token, role\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id"
	suite.mockSql.ExpectQuery(query).
		WithArgs(user.Name, user.Email, user.Token, user.Role).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.Save(ctx, tx, user, c)
	suite.NoError(err)
	suite.Equal("1", result.ID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestSave_Error() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	user := entity.User{
		Name:  "John Doe",
		Email: "john@example.com",
		Token: "token123",
		Role:  "user",
	}

	query := "INSERT INTO users \\(name, email, token, role\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id"
	suite.mockSql.ExpectQuery(query).
		WithArgs(user.Name, user.Email, user.Token, user.Role).
		WillReturnError(errors.New("Insert Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Save(ctx, tx, user, c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE email = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow(1, "John Doe", "john@example.com", "token123", "user")
	suite.mockSql.ExpectQuery(query).WithArgs("john@example.com").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindByEmail(ctx, tx, "john@example.com", c)
	suite.NoError(err)
	suite.Equal("John Doe", result.Name)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail_Error() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE email = \$1`
	suite.mockSql.ExpectQuery(query).WithArgs("john@example.com").WillReturnError(errors.New("Query Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByEmail(ctx, tx, "john@example.com", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail_ErrorOnQuery() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE email = \$1`
	suite.mockSql.ExpectQuery(query).WithArgs("john@example.com").WillReturnError(errors.New("Query Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByEmail(ctx, tx, "john@example.com", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail_ErrorOnScan() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE email = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow(1, nil, "john@example.com", "token123", "user") // nil will cause scan error
	suite.mockSql.ExpectQuery(query).WithArgs("john@example.com").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByEmail(ctx, tx, "john@example.com", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail_UserNotFound() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE email = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"})
	suite.mockSql.ExpectQuery(query).WithArgs("nonexistent@example.com").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindByEmail(ctx, tx, "nonexistent@example.com", c)
	suite.Error(err)
	suite.Empty(result)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByID_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow(1, "John Doe", "john@example.com", "token123", "user")
	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindByID(ctx, tx, "1", c)
	suite.NoError(err)
	suite.Equal("John Doe", result.Name)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByID_ErrorOnQuery() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE id = \$1`
	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnError(errors.New("Query Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(ctx, tx, "1", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByID_ErrorOnScan() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow(1, nil, "john@example.com", "token123", "user") // nil will cause scan error
	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(ctx, tx, "1", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindByID_UserNotFound() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users WHERE id = \$1`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"})
	suite.mockSql.ExpectQuery(query).WithArgs("nonexistent-id").WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindByID(ctx, tx, "nonexistent-id", c)
	suite.Error(err)
	suite.Empty(result)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindAll_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow("1", "John Doe", "john@example.com", "token123", "user").
		AddRow("2", "Jane Doe", "jane@example.com", "token456", "admin")
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindAll(ctx, tx, c)
	suite.NoError(err)
	suite.Len(result, 2)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindAll_ErrorOnQuery() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users`
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("Query Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(ctx, tx, c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestFindAll_ErrorOnScan() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `SELECT id, name, email, token, role FROM users`
	rows := sqlmock.NewRows([]string{"id", "name", "email", "token", "role"}).
		AddRow(1, nil, "john@example.com", "token123", "user") // nil akan menyebabkan error pada Scan
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(ctx, tx, c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdate_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	user := entity.User{
		ID:   "1",
		Role: "admin",
	}

	query := `UPDATE users SET role = \$1 WHERE id = \$2`
	suite.mockSql.ExpectExec(query).
		WithArgs(user.Role, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.Update(ctx, tx, user, c)
	suite.NoError(err)
	suite.Equal("admin", result.Role)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdate_Error() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	user := entity.User{
		ID:   "1",
		Role: "admin",
	}

	query := `UPDATE users SET role = \$1 WHERE id = \$2`
	suite.mockSql.ExpectExec(query).
		WithArgs(user.Role, user.ID).
		WillReturnError(errors.New("Update Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.Update(ctx, tx, user, c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDelete_Success() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `DELETE FROM users WHERE id = \$1`
	suite.mockSql.ExpectExec(query).WithArgs("1").WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(ctx, tx, "1", c)
	suite.NoError(err)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDelete_Error() {
	suite.mockSql.ExpectBegin()
	tx, err := suite.db.Begin()
	suite.NoError(err)

	query := `DELETE FROM users WHERE id = \$1`
	suite.mockSql.ExpectExec(query).WithArgs("1").WillReturnError(errors.New("Delete Error"))

	ctx := context.Background()
	c, _ := gin.CreateTestContext(nil)
	err = suite.repo.Delete(ctx, tx, "1", c)
	suite.Error(err)

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	suite.NoError(suite.mockSql.ExpectationsWereMet())
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
