package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/a5932016/go-ddd-example/config"
	"github.com/a5932016/go-ddd-example/util/log"
	"github.com/a5932016/go-ddd-example/util/mGin"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var Validate *validator.Validate = validator.New()

func init() {
	// Register custom validation for name characters
	Validate.RegisterValidation("name-chars", validateNameChars)
}

func validateNameChars(fl validator.FieldLevel) bool {
	// Define dangerous character regex (replace with yours)
	var dangerousCharsRegex = regexp.MustCompile("[^a-zA-Z0-9_() -]")
	return !dangerousCharsRegex.MatchString(fl.Field().String())
}

// RunServer provide run http or https protocol.
func (rH Handler) RunServer() (err error) {
	if err = rH.perRepo.LoadPolicy(); err != nil {
		return
	}

	var (
		httpSrv = &http.Server{
			Addr:           ":" + config.Env.Core.Port,
			Handler:        rH.routerEngine(),
			ReadTimeout:    2 * time.Second, // Reduced for quicker detection
			WriteTimeout:   5 * time.Second, // Reduced for faster response closure
			MaxHeaderBytes: 1 << 16,         // 64KB header limit
		}
		errCh = make(chan error)
	)

	// http server
	go func() {
		binding.Validator = new(mGin.DefaultValidator)
		mGin.SetResponseCodePrefix(1)

		log.Info("HTTP server is running on " + config.Env.Core.Port + " port.")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- errors.Wrap(err, "listen and serve http")
		}
	}()

	shutdown := func(httpSrv *http.Server) {
		log.Warning("Gracefully Shutdown Server ...")

		var (
			finish      []interface{}
			finishCount int
			finishCh    = make(chan interface{})
		)

		timeout, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer timeoutCancel()

		// Shutdown httpSrv
		finishCount++
		go func() {
			var err error
			defer func() {
				if err != nil {
					finishCh <- err
				} else {
					finishCh <- struct{}{}
				}
			}()
			if err = httpSrv.Shutdown(timeout); err != nil {
				err = errors.Wrap(err, "Shutdown httpSrv")
			}
		}()

		for {
			select {
			case f := <-finishCh:
				if err, ok := f.(error); ok {
					log.Error(err)
				}
				if finish = append(finish, f); len(finish) == finishCount {
					return
				}
			case <-timeout.Done():
				log.Error("Gracefully Shutdown Timeout.")
				return
			}
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-errCh:
		return err
	case <-quit:
		shutdown(httpSrv)
		return nil
	}
}

func (rH Handler) routerEngine() *gin.Engine {
	// set server mode
	gin.SetMode(config.Env.Core.Mode)

	r := gin.New()
	r.RedirectTrailingSlash = false
	middleware := []gin.HandlerFunc{
		gin.Recovery(),
		CORSMiddleware(),
		RateLimitMiddleware(rH),
		mGin.RequestBodyToContextMiddleware(),
	}

	r.Use(middleware...)
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"text": "Welcome to API server.",
		})
	})

	// app
	routers := rH.getRouter()
	for i := range routers {
		r.Handle(
			routers[i].method,
			routers[i].endpoint,
			rH.permissionMiddleware(routers[i].allowancePair).GinFunc(),
			routers[i].worker,
		)
	}

	return r
}
