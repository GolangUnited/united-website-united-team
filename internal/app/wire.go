//go:build wireinject

package app

import (
	"context"
	_ "github.com/google/subcommands"
	"github.com/google/wire"
	httpAdapter "github.com/zhuravlev-pe/course-watch/internal/adapter/http"
	httpV1 "github.com/zhuravlev-pe/course-watch/internal/adapter/http/v1"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/repository"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/internal/app/config"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
)

var commonSet = wire.NewSet(
	//config.GetConfig,
	createIdGen,
)

var repoSet = wire.NewSet(
	repository.NewUsersRepo,
	wire.Bind(new(service.UsersRepository), new(*repository.UsersRepo)),
	fake_repo.NewCourses,
	//wire.Bind(new(service.CoursesRepository), new(*fake_repo.CoursesRepo)),
	//wire.Struct(new(service.Repositories), "*"),
)

var usersServiceSet = wire.NewSet(
	service.NewUsersService,
	wire.Bind(new(httpV1.UserService), new(*service.UserService)),
)

var coursesServiceSet = wire.NewSet(
	service.NewCoursesService,
	wire.Bind(new(httpV1.CourseService), new(*service.CourseService)),
)

func injectHandler(ctx context.Context, cfg *config.Config) (*httpAdapter.Handler, func(), error) {
	wire.Build(commonSet, createPgClient, createAuthenticator, repoSet, usersServiceSet, coursesServiceSet,
		httpAdapter.NewHandler,
	)
	return nil, nil, nil
}
