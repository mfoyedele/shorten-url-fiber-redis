package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mfoyedele/shorten-url-fiber-redis/database"
	"github.com/mfoyedele/shorten-url-fiber-redis/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit`
	XRateLimitReset time.Duration `json:"rate_limit_reset`
}

func ShortenURL(c *fiber.Ctx) {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	//implement rate limiting

	r2 := database.CreateClient(1)

	//check if the input is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Inavlid URL"})
	}

	//check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "you can't hack the system(:"})
	}

	//enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

}
