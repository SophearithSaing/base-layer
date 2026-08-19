package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func ListProjects() {}

func (s *Service) CreateProject(ctx context.Context, payload Project) error {
	project := Project{
		ID:                 bson.NewObjectID(),
		Title:              payload.Title,
		Description:        payload.Description,
		Legend:             payload.Legend,
		Phases:             payload.Phases,
		Capstones:          payload.Capstones,
		RecommendedOrder:   payload.RecommendedOrder,
		MasteryDefinitions: payload.MasteryDefinitions,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	return s.repo.CreateProject(ctx, project)
}

func GetProjectByID() {}

func UpdateProject() {}

func StartProject() {}

func ListProgresses() {}

func GetProgressByID() {}

func UpdateProgress() {}

func UpdateCompletedItems() {}
