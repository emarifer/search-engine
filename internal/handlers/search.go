package handlers

import (
	"github.com/emarifer/search-engine/internal/handlers/dto"
	"github.com/emarifer/search-engine/internal/services"
	"github.com/gofiber/fiber/v2"
)

/********** Handlers for Search endpoint **********/

type SearchService interface {
	SearchFullText(v string) ([]services.CrawledUrl, error)
}

func NewSearchHandler(s SearchService) SearchHandler {

	return SearchHandler{
		Search: s,
	}
}

type SearchHandler struct {
	Search SearchService
}

func (sh *SearchHandler) searchHandler(c *fiber.Ctx) error {
	search := dto.SearchJsonDto{}
	err := c.BodyParser(&search)
	if err != nil || search.Term == "" {
		c.Status(fiber.StatusBadRequest)
		c.Append("content-type", "application/json")

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid input",
			"data":    nil,
		})
	}

	data, err := sh.Search.SearchFullText(search.Term)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Append("content-type", "application/json")

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Could not find a response",
			"data":    nil,
		})
	}

	c.Status(fiber.StatusOK)
	c.Append("content-type", "application/json")

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Search results",
		"data":    data,
	})
}
