package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/diezfx/split-app-backend/internal/storage"
	"github.com/google/uuid"
)

type Service struct {
	projStorage ProjectStorage
}

// AddTransaction implements api.ProjectService.
func (s *Service) AddTransaction(ctx context.Context, projID uuid.UUID, transaction Transaction) error {
	_, err := s.projStorage.GetProjectByID(ctx, projID)
	if errors.Is(err, storage.ErrNotFound) {
		return ErrProjectNotFound
	}
	if err != nil {
		return fmt.Errorf("get project:%w", err)
	}
	err = s.projStorage.AddTransaction(ctx, projID, ToStorageTransaction(transaction))
	if err != nil {
		return fmt.Errorf("add transaction: %w", err)
	}
	return nil
}

func New(projStorage ProjectStorage) *Service {
	return &Service{projStorage: projStorage}
}

func (s *Service) GetProjectByID(ctx context.Context, id uuid.UUID) (Project, error) {
	proj, err := s.projStorage.GetProjectByID(ctx, id)
	if errors.Is(err, storage.ErrNotFound) {
		return Project{}, ErrProjectNotFound
	}
	if err != nil {
		return Project{}, fmt.Errorf("get project:%w", err)
	}
	return FromStorageProject(proj), nil
}

func (s *Service) GetProjects(ctx context.Context) ([]Project, error) {
	projs, err := s.projStorage.GetProjects(ctx)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, ErrProjectNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get project:%w", err)
	}

	var projectList = make([]Project, 0, len(projs))
	for _, p := range projs {
		projectList = append(projectList, FromStorageProject(p))
	}

	return projectList, nil
}

func (s *Service) AddProject(ctx context.Context, project Project) (Project, error) {
	_, err := s.projStorage.GetProjectByID(ctx, project.ID)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return Project{}, fmt.Errorf("add project: %w", err)
	}
	if !errors.Is(err, storage.ErrNotFound) {
		return Project{}, fmt.Errorf("add project: %w", err)
	}

	proj, err := s.projStorage.AddProject(ctx, ToStorageProject(project))
	if err != nil {
		return Project{}, fmt.Errorf("add project: %w", err)
	}

	return FromStorageProject(proj), nil
}
