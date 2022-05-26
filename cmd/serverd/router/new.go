package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swagger_files "github.com/swaggo/files"
	gin_swagger "github.com/swaggo/gin-swagger"
	"logur.dev/logur"
	"os"
	session_handlers "scaffold-api-server/internal/handlers/session-handlers"
	system_handlers "scaffold-api-server/internal/handlers/system-handlers"
	user_handlers "scaffold-api-server/internal/handlers/user-handlers"
	"scaffold-api-server/internal/repository"
	"scaffold-api-server/internal/services/cache"
	"scaffold-api-server/internal/services/database"
	"scaffold-api-server/pkg/auth"
	build_info "scaffold-api-server/pkg/build-info"
	"scaffold-api-server/pkg/web"
)

// New hooks all routes to handlers
func New(
	buildInfo build_info.BuildInfo,
	logger logur.LoggerFacade,
	db *database.DB,
	cacheClient *cache.Client,
	isDevEnv bool,
	address string,
) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/api/v1/liveness"),
		gin.Recovery(),
	)
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	// using cors on dev mode
	if isDevEnv {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = []string{os.Getenv("ORIGIN")}
	}
	r.Use(cors.New(config))
	r.Use(auth.Middleware(logger))
	r.Use(web.RequestIDMiddleware())

	// setup repo
	repo := repository.New(logger, db, cacheClient)

	// setup handler
	systemHandler := system_handlers.New(buildInfo)
	sessionHandler := session_handlers.New(logger, repo)
	userHandler := user_handlers.New(logger, repo)

	rootURL := r.Group("api/v1")
	rootURL.GET("/liveness", systemHandler.Liveness)

	// For profiling
	pprof.Register(r, "api/v1/debug/pprof")

	// Auth endpoint
	sessionURL := r.Group("api/v1/sessions")
	sessionURL.POST("/login", sessionHandler.Login)
	sessionURL.POST("/logout", sessionHandler.Logout)

	// Business endpoint
	userURL := r.Group("api/v1/investors")
	userURL.POST("/signup", userHandler.Signup)

	// Swagger API Docs for QA/Dev
	if isDevEnv {
		_ = gin_swagger.URL(address + " /swagger/doc.json")
		r.GET("/swagger/*any", gin_swagger.WrapHandler(swagger_files.Handler))
	}

	return r
}
