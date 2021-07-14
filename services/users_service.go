package services

import (
	"github.com/nicoletafratila/bookstore_users-api/domain/users"
	"github.com/nicoletafratila/bookstore_users-api/utils/crypto_utils"
	"github.com/nicoletafratila/bookstore_users-api/utils/date_utils"
	"github.com/nicoletafratila/bookstore_utils-go/rest_errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {
}

type usersServiceInterface interface {
	Get(int64) (*users.User, rest_errors.RestErr)
	Create(users.User) (*users.User, rest_errors.RestErr)
	Update(bool, users.User) (*users.User, rest_errors.RestErr)
	Delete(int64) rest_errors.RestErr
	Search(string) (users.Users, rest_errors.RestErr)
	Login(request users.LoginRequest) (*users.User, rest_errors.RestErr)
}

func (s *usersService) Get(userId int64) (*users.User, rest_errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) Create(user users.User) (*users.User, rest_errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDbFormat()
	user.Password = crypto_utils.GetMd5(user.Password)

	if err := user.Create(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) Update(isPartial bool, user users.User) (*users.User, rest_errors.RestErr) {
	current, err := s.Get(user.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *usersService) Delete(userId int64) rest_errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) Search(status string) (users.Users, rest_errors.RestErr) {
	dao := &users.User{}
	return dao.SearchByStatus(status)
}

func (s *usersService) Login(request users.LoginRequest) (*users.User, rest_errors.RestErr) {
	dao := &users.User{
		Email:    request.Email,
		Password: crypto_utils.GetMd5(request.Password),
	}
	if err := dao.SearchByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
