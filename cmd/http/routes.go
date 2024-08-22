package http

import (
	"fmt"
	"net/http"

	"github.com/NidzamuddinMuzakki/movies-abishar/cmd/middleware"
	"github.com/NidzamuddinMuzakki/movies-abishar/common/util"
	"github.com/NidzamuddinMuzakki/movies-abishar/config"
	commonCache "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/cache"
	commonMiddleware "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/middleware/gin"
	common "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/registry"
	commonResponse "github.com/NidzamuddinMuzakki/movies-abishar/go-lib-common/response"
	delivery "github.com/NidzamuddinMuzakki/movies-abishar/handler"
	"github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	Register() *gin.Engine
}

type router struct {
	engine   *gin.Engine
	common   common.IRegistry
	delivery delivery.IRegistry
}

func NewRouter(
	common common.IRegistry,
	delivery delivery.IRegistry,
) Router {
	return &router{
		engine:   gin.Default(),
		common:   common,
		delivery: delivery,
	}
}

// @title          movies-abishar Swagger API
// @version        1.0
// @description    movies-abishar Swagger API
// @termsOfService http://swagger.io/terms/

// @contact.name  API Support
// @contact.url   http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func (r *router) Register() *gin.Engine {
	// Middleware
	r.engine.Use(
		commonMiddleware.CORS(),
		commonMiddleware.RequestID(),
		r.common.GetPanicRecoveryMiddleware().PanicRecoveryMiddleware(),
	)

	// handle no-route error (404 not found)
	commonResponse.RouteNotFound(r.engine)

	// Landing
	r.engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, http.StatusText(http.StatusOK))
	})

	// Health Check
	r.engine.GET("/health", r.delivery.GetHealth().Check)
	r.engine.GET("/files/:name", func(c *gin.Context) {
		name, err := c.Params.Get("name")
		fmt.Println(err, "get data file")
		names := fmt.Sprintf("%v", name)
		c.File("files/" + names)

	})
	// v1
	// Configuration
	// r.swagger()
	r.v1()

	return r.engine
}

// func (r *router) swagger() {
// 	docs.SwaggerInfo.Schemes = []string{"http", "https"}
// 	// Route: /docs/index.html
// 	r.engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// }

func (r *router) v1() {

	v1 := r.engine.Group("/v1")
	caches, err := commonCache.NewCache(
		commonCache.WithDriver(commonCache.RedisDriver),
		commonCache.WithHost(config.Cold.RedisHost),
		commonCache.WithDatabase("0"),
		commonCache.WithUsername(config.Cold.RedisUsername),
		commonCache.WithPassword(config.Cold.RedisPassword),
	)
	util.PanicIfError(err)
	common := common.NewRegistry(common.WithCache(caches))
	middlewareImpl := middleware.NewMiddleware(common)

	v1.POST("/users/register", r.delivery.GetUsers().CreateUsers)
	v1.POST("/users/login", r.delivery.GetUsers().LoginUsers)
	v1.POST("/users/logout", middlewareImpl.AuthJWT(), r.delivery.GetUsers().LogoutUsers)

	v1.GET("/movies/list", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetListMovies)
	v1.POST("/movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().CreateMovies)
	v1.POST("/movies/vote/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().VoteMovies)
	v1.POST("/movies/unvote/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().UnVoteMovies)

	v1.PUT("/movies/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().UpdateMovies)
	v1.GET("/movies/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetDetailMovies)
	v1.GET("/movies/users/:id/:duration", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetDetailMoviesUsers)

	v1.GET("/movies/vote-users", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetVoteMovies)

	v1.GET("/movies/most-viewed-movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostViewedMovies)
	v1.GET("/movies/most-vote-movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostVoteMovies)
	v1.GET("/movies/most-viewed-genre", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostViewedGenre)

}
