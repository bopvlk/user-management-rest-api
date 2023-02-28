package mappers

import (
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"github.com/labstack/echo/v4"
	"strconv"
)

func MapContextToPagination(c echo.Context) *models.Pagination {

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 5
		c.Logger().Error(err)
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		page = 1
		c.Logger().Error(err)
	}

	sort := c.QueryParam("sort")
	if sort == "" {
		sort = "id desc"
	}

	return &models.Pagination{Limit: limit, Page: page, Sort: sort}
}
