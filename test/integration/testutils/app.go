package testutils

import (
	"net/http"
	"testing"

	"ticketer/internal/app"
	"ticketer/internal/auth"
	authpostgres "ticketer/internal/auth/postgres"
	"ticketer/internal/availability"
	"ticketer/internal/booking"
	bookingpostgres "ticketer/internal/booking/postgres"
	"ticketer/internal/catalog"
	catalogpostgres "ticketer/internal/catalog/postgres"
	"ticketer/internal/core/lock"
	"ticketer/internal/pricing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

// TestApp encapsulates the test HTTP server and key services.
type TestApp struct {
	Mux         *http.ServeMux
	AuthService auth.Service
}

// SetupTestApp uses Fx to wire dependencies but without starting an actual net.Listener.
func SetupTestApp(t *testing.T, pool *pgxpool.Pool) *TestApp {
	var mux *http.ServeMux
	var authSvc auth.Service

	// Disable fx logging for tests
	app := fxtest.New(t,
		fx.NopLogger,
		fx.Supply(pool), // Supply the testcontainer DB
		fx.Provide(
			zap.NewNop, // No-op logger
			fx.Annotate(lock.NewInMemoryLockService, fx.As(new(lock.LockService))),
		),
		// Repositories
		fx.Provide(
			fx.Annotate(catalogpostgres.NewMovieRepository, fx.As(new(catalog.MovieRepository))),
			fx.Annotate(catalogpostgres.NewShowRepository, fx.As(new(catalog.ShowRepository))),
			fx.Annotate(catalogpostgres.NewShowSeatRepository, fx.As(new(catalog.ShowSeatRepository))),
			fx.Annotate(catalogpostgres.NewTheaterRepository, fx.As(new(catalog.TheaterRepository))),
			fx.Annotate(bookingpostgres.NewBookingRepository, fx.As(new(booking.BookingRepository))),
			fx.Annotate(authpostgres.NewUserRepository, fx.As(new(auth.UserRepository))),
		),
		// Services
		fx.Provide(
			fx.Annotate(availability.New, fx.As(new(availability.Service))),
			fx.Annotate(pricing.New, fx.As(new(pricing.Service))),
			fx.Annotate(booking.NewBookingService, fx.As(new(booking.Service))),
			catalog.NewTheaterService,
			catalog.NewMovieService,
			catalog.NewShowService,
			fx.Annotate(auth.NewService, fx.As(new(auth.Service))),
		),
		// Handlers
		fx.Provide(
			fx.Annotate(booking.NewHandler, fx.As(new(booking.RouteRegistrar)), fx.ResultTags(`group:"routes"`)),
			fx.Annotate(catalog.NewHandler, fx.As(new(booking.RouteRegistrar)), fx.ResultTags(`group:"routes"`)),
			fx.Annotate(availability.NewHandler, fx.As(new(booking.RouteRegistrar)), fx.ResultTags(`group:"routes"`)),
			fx.Annotate(auth.NewHandler, fx.As(new(booking.RouteRegistrar)), fx.ResultTags(`group:"routes"`)),
		),
		fx.Provide(app.NewHTTPServer),
		fx.Populate(&mux, &authSvc),
	)

	app.RequireStart()
	t.Cleanup(func() {
		app.RequireStop()
	})

	return &TestApp{
		Mux:         mux,
		AuthService: authSvc,
	}
}
