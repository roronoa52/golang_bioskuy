package repository

import (
	"bioskuy/api/v1/genretomovie/entity"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenreToMovieRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
	repo    GenreToMovieRepository
}

var mockingGenreToMovie = entity.GenreToMovie{
	ID:      uuid.New().String(),
	GenreID: uuid.New().String(),
	MovieID: uuid.New().String(),
}

func (suite *GenreToMovieRepositoryTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(suite.T(), err)
	suite.mockDb = db
	suite.mockSql = mock
	suite.repo = NewGenreToMovieRepository()
}

func TestGenreToMovieRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(GenreToMovieRepositoryTestSuite))
}

func (suite *GenreToMovieRepositoryTestSuite) TestSave_Success() {
	ctx := context.Background()
	c := gin.Context{}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`INSERT INTO genre_to_movies \(genre_id, movie_id\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(mockingGenreToMovie.GenreID, mockingGenreToMovie.MovieID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(mockingGenreToMovie.ID))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.BeginTx(ctx, nil)
	assert.NoError(suite.T(), err)

	result, err := suite.repo.Save(ctx, tx, mockingGenreToMovie, &c)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockingGenreToMovie.ID, result.ID)
}

func (suite *GenreToMovieRepositoryTestSuite) TestSave_Failed() {
	ctx := context.Background()
	c := gin.Context{}

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(`INSERT INTO genre_to_movies \(genre_id, movie_id\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(mockingGenreToMovie.GenreID, mockingGenreToMovie.MovieID).
		WillReturnError(errors.New("Insert GenreToMovie Failed"))
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.BeginTx(ctx, nil)
	assert.NoError(suite.T(), err)

	_, err = suite.repo.Save(ctx, tx, mockingGenreToMovie, &c)
	assert.Error(suite.T(), err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestFindByID_Success() {
	query := regexp.QuoteMeta(`
		SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
		       m.price as movie_price, m.duration as movie_duration, m.status as movie_status
		FROM genre_to_movies gtm 
		JOIN genres g ON gtm.genre_id = g.id 
		JOIN movies m ON gtm.movie_id = m.id
		WHERE gtm.id = $1
	`)

	genreToMovie := entity.GenreToMovie{
		ID:               "1",
		GenreID:          "genre1",
		MovieID:          "movie1",
		GenreName:        "Action",
		MovieTitle:       "Movie 1",
		MovieDescription: "Description 1",
		MoviePrice:       100,
		MovieDuration:    120,
		MovieStatus:      "active",
	}

	rows := sqlmock.NewRows([]string{
		"id", "genre_id", "movie_id", "genre_name", "movie_title", "movie_description",
		"movie_price", "movie_duration", "movie_status",
	}).AddRow(
		genreToMovie.ID, genreToMovie.GenreID, genreToMovie.MovieID, genreToMovie.GenreName,
		genreToMovie.MovieTitle, genreToMovie.MovieDescription, genreToMovie.MoviePrice,
		genreToMovie.MovieDuration, genreToMovie.MovieStatus,
	)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindByID(context.Background(), tx, "1", ginContext)
	suite.NoError(err)
	suite.Equal(genreToMovie.ID, result.ID)
	suite.Equal(genreToMovie.MovieID, result.MovieID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestFindByID_Error() {
	query := regexp.QuoteMeta(`
		SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
		       m.price as movie_price, m.duration as movie_duration, m.status as movie_status
		FROM genre_to_movies gtm 
		JOIN genres g ON gtm.genre_id = g.id 
		JOIN movies m ON gtm.movie_id = m.id
		WHERE gtm.id = $1
	`)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WithArgs("1").WillReturnError(errors.New("query error"))

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindByID(context.Background(), tx, "1", ginContext)
	suite.Error(err)
	suite.EqualError(err, "query error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestFindAll_Success() {
	query := regexp.QuoteMeta(`
		SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
		       m.price as movie_price, m.duration as movie_duration, m.status as movie_status
		FROM genre_to_movies gtm 
		JOIN genres g ON gtm.genre_id = g.id 
		JOIN movies m ON gtm.movie_id = m.id
	`)

	genreToMovies := []entity.GenreToMovie{
		{
			ID:               "1",
			GenreID:          "genre1",
			MovieID:          "movie1",
			GenreName:        "Action",
			MovieTitle:       "Movie 1",
			MovieDescription: "Description 1",
			MoviePrice:       100,
			MovieDuration:    120,
			MovieStatus:      "active",
		},
		{
			ID:               "2",
			GenreID:          "genre2",
			MovieID:          "movie2",
			GenreName:        "Comedy",
			MovieTitle:       "Movie 2",
			MovieDescription: "Description 2",
			MoviePrice:       120,
			MovieDuration:    110,
			MovieStatus:      "inactive",
		},
	}

	rows := sqlmock.NewRows([]string{
		"id", "genre_id", "movie_id", "genre_name", "movie_title", "movie_description",
		"movie_price", "movie_duration", "movie_status",
	}).AddRow(
		genreToMovies[0].ID, genreToMovies[0].GenreID, genreToMovies[0].MovieID, genreToMovies[0].GenreName,
		genreToMovies[0].MovieTitle, genreToMovies[0].MovieDescription, genreToMovies[0].MoviePrice,
		genreToMovies[0].MovieDuration, genreToMovies[0].MovieStatus,
	).AddRow(
		genreToMovies[1].ID, genreToMovies[1].GenreID, genreToMovies[1].MovieID, genreToMovies[1].GenreName,
		genreToMovies[1].MovieTitle, genreToMovies[1].MovieDescription, genreToMovies[1].MoviePrice,
		genreToMovies[1].MovieDuration, genreToMovies[1].MovieStatus,
	)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WillReturnRows(rows)

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	result, err := suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.NoError(err)
	suite.Len(result, 2)
	suite.Equal(genreToMovies[0].ID, result[0].ID)
	suite.Equal(genreToMovies[1].GenreID, result[1].GenreID)

	suite.mockSql.ExpectCommit()
	err = tx.Commit()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestFindAll_Error() {
	query := regexp.QuoteMeta(`
		SELECT gtm.id, gtm.genre_id, gtm.movie_id, g.name as genre_name, m.title as movie_title, m.description as movie_description, 
		       m.price as movie_price, m.duration as movie_duration, m.status as movie_status
		FROM genre_to_movies gtm 
		JOIN genres g ON gtm.genre_id = g.id 
		JOIN movies m ON gtm.movie_id = m.id
	`)

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectQuery(query).WillReturnError(errors.New("query error"))

	tx, err := suite.mockDb.Begin()
	suite.NoError(err)

	ginContext, _ := gin.CreateTestContext(nil)
	_, err = suite.repo.FindAll(context.Background(), tx, ginContext)
	suite.Error(err)
	suite.EqualError(err, "query error")

	suite.mockSql.ExpectRollback()
	err = tx.Rollback()
	suite.NoError(err)

	err = suite.mockSql.ExpectationsWereMet()
	suite.NoError(err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestDelete_Success() {
	ctx := context.Background()
	c := gin.Context{}
	id := mockingGenreToMovie.ID

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`DELETE FROM genre_to_movies WHERE id = \$1`).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	suite.mockSql.ExpectCommit()

	tx, err := suite.mockDb.BeginTx(ctx, nil)
	assert.NoError(suite.T(), err)

	err = suite.repo.Delete(ctx, tx, id, &c)
	assert.NoError(suite.T(), err)
}

func (suite *GenreToMovieRepositoryTestSuite) TestDelete_Failed() {
	ctx := context.Background()
	c := gin.Context{}
	id := uuid.New().String()

	suite.mockSql.ExpectBegin()
	suite.mockSql.ExpectExec(`DELETE FROM genre_to_movies WHERE id = \$1`).
		WithArgs(id).
		WillReturnError(errors.New("Delete Failed"))
	suite.mockSql.ExpectRollback()

	tx, err := suite.mockDb.BeginTx(ctx, nil)
	assert.NoError(suite.T(), err)

	err = suite.repo.Delete(ctx, tx, id, &c)
	assert.Error(suite.T(), err)
}
