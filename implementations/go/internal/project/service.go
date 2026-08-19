package project

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func ListProjects() {}

func CreateProject() {}

func GetProjectByID() {}

func UpdateProject() {}

func StartProject() {}

func ListProgresses() {}

func GetProgressByID() {}

func UpdateProgress() {}

func UpdateCompletedItems() {}
