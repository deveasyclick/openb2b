## 📖 API Documentation (Swagger)

OpenB2B provides a **Swagger UI** to explore and test the API endpoints interactively.

---

### Accessing Swagger UI

If your host is running on `localhost:3000`, you can access the Swagger UI at [http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html).


- This serves the **Swagger UI**, which allows you to view all available API endpoints, request parameters, and response formats.  
- The UI dynamically fetches the API specification JSON from:

[http://localhost:3000/swagger/doc.json](http://localhost:3000/swagger/doc.json)

> Make sure the backend server is running on the port specified in your `.env` file (`PORT`).

---

### Generating Swagger Documentation

If you add new endpoints or update comments, regenerate the Swagger docs using:

```bash
# From the root directory
make swagger
```

This will generate or update the `docs` folder containing `doc.json` and `docs.go`.

---

### Swagger Annotations

- OpenB2B uses `swaggo/swag` annotations in your Go code to generate the documentation.
- Example:

```go
// @title OpenB2B API
// @version 1.0
// @description OpenB2B API documentation
// @host localhost:8080
// @BasePath /api/v1
```

- Each endpoint can have comments like:

```go
// @Summary Get all users
// @Description Retrieve a list of users in the system
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Router /users [get]
```

---

### Notes

- The Swagger **host URL** is dynamically set in the backend using your configured `PORT`.
- You can use Swagger UI to **try requests directly from the browser**.
- If running the frontend on a separate port (e.g., `3000`), make sure to access Swagger directly via the backend port to avoid CORS issues.

---

### References

- [Swagger UI for Go (`swaggo/http-swagger`)](https://github.com/swaggo/http-swagger)
- [Swaggo annotations documentation](https://github.com/swaggo/swag#general-api-info)
