package app

import (
	"log/slog"
	"net/http"
	"time"

	authhandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/handler"
	authrepository "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/repository"
	authservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/service"
	bookinghandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/handler"
	bookingrepository "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/repository"
	bookingservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/booking/service"
	conferencemock "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/conference/mock"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/config"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/clock"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/httpcommon"
	jwtplatform "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/jwt"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/middleware"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/password"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/platform/postgres"
	roomhandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/handler"
	roomrepository "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/repository"
	roomservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/room/service"
	schedulehandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/handler"
	schedulerepository "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/repository"
	scheduleservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/schedule/service"
	slothandler "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/handler"
	slotrepository "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/repository"
	slotservice "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/slot/service"
	"github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type modules struct {
	handler          http.Handler
	slotRefillWorker *worker.SlotRefillWorker
}

func buildModules(cfg config.Config, logger *slog.Logger, db *pgxpool.Pool) modules {
	systemClock := clock.New()
	transactor := postgres.NewTransactor(db)
	jwtManager := jwtplatform.New(cfg.JWT.Secret, time.Duration(cfg.JWT.TTLSeconds)*time.Second)
	passwordManager := password.New(12)

	authRepo := authrepository.New(db)
	roomRepo := roomrepository.New(db)
	scheduleRepo := schedulerepository.New(db)
	slotRepo := slotrepository.New(db)
	bookingRepo := bookingrepository.New(db)

	slotSvc := slotservice.New(roomRepo, scheduleRepo, slotRepo, transactor, systemClock)
	scheduleSvc := scheduleservice.New(
		roomRepo,
		scheduleRepo,
		slotSvc,
		transactor,
		systemClock,
		cfg.Slots.GenerationWindowDays,
	)
	authSvc := authservice.New(authRepo, jwtManager, passwordManager, transactor)
	roomSvc := roomservice.New(roomRepo)
	conferenceService := conferencemock.New(cfg.Conference.BaseURL, cfg.Conference.MockMode)
	bookingSvc := bookingservice.New(bookingRepo, slotRepo, conferenceService, systemClock, transactor)

	authHandler := authhandler.New(authSvc)
	roomHandler := roomhandler.New(roomSvc)
	scheduleHandler := schedulehandler.New(scheduleSvc)
	slotHandler := slothandler.New(slotSvc)
	bookingHandler := bookinghandler.New(bookingSvc)

	mux := http.NewServeMux()
	authenticated := middleware.Authenticate(jwtManager)

	mux.Handle("/_info", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpcommon.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}))

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	mux.Handle("POST /dummyLogin", authHandler.DummyLogin())
	mux.Handle("POST /register", authHandler.Register())
	mux.Handle("POST /login", authHandler.Login())

	mux.Handle("GET /rooms/list", authenticated(roomHandler.List()))
	mux.Handle("POST /rooms/create", authenticated(roomHandler.Create()))
	mux.Handle("POST /rooms/{roomId}/schedule/create", authenticated(scheduleHandler.Create()))
	mux.Handle("GET /rooms/{roomId}/slots/list", authenticated(slotHandler.ListAvailable()))

	mux.Handle("POST /bookings/create", authenticated(bookingHandler.Create()))
	mux.Handle("GET /bookings/list", authenticated(bookingHandler.ListAll()))
	mux.Handle("GET /bookings/my", authenticated(bookingHandler.ListMyFuture()))
	mux.Handle("POST /bookings/{bookingId}/cancel", authenticated(bookingHandler.Cancel()))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		httpcommon.WriteError(w, http.StatusNotFound, httpcommon.CodeNotFound, "route not found")
	}))

	return modules{
		handler: middleware.Recover(logger)(middleware.Logging(logger)(mux)),
		slotRefillWorker: worker.NewSlotRefillWorker(
			logger,
			scheduleRepo,
			slotSvc,
			systemClock,
			cfg.Slots.GenerationWindowDays,
			time.Duration(cfg.Slots.RefillIntervalSeconds)*time.Second,
		),
	}
}
