package repositories

import (
	modelentities "blogging-platform-api/models/entities"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogRepository interface {
	Create(pool *pgxpool.Pool, ctx context.Context, blog modelentities.Blog) (insertedId int, err error)
	Update(pool *pgxpool.Pool, ctx context.Context, blog modelentities.Blog) (rowsAffected int64, err error)
	FindById(pool *pgxpool.Pool, ctx context.Context, id int) (blog modelentities.Blog, err error)
	Delete(pool *pgxpool.Pool, ctx context.Context, id int) (rowsAffected int64, err error)
	FindAll(pool *pgxpool.Pool, ctx context.Context, term string) (blogs []modelentities.Blog, err error)
}

type BlogRepositoryImplementation struct {
}

func NewBlogRepository() BlogRepository {
	return &BlogRepositoryImplementation{}
}

func (repository *BlogRepositoryImplementation) Create(pool *pgxpool.Pool, ctx context.Context, blog modelentities.Blog) (insertedId int, err error) {
	query := `INSERT INTO blogs (title,content,category,tags,created_at,updated_at) 
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id;`
	err = pool.QueryRow(ctx, query, blog.Title, blog.Content, blog.Category, blog.Tags, blog.CreatedAt, blog.UpdatedAt).Scan(&insertedId)
	return
}

func (repository *BlogRepositoryImplementation) Update(pool *pgxpool.Pool, ctx context.Context, blog modelentities.Blog) (rowsAffected int64, err error) {
	query := `UPDATE blogs SET title = $1, content = $2, category = $3, tags = $4, updated_at = $5 WHERE id = $6;`
	result, err := pool.Exec(ctx, query, blog.Title, blog.Content, blog.Category, blog.Tags, blog.UpdatedAt, blog.Id)
	if err != nil {
		return
	}
	rowsAffected = result.RowsAffected()
	return
}

func (repository *BlogRepositoryImplementation) FindById(pool *pgxpool.Pool, ctx context.Context, id int) (blog modelentities.Blog, err error) {
	query := `SELECT id,title,content,category,tags,created_at,updated_at FROM blogs WHERE id = $1;`
	err = pool.QueryRow(ctx, query, id).Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Category, &blog.Tags, &blog.CreatedAt, &blog.UpdatedAt)
	return
}

func (repository *BlogRepositoryImplementation) Delete(pool *pgxpool.Pool, ctx context.Context, id int) (rowsAffected int64, err error) {
	query := `DELETE FROM blogs WHERE id = $1;`
	result, err := pool.Exec(ctx, query, id)
	if err != nil {
		return
	}
	rowsAffected = result.RowsAffected()
	return
}

func (repository *BlogRepositoryImplementation) FindAll(pool *pgxpool.Pool, ctx context.Context, term string) (blogs []modelentities.Blog, err error) {
	whereTerm := ""
	params := []interface{}{}
	if term != "" {
		whereTerm = " WHERE tags ILIKE '%' || $1 || '%'"
		params = append(params, term)
	}
	query := `SELECT id,title,content,category,tags,created_at,updated_at FROM blogs ` + whereTerm + `;`
	rows, err := pool.Query(ctx, query, params...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var blog modelentities.Blog
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Category, &blog.Tags, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			blogs = []modelentities.Blog{}
			return
		}
		blogs = append(blogs, blog)
	}
	if rows.Err() != nil {
		blogs = []modelentities.Blog{}
		err = rows.Err()
		return
	}
	return
}
