# movies-abishar

# dbdiagram https://dbdiagram.io/d/66c6d569a346f9518cc0f2ab

# migration goose = GOOSE_DRIVER=mysql GOOSE_DBSTRING="sql12727201:BWkG41eVe1@tcp(sql12.freesqldatabase.com:3306)/sql12727201?parseTime=true" ./goose up

# run = go run main.go
# akun admin -> username : admin , password : 123
# architecture repository pattern
# routes -> middleware(validateToken,panicRecovery) -> handler(bindStruct, Validation) -> service -> repository -> db
# routes
#v1.POST("/users/register", r.delivery.GetUsers().CreateUsers)
#v1.POST("/users/login", r.delivery.GetUsers().LoginUsers)
#v1.POST("/users/logout", middlewareImpl.AuthJWT(), r.delivery.GetUsers().LogoutUsers)

# v1.GET("/movies/list", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetListMovies)
# v1.POST("/movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().CreateMovies)
#	v1.POST("/movies/vote/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().VoteMovies)
#	v1.POST("/movies/unvote/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().UnVoteMovies)

#	v1.PUT("/movies/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().UpdateMovies)
#	v1.GET("/movies/:id", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetDetailMovies)
#	v1.GET("/movies/users/:id/:duration", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetDetailMoviesUsers)

#	v1.GET("/movies/vote-users", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetVoteMovies)

#	v1.GET("/movies/most-viewed-movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostViewedMovies)
#	v1.GET("/movies/most-vote-movies", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostVoteMovies)
#	v1.GET("/movies/most-viewed-genre", middlewareImpl.AuthJWT(), r.delivery.GetMovies().GetMostViewedGenre)
