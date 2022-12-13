package core

type Course struct {
	Id          string
	Title       string
	Description string
}

type CreateCourseInput struct {
	Title       string
	Description string
}
