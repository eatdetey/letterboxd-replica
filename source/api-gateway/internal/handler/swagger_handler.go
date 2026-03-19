package handler

import "github.com/gofiber/fiber/v3"

const openAPISpecYAML = `openapi: 3.0.3
info:
  title: Letterboxd Replica API Gateway
  version: 1.0.0
  description: HTTP gateway over gRPC services (user, movie, review).
servers:
  - url: /api
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    refreshCookie:
      type: apiKey
      in: cookie
      name: refresh_token
paths:
  /v1/auth/register:
    post:
      summary: Register user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [username, password, email]
              properties:
                username: { type: string }
                password: { type: string }
                email: { type: string }
      responses:
        "200":
          description: User + access token (refresh token is set in cookie)
  /v1/auth/login:
    post:
      summary: Login user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [username, password]
              properties:
                username: { type: string }
                password: { type: string }
      responses:
        "200":
          description: User + access token (refresh token is set in cookie)
  /v1/auth/refresh:
    post:
      summary: Refresh access token by refresh cookie
      security:
        - refreshCookie: []
      responses:
        "200":
          description: New access token in response body and new refresh cookie
  /v1/users:
    get:
      summary: Get users by ids or usernames
      parameters:
        - in: query
          name: ids
          schema: { type: string }
        - in: query
          name: usernames
          schema: { type: string }
        - in: query
          name: limit
          schema: { type: integer }
        - in: query
          name: offset
          schema: { type: integer }
      responses:
        "200":
          description: Users list
  /v1/movies:
    get:
      summary: Get movies
      parameters:
        - in: query
          name: limit
          schema: { type: integer }
        - in: query
          name: offset
          schema: { type: integer }
        - in: query
          name: search
          schema: { type: string }
        - in: query
          name: genre
          schema: { type: string }
        - in: query
          name: playlist_id
          schema: { type: string }
        - in: query
          name: enrich_playlists
          schema: { type: boolean }
      responses:
        "200":
          description: Movies list
  /v1/movies/{id}/reviews:
    get:
      summary: Get movie reviews
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
      responses:
        "200":
          description: Reviews list
    post:
      summary: Add review to movie
      security:
        - bearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [text]
              properties:
                text: { type: string }
      responses:
        "201":
          description: Created review
  /v1/playlists:
    get:
      summary: Get current user playlists
      security:
        - bearerAuth: []
      responses:
        "200":
          description: Playlists
    post:
      summary: Create playlist
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name: { type: string }
      responses:
        "201":
          description: Created playlist
  /v1/playlists/{id}:
    put:
      summary: Rename playlist
      security:
        - bearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name: { type: string }
      responses:
        "200":
          description: Renamed playlist
    delete:
      summary: Delete playlist
      security:
        - bearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
      responses:
        "204":
          description: Deleted
  /v1/playlists/{id}/movies:
    post:
      summary: Add movie to playlist
      security:
        - bearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [movie_id]
              properties:
                movie_id: { type: string }
      responses:
        "204":
          description: Added
  /v1/playlists/{id}/movies/{movie_id}:
    delete:
      summary: Remove movie from playlist
      security:
        - bearerAuth: []
      parameters:
        - in: path
          name: id
          required: true
          schema: { type: string }
        - in: path
          name: movie_id
          required: true
          schema: { type: string }
      responses:
        "204":
          description: Removed
`

const swaggerUIHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1"/>
  <title>API Gateway Swagger</title>
  <style>
    body { margin: 0; font-family: sans-serif; }
    #swagger-ui { min-height: 100vh; }
    .swagger-fallback { padding: 16px; }
    .swagger-fallback pre {
      white-space: pre-wrap;
      word-break: break-word;
      background: #f5f5f5;
      border-radius: 8px;
      padding: 12px;
      overflow: auto;
    }
  </style>
</head>
<body>
  <div id="swagger-ui">Loading Swagger UI...</div>
  <script>
    (function() {
      const root = document.getElementById("swagger-ui");
      const cssCandidates = [
        "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css",
        "https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"
      ];
      const jsCandidates = [
        "https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js",
        "https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"
      ];

      function loadCSS(src) {
        return new Promise((resolve, reject) => {
          const link = document.createElement("link");
          link.rel = "stylesheet";
          link.href = src;
          link.onload = resolve;
          link.onerror = reject;
          document.head.appendChild(link);
        });
      }

      function loadScript(src) {
        return new Promise((resolve, reject) => {
          const script = document.createElement("script");
          script.src = src;
          script.onload = resolve;
          script.onerror = reject;
          document.head.appendChild(script);
        });
      }

      function tryLoad(list, loader) {
        return new Promise((resolve, reject) => {
          let idx = 0;
          const next = () => {
            if (idx >= list.length) {
              reject(new Error("all candidates failed"));
              return;
            }
            const current = list[idx++];
            loader(current).then(resolve).catch(next);
          };
          next();
        });
      }

      function escapeHTML(text) {
        return text
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;");
      }

      function renderFallback(message) {
        fetch("/swagger/openapi.yaml")
          .then((r) => r.text())
          .then((spec) => {
            root.innerHTML =
              "<div class='swagger-fallback'>" +
              "<h2>Swagger UI недоступен</h2>" +
              "<p>" + message + "</p>" +
              "<p>Ниже OpenAPI спецификация:</p>" +
              "<pre>" + escapeHTML(spec) + "</pre>" +
              "</div>";
          })
          .catch(() => {
            root.innerHTML =
              "<div class='swagger-fallback'>" +
              "<h2>Swagger UI недоступен</h2>" +
              "<p>" + message + "</p>" +
              "<p>И OpenAPI спецификацию тоже не удалось загрузить.</p>" +
              "</div>";
          });
      }

      tryLoad(cssCandidates, loadCSS)
        .catch(() => {})
        .finally(() => {
          tryLoad(jsCandidates, loadScript)
            .then(() => {
              if (typeof SwaggerUIBundle !== "function") {
                renderFallback("Скрипт Swagger UI загружен, но не инициализировался.");
                return;
              }
              SwaggerUIBundle({
                url: "/swagger/openapi.yaml",
                dom_id: "#swagger-ui",
                deepLinking: true,
                displayRequestDuration: true
              });
            })
            .catch(() => {
              renderFallback("Не удалось загрузить Swagger UI assets с CDN.");
            });
        });
    })();
  </script>
</body>
</html>`

func SwaggerSpec(c fiber.Ctx) error {
	c.Type("yaml", "utf-8")
	return c.SendString(openAPISpecYAML)
}

func SwaggerUI(c fiber.Ctx) error {
	c.Type("html", "utf-8")
	return c.SendString(swaggerUIHTML)
}
