package main

import (
	"backend/api/image"
	"github.com/labstack/echo/v4"
	"net/http"
)

// тут будет загрузка главной страницы с ресторанами
func mainPage(c echo.Context) error {
	return c.HTML(http.StatusOK, `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Single file upload</title>
</head>
<body>
<h1>Upload single file with fields</h1>

<form action="/upload" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="avatar"><br><br>
    <input type="submit" value="Submit">
</form>
</body>
</html>
`)
	//return c.String(http.StatusOK, "It will be the main page")
}

func router(e *echo.Echo) {
	e.GET("/", mainPage)
	//e.GET("/:id", restaurantPage) // урл на получение странички ресторана номер id
	//e.POST("/signup", createUser) // урл на регистрацию пользователя
	//e.POST("/signin", logUser) // урл на авторизацию
	//e.POST("/edituser", updateUser) // обновить пользователя после редактирования профиля
	e.POST("/upload", image.UploadAvatar)
	e.GET("/avatar", image.DownloadAvatar)
}

func main() {
	//users := make([]user, 0, 0)
	//restaurants := make([]restaurant, 0, 0)

	e := echo.New()
	router(e)

	e.Logger.Fatal(e.Start(":5000"))
}

