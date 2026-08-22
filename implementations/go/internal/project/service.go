package project

import (
	"baselayer/internal/auth"
	"context"
	"maps"
	"math"
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

func (s *Service) UpdateProject(ctx context.Context, id string, payload UpdateProjectPayload) (*Project, error) {
	update := UpdatePayload[UpdateProjectPayload]{
		Payload:   payload,
		UpdatedAt: time.Now(),
	}
	return s.repo.UpdateProject(ctx, id, bson.M{"$set": update})
}

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

func (s *Service) ListProgresses(ctx context.Context) (*[]ProjectProgress, error) {
	rawUserID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := bson.ObjectIDFromHex(rawUserID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "userId", Value: userID}}
	sort := bson.D{{Key: "createdAt", Value: -1}}
	return s.repo.SearchProgresses(ctx, filter, sort)
}

func (s *Service) GetProgressByID(ctx context.Context, id string) (*ProjectProgress, error) {
	userID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	progress, err := s.repo.GetProgressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if progress.UserID.Hex() != userID {
		return nil, ErrUserDontHavePermissionToView
	}
	return progress, nil
}

func (s *Service) UpdateProgress(ctx context.Context, id string, payload UpdateProjectProgressPayload) (*ProjectProgress, error) {
	rawUserID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	userID, err := bson.ObjectIDFromHex(rawUserID)
	if err != nil {
		return nil, err
	}
	progressID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id":    progressID,
		"userId": userID,
	}
	update := UpdatePayload[UpdateProjectProgressPayload]{
		Payload:   payload,
		UpdatedAt: time.Now(),
	}
	return s.repo.UpdateProgress(ctx, filter, bson.M{"$set": update})
}

func (s *Service) UpdateCompletedItems(ctx context.Context, id string, payload UpdateCompletedItemsPayload) (*ProjectProgress, error) {
	rawUserID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return nil, err
	}
	progressID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	progress, err := s.repo.GetProgressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if progress.UserID.Hex() != rawUserID {
		return nil, ErrUserDontHavePermissionToView
	}
	project, err := s.repo.GetProjectByID(ctx, progress.ProjectID.Hex())
	if err != nil {
		return nil, err
	}

	completedItems := progress.CompletedItems
	if completedItems == nil {
		completedItems = make(map[string]bool)
	}
	maps.Copy(completedItems, *payload.CompletedItems)
	completedTasks := 0
	for _, v := range completedItems {
		if v {
			completedTasks++
		}
	}
	phaseTasks := 0
	for _, phase := range project.Phases {
		phaseTasks += len(phase.Concepts) + len(phase.Tools) + len(phase.Practice)
	}
	capstoneTasks := 0
	for _, capstone := range project.Capstones {
		capstoneTasks += len(capstone.Build) + len(capstone.Concepts) + len(capstone.Tools)
	}
	totalTasks := phaseTasks + capstoneTasks
	percentage := math.Round(float64(completedTasks) / float64(totalTasks) * 100)

	filter := bson.M{"_id": progressID}
	update := bson.M{
		"$set": bson.M{
			"completedItems": completedItems,
			"progress":       percentage,
			"updatedAt":      time.Now(),
		},
	}
	return s.repo.UpdateProgress(ctx, filter, update)
}
