package handler

import (
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) Swagger(c fiber.Ctx) error {

	return c.Type("html").SendString(`
<!DOCTYPE html>
<html>
<head>
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
</head>
<body>
  <div id="swagger-ui"></div>

  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
  <script>
    const ui = SwaggerUIBundle({
      url: 'docs/swagger.json',
      dom_id: '#swagger-ui',
    });
  </script>
</body>
</html>
`)
}
