package repository

import (
	"database/sql"
	"github.com/borscht/backend/database/local"
	"github.com/borscht/backend/internal/models"
	"github.com/borscht/backend/internal/user"
	_errors "github.com/borscht/backend/utils"
	"io"
	"mime/multipart"
	"os"
)

type userRepo struct {
	db local.Database
	DB *sql.DB
}

func NewUserRepo() user.UserRepo {
	return &userRepo{
		db: local.GetInstance(),
	}
}

func (u *userRepo) Create(newUser models.User) (int32, error) {
	//for _, curUser := range *u.db.GetModels().Users {
	//	if (curUser.Phone == newUser.Phone) && curUser.Password == newUser.Password {
	//		return _errors.NewCustomError(http.StatusUnauthorized, "User with this number already exists") // такой юзер уже есть
	//	}
	//}
	//
	//// записываем нового
	//*u.db.GetModels().Users = append(*u.db.GetModels().Users, newUser)
	//return nil
	// todo exists checking
	var uid int32
	err := u.DB.QueryRow("insert into users (name, number, email, password) values ($1, $2, $3, $4) returning uid",
		newUser.Name, newUser.Phone, newUser.Email, newUser.Password).Scan(&uid)
	if err != nil {
		// TODO
	}

	return uid, nil
}

func (u *userRepo) CheckUserExists(userToCheck models.UserAuth) (models.User, error) {
	//for _, curUser := range *u.db.GetModels().Users {
	//	if (curUser.Email == check.Login || curUser.Phone == check.Login) && curUser.Password == check.Password {
	//		return curUser, nil
	//	}
	//}
	//
	//return models.User{}, _errors.Authorization("not curUser bd")
	DBuser, err := u.DB.Query("select uid, name, avatar from users where (phone=&1 or email=$1) and password=$2",
		userToCheck.Login, userToCheck.Password)
	if err != nil {
		// todo
	}
	user := new(models.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user.Uid,
			&user.Name,
			&user.Avatar,
		)
		if err != nil {
			// TODO
		}
	}

	return *user, nil
}

func (u *userRepo) GetByUid(uid int32) (models.User, error) {
	//for _, curUser := range *u.db.GetModels().Users {
	//	if curUser.Phone == number {
	//		return curUser, nil
	//	}
	//}
	//
	//return models.User{}, _errors.Authorization("curUser not found")

	DBuser, err := u.DB.Query("select name, phone, email, avatar, address, mainAddress from users where uid=$1", uid)
	if err != nil {
		return models.User{}, _errors.Authorization("user not found")
	}
	user := new(models.User)
	for DBuser.Next() {
		err = DBuser.Scan(
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.MainAddress,
			&user.Avatar,
		)
		if err != nil {
			// TODO
		}
	}
	return *user, nil
}

func (u *userRepo) Update(newUser models.UserData) error {
	// TODO

	//for i, curUser := range *u.db.GetModels().Users {
	//	if curUser.Email == newUser.Email && curUser.Phone != ctx.Value("User").(models.User).Phone { // если у кого-то другого уже есть такой email
	//		return _errors.NewCustomError(http.StatusBadRequest, "curUser with this email already exists")
	//	}
	//	if curUser.Phone == newUser.Phone && curUser.Phone != ctx.Value("User").(models.User).Phone { // если у кого-то другого уже есть такой телефон
	//		return _errors.NewCustomError(http.StatusBadRequest, "User with this number already exists")
	//	}
	//
	//	if curUser.Phone == ctx.Value("User").(models.User).Phone {
	//		if newUser.Password != "" {
	//			if newUser.PasswordOld != curUser.Password {
	//				fmt.Println(newUser.PasswordOld, " ", curUser.Password)
	//				return _errors.NewCustomError(http.StatusBadRequest, "invalid old password")
	//			}
	//			(*u.db.GetModels().Users)[i].Password = newUser.Password
	//		}
	//
	//		(*u.db.GetModels().Users)[i].Phone = newUser.Phone
	//		(*u.db.GetModels().Users)[i].Email = newUser.Email
	//		(*u.db.GetModels().Users)[i].Name = newUser.Name
	//
	//		fmt.Println(*u.db.GetModels().Users)
	//		//return cc.SendResponse(profileEdits)
	//		return nil
	//	}
	//}
	//
	//return _errors.Authorization("curUser not found")
}

func (u *userRepo) UploadAvatar(image *multipart.FileHeader, filename string) error {
	// Читаем файл из пришедшего запроса
	src, err := image.Open()
	if err != nil {
		return _errors.FailServer(err.Error())
	}
	defer src.Close()

	// создаем файл у себя
	dst, err := os.Create(filename)
	if err != nil {
		return _errors.FailServer(err.Error())
	}
	defer dst.Close()

	// копируем один в другой
	if _, err = io.Copy(dst, src); err != nil {
		return _errors.FailServer(err.Error())
	}

	return nil
}
