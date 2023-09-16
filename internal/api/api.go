package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/diezfx/split-app-backend/internal/config"
	"github.com/diezfx/split-app-backend/internal/service"
	"github.com/diezfx/split-app-backend/pkg/auth"
	"github.com/diezfx/split-app-backend/pkg/logger"
	"github.com/diezfx/split-app-backend/pkg/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIHandler struct {
	projectService ProjectService
}

func newAPIHandler(projectService ProjectService) *APIHandler {
	return &APIHandler{projectService: projectService}
}

func InitAPI(cfg *config.Config, projectService ProjectService) *http.Server {
	mr := gin.New()
	mr.Use(gin.Recovery())
	mr.Use(middleware.HTTPLoggingMiddleware())
	mr.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "OPTION"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
	r := mr.Group("/api/v1.0/")
	if !cfg.IsLocal() {
		r.Use(auth.AuthMiddleware(cfg.Auth))
	}
	apiHandler := newAPIHandler(projectService)
	r.GET("projects/:id", apiHandler.getProjectByIDHandler)
	r.GET("projects", apiHandler.getProjectsHandler)
	r.POST("projects", apiHandler.addProjectHandler)
	r.GET("users/:id/costs", apiHandler.getUserCostsHandler)
	r.POST("projects/:id/transactions", apiHandler.addTransactionHandler)
	r.GET("projects/:id/users", apiHandler.getProjectUsersHandler)
	r.POST("projects/:id/users", apiHandler.addProjectUserHandler)
	r.GET("projects/:id/costs", apiHandler.getProjectCostsHandler)

	return &http.Server{
		Handler: mr,
		Addr:    "localhost:5002",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func (api *APIHandler) getProjectUsersHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("invalid id givens: %w", errInvalidInput))
	}
	users, err := api.projectService.GetProjectUsers(ctx, id)
	if err != nil {
		handleError(ctx, fmt.Errorf("getUsers: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (api *APIHandler) getProjectCostsHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("invalid id givens: %w", errInvalidInput))
	}
	costs, err := api.projectService.GetCostsByProject(ctx, id)
	if err != nil {
		handleError(ctx, fmt.Errorf("getUsers: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, ProjectCostsFromService(costs))
}

func (api *APIHandler) getUserCostsHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		handleError(ctx, fmt.Errorf("invalid id given: %w", errInvalidInput))
		return
	}
	users, err := api.projectService.GetCostsByUser(ctx, id)
	if err != nil {
		handleError(ctx, fmt.Errorf("getUsers: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, UserCostsFromService(users))
}

func (api *APIHandler) addProjectUserHandler(ctx *gin.Context) {
	projectIDStr := ctx.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("invalid id given: %w: %w", err, errInvalidInput))
		return
	}
	projectUser := User{}
	err = ctx.BindJSON(&projectUser)
	if err != nil || projectUser.ID == "" {
		handleError(ctx, fmt.Errorf("invalid body: %w: %w", err, errInvalidInput))
		return
	}
	err = api.projectService.AddProjectUser(ctx, projectID, projectUser.ID)
	if err != nil {
		handleError(ctx, fmt.Errorf("getUsers: %w", err))
		return
	}

	ctx.Status(http.StatusCreated)
}

func (api *APIHandler) getProjectsHandler(ctx *gin.Context) {
	var queryParams GetProjectsQueryParams
	err := ctx.BindQuery(&queryParams)
	if err != nil {
		handleError(ctx, fmt.Errorf("getProjectHandler parse query params : %w: %w", errInvalidInput, err))
		return
	}

	proj, err := api.projectService.GetProjects(ctx)
	if err != nil {
		handleError(ctx, err)
		return
	}

	projectList := make([]Project, 0, len(proj))
	for _, p := range proj {
		projectList = append(projectList, ProjectFromServiceProject(p))
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, projectList)
}

func (api *APIHandler) getProjectByIDHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse id: %w: %w", errInvalidInput, err))
		return
	}

	proj, err := api.projectService.GetProjectByID(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.JSON(http.StatusOK, ProjectFromServiceProject(proj))
}

func (api *APIHandler) addTransactionHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(ctx, fmt.Errorf("invalid projectId: %w: %w", errInvalidInput, err))
		return
	}

	var transaction AddTransaction
	if err = ctx.BindJSON(&transaction); err != nil {
		handleError(ctx, fmt.Errorf("parse add transaction body: %w: %w", errInvalidInput, err))
		return
	}

	svcTransaction, err := transaction.Validate()
	if err != nil {
		handleError(ctx, fmt.Errorf("validate transaction: %w: %w", errInvalidInput, err))
		return
	}

	err = api.projectService.AddTransaction(ctx, id, svcTransaction)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (api *APIHandler) addProjectHandler(ctx *gin.Context) {
	var body AddProject
	err := ctx.BindJSON(&body)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse add project body: %w: %w", errInvalidInput, err))
		return
	}

	idParsed, err := uuid.Parse(body.ID)
	if err != nil {
		handleError(ctx, fmt.Errorf("parse id in body: %w: %w", errInvalidInput, err))
		return
	}

	project := service.Project{ID: idParsed, Name: body.Name, Members: body.Members}

	proj, err := api.projectService.AddProject(ctx, project)
	if err != nil {
		handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, proj)
}

func handleError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, errInvalidInput):
		logger.Info(ctx).Err(err).Msg("request failed with invalid input")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Reason:    "invalid input",
		})
	case errors.Is(err, service.ErrProjectNotFound):
		logger.Info(ctx).Err(err).Msg("not found")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusNotFound,
			Reason:    "not found",
		})
	default:
		logger.Error(ctx, err).Msg("unexpected error occurred")
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Reason:    "unexpected",
		})
	}
}
