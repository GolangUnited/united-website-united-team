package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zhuravlev-pe/course-watch/internal/adapter/repository/fake_repo"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
)

func NewRepositories(client *pgxpool.Pool) *service.Repositories {
	return &service.Repositories{
		Users:   NewUsersRepo(client),
		Courses: fake_repo.NewCourses(),
	}
}
