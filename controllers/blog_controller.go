package controllers

import (
	modelrequests "blogging-platform-api/models/requests"
	"blogging-platform-api/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BlogController interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	FindById(c echo.Context) error
	FindAll(c echo.Context) error
}

type BlogControllerImplementation struct {
	BlogService services.BlogService
}

func NewBlogController(blogService services.BlogService) BlogController {
	return &BlogControllerImplementation{
		BlogService: blogService,
	}
}

func (controller *BlogControllerImplementation) Create(c echo.Context) error {
	var createRequest modelrequests.CreateRequest
	err := c.Bind(&createRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.BlogService.Create(c.Request().Context(), createRequest)
	return c.JSON(httpCode, response)
}

func (controller *BlogControllerImplementation) Update(c echo.Context) error {
	var updateRequest modelrequests.UpdateRequest
	err := c.Bind(&updateRequest)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.BlogService.Update(c.Request().Context(), id, updateRequest)
	return c.JSON(httpCode, response)
}

func (controller *BlogControllerImplementation) Delete(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, _ := controller.BlogService.Delete(c.Request().Context(), id)
	return c.NoContent(httpCode)
}

func (controller *BlogControllerImplementation) FindById(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	httpCode, response := controller.BlogService.FindById(c.Request().Context(), id)
	return c.JSON(httpCode, response)
}

func (controller *BlogControllerImplementation) FindAll(c echo.Context) error {
	term := c.QueryParam("term")
	httpCode, response := controller.BlogService.FindAllPosts(c.Request().Context(), term)
	return c.JSON(httpCode, response)
}
