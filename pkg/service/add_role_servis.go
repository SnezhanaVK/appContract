package service

import (
	db "appContract/pkg/db/repository"
	"appContract/pkg/models"
	"errors"
)

func AddRole(id_user int, id_role int) error{
	if id_role == 1 {
		return db.DBaddUserAdmin(models.Users{Id_user: id_user})

	} else if id_role == 2 {
		return db.DBaddUserMeneger(models.Users{Id_user: id_user})
	} else if id_role == 3 {
		return db.DBaddUserUser(models.Users{Id_user: id_user})
	}else{
		return errors.New("role not found")
	}

}


