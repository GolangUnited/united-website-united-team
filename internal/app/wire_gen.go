// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"context"
	"github.com/google/wire"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/http"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/repository"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/internal/app/config"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
)

import (
	_ "github.com/google/subcommands"
	_ "github.com/joho/godotenv/autoload"
)

// Injectors from wire.go:

func injectHandler(ctx context.Context, cfg *config.Config) (*http.Handler, func(), error) {
	pool, cleanup, err := createPgClient(cfg, ctx)
	if err != nil {
		return nil, nil, err
	}
	usersRepo := repository.NewUsersRepo(pool)
	idGenerator, err := createIdGen(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userService := service.NewUsersService(usersRepo, idGenerator)
	coursesRepository := fake_repo.NewCourses()
	courseService := service.NewCoursesService(coursesRepository, idGenerator)
	bearerAuthenticator, err := createAuthenticator(cfg)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	handler := v1.NewHandler(userService, courseService, bearerAuthenticator)
	httpHandler := http.NewHandler(handler)
	return httpHandler, func() {
		cleanup()
	}, nil
}

// wire.go:

var commonSet = wire.NewSet(

	createIdGen,
)

var repoSet = wire.NewSet(repository.NewUsersRepo, wire.Bind(new(service.UsersRepository), new(*repository.UsersRepo)), fake_repo.NewCourses)

var usersServiceSet = wire.NewSet(service.NewUsersService, wire.Bind(new(v1.UserService), new(*service.UserService)))

var coursesServiceSet = wire.NewSet(service.NewCoursesService, wire.Bind(new(v1.CourseService), new(*service.CourseService)))
