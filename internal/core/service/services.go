package service

type Services struct {
	Users   *UserService
	Courses *CourseService
}

func NewServices(deps Deps) *Services {
	return &Services{
		Users:   NewUsersService(deps.Repos.Users, deps.IdGen),
		Courses: NewCoursesService(deps.Repos.Courses, deps.IdGen),
	}
}
