package project

import (
	"baselayer/internal/auth"
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

func (s *Service) ListProjects(ctx context.Context) ([]Project, error) {
	filter := bson.D{}
	sort := bson.D{{Key: "createdAt", Value: -1}}
	return s.repo.SearchProjects(ctx, filter, sort)
}

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

func (s *Service) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return project, err
}

func UpdateProject() {}

func (s *Service) StartProject(ctx context.Context, id string) (string, error) {
	rawUserID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return "", err
	}
	userID, err := bson.ObjectIDFromHex(rawUserID)
	if err != nil {
		return "", err
	}
	project, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		return "", err
	}
	now := time.Now()
	progress := ProjectProgress{
		ID:          bson.NewObjectID(),
		UserID:      userID,
		ProjectID:   project.ID,
		Title:       project.Title,
		Description: project.Description,
		Progress:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return s.repo.CreateProgress(ctx, progress)
}

func ListProgresses() {}

func GetProgressByID() {}

func UpdateProgress() {}

func UpdateCompletedItems() {}
