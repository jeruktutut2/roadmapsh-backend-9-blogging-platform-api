package services

import (
	modelentities "blogging-platform-api/models/entities"
	modelrequests "blogging-platform-api/models/requests"
	modelresponses "blogging-platform-api/models/responses"
	"blogging-platform-api/repositories"
	"blogging-platform-api/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type BlogService interface {
	Create(ctx context.Context, createRequest modelrequests.CreateRequest) (httpCode int, response interface{})
	Update(ctx context.Context, idBlog int, updateRequest modelrequests.UpdateRequest) (httpCode int, response interface{})
	Delete(ctx context.Context, idBlog int) (httpCode int, response interface{})
	FindById(ctx context.Context, idBlog int) (httpCode int, response interface{})
	FindAllPosts(ctx context.Context, term string) (httpCode int, response interface{})
	// Filter()
}

type BlogServiceImplementation struct {
	PostgresUtil   utils.PostgresUtil
	Validate       *validator.Validate
	BlogRepository repositories.BlogRepository
}

func NewBlogService(postgresUtil utils.PostgresUtil, validate *validator.Validate, blogRepository repositories.BlogRepository) BlogService {
	return &BlogServiceImplementation{
		PostgresUtil:   postgresUtil,
		Validate:       validate,
		BlogRepository: blogRepository,
	}
}

func (service *BlogServiceImplementation) Create(ctx context.Context, createRequest modelrequests.CreateRequest) (httpCode int, response interface{}) {
	err := service.Validate.Struct(createRequest)
	if err != nil {
		httpCode = http.StatusBadRequest
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	var blog modelentities.Blog
	blog.Title = pgtype.Text{Valid: true, String: createRequest.Title}
	blog.Content = pgtype.Text{Valid: true, String: createRequest.Content}
	blog.Category = pgtype.Text{Valid: true, String: createRequest.Category}
	blog.Tags = pgtype.Text{Valid: true, String: strings.Join(createRequest.Tags, ", ")}
	blog.CreatedAt = pgtype.Int8{Valid: true, Int64: time.Now().UnixMilli()}
	blog.UpdatedAt = pgtype.Int8{Valid: true, Int64: time.Now().UnixMilli()}
	insertedId, err := service.BlogRepository.Create(service.PostgresUtil.GetPool(), ctx, blog)
	if err != nil {
		httpCode = http.StatusBadRequest
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	var createResponse modelresponses.CreateResponse
	createResponse.Id = insertedId
	createResponse.Title = createRequest.Title
	createResponse.Content = createRequest.Content
	createResponse.Category = createRequest.Category
	createResponse.Tags = createRequest.Tags
	createResponse.CreatedAt = time.Unix(blog.CreatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
	createResponse.UpdatedAt = time.Unix(blog.UpdatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")

	httpCode = http.StatusCreated
	response = createResponse
	return
}

func (service *BlogServiceImplementation) Update(ctx context.Context, idBlog int, updateRequest modelrequests.UpdateRequest) (httpCode int, response interface{}) {
	err := service.Validate.Struct(updateRequest)
	if err != nil {
		httpCode = http.StatusBadRequest
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	var blog modelentities.Blog
	blog.Id = pgtype.Int4{Valid: true, Int32: int32(idBlog)}
	blog.Title = pgtype.Text{Valid: true, String: updateRequest.Title}
	blog.Content = pgtype.Text{Valid: true, String: updateRequest.Content}
	blog.Category = pgtype.Text{Valid: true, String: updateRequest.Category}
	blog.Tags = pgtype.Text{Valid: true, String: strings.Join(updateRequest.Tags, ", ")}
	blog.UpdatedAt = pgtype.Int8{Valid: true, Int64: time.Now().UnixMilli()}
	rowsAffected, err := service.BlogRepository.Update(service.PostgresUtil.GetPool(), ctx, blog)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	if rowsAffected != 1 {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse("rows affeted update not one")
		return
	}
	blog, err = service.BlogRepository.FindById(service.PostgresUtil.GetPool(), ctx, idBlog)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	var updateResponse modelresponses.UpdateResponse
	updateResponse.Id = int(blog.Id.Int32)
	updateResponse.Title = blog.Title.String
	updateResponse.Content = blog.Content.String
	updateResponse.Category = blog.Category.String
	updateResponse.Tags = strings.Split(blog.Tags.String, ",")
	updateResponse.CreatedAt = time.Unix(blog.CreatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
	updateResponse.UpdatedAt = time.Unix(blog.UpdatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
	httpCode = http.StatusOK
	response = updateResponse
	return
}

func (service *BlogServiceImplementation) Delete(ctx context.Context, idBlog int) (httpCode int, response interface{}) {
	_, err := service.BlogRepository.FindById(service.PostgresUtil.GetPool(), ctx, idBlog)
	if err != nil && err != pgx.ErrNoRows {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	} else if err == pgx.ErrNoRows {
		httpCode = http.StatusNotFound
		response = modelresponses.ToErrorResponse("not found")
		return
	}
	rowsAffected, err := service.BlogRepository.Delete(service.PostgresUtil.GetPool(), ctx, idBlog)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	if rowsAffected != 1 {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse("rows affected not one")
		return
	}
	httpCode = http.StatusNoContent
	response = ""
	return
}

func (service *BlogServiceImplementation) FindById(ctx context.Context, idBlog int) (httpCode int, response interface{}) {
	blog, err := service.BlogRepository.FindById(service.PostgresUtil.GetPool(), ctx, idBlog)
	if err != nil && err != pgx.ErrNoRows {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	} else if err == pgx.ErrNoRows {
		httpCode = http.StatusBadRequest
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}
	var findByIdResponse modelresponses.FindByIdResponse
	findByIdResponse.Id = int(blog.Id.Int32)
	findByIdResponse.Title = blog.Title.String
	findByIdResponse.Content = blog.Content.String
	findByIdResponse.Category = blog.Category.String
	findByIdResponse.Tags = strings.Split(blog.Tags.String, ",")
	findByIdResponse.CreatedAt = time.Unix(blog.CreatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
	findByIdResponse.UpdatedAt = time.Unix(blog.UpdatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
	httpCode = http.StatusOK
	response = findByIdResponse
	return
}

func (service *BlogServiceImplementation) FindAllPosts(ctx context.Context, term string) (httpCode int, response interface{}) {
	blogs, err := service.BlogRepository.FindAll(service.PostgresUtil.GetPool(), ctx, term)
	if err != nil {
		httpCode = http.StatusInternalServerError
		response = modelresponses.ToErrorResponse(err.Error())
		return
	}

	var findAll []modelresponses.FindResponse
	for _, blog := range blogs {
		var findResponse modelresponses.FindResponse
		findResponse.Id = int(blog.Id.Int32)
		findResponse.Title = blog.Title.String
		findResponse.Content = blog.Content.String
		findResponse.Category = blog.Category.String
		findResponse.Tags = strings.Split(blog.Tags.String, ",")
		findResponse.CreatedAt = time.Unix(blog.CreatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
		findResponse.UpdatedAt = time.Unix(blog.UpdatedAt.Int64/1000, 0).UTC().Format("2006-01-02T15:04:05Z")
		findAll = append(findAll, findResponse)
	}
	httpCode = http.StatusOK
	response = findAll
	return
}
