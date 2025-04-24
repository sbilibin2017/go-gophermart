package apps

// import (
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/jmoiron/sqlx"

// 	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
// 	"github.com/sbilibin2017/go-gophermart/internal/repositories"
// 	"github.com/sbilibin2017/go-gophermart/internal/routers"
// 	"github.com/sbilibin2017/go-gophermart/internal/services"
// 	"github.com/sbilibin2017/go-gophermart/internal/validators"
// )

// func ConfigureAccrualServer(db *sqlx.DB, srv *http.Server) {

// 	rsRepo := repositories.NewRewardSaveRepository(e)
// 	rfRepo := repositories.NewRewardFilterILikeRepository(q)
// 	repoRewardExists := repositories.NewRewardExistsRepository(q)

// 	oeRepo := repositories.NewOrderExistsRepository(q)
// 	ofRepo := repositories.NewOrderGetRepository(q)
// 	osRepo := repositories.NewOrderSaveRepository(e)

// 	val := validator.New()
// 	validators.RegisterLuhnValidator(val)
// 	validators.RegisterRewardTypeValidator(val)

// 	rrSvc := services.NewRewardRegisterService(val, repoRewardExists, rsRepo)
// 	oaSvc := services.NewOrderAcceptService(val, oeRepo, osRepo, rfRepo)
// 	ogSvc := services.NewOrderGetService(val, ofRepo)

// 	r := chi.NewRouter()
// 	accrualR := routers.NewAccrualRouter(db, rrSvc, oaSvc, ogSvc)
// 	r.Mount("/api", accrualR)

// 	srv.Handler = r
// }
