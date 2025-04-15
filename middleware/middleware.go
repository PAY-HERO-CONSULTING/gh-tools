package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	// "net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	e_errors "emperror.dev/errors"
	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/ctxutils"
	"github.com/PAY-HERO-CONSULTING/gh-tools/custom_types"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logger"
	"github.com/PAY-HERO-CONSULTING/gh-tools/logtrail"
	"github.com/PAY-HERO-CONSULTING/gh-tools/sessions"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/gin"
)

const (
	requestIDHeaderKey = "x-request-id"
	userAgentHeaderKey = "user-agent"
	maxLatency         = time.Second * 5
)

func DefaultMiddlewares(
	authSession sessions.AuthSession,
	requestValidator sessions.RequestValidator,
	driver logtrail.Driver,
	opts ...logtrail.Option,
) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// Comes firsr
		securer(),
		compressor(),
		corsMiddleware(requestValidator),

		// place compressor, cors and securer
		setRequestId(),
		setupAppContext(authSession),
		logTrailMiddleware(driver, opts...),

		// remains as it
		maintenanceMode(),
		latencyMiddleware(),

		// panic recovery should always be lasr
		panicRecovery(),
	}
}

func corsMiddleware(
	requestValidator sessions.RequestValidator,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token, X-SET-AUTH-ACCOUNT-ID, X-AUTH-ACCOUNT-ID")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token, Authorization, X-Requested-With, X-AUTH-Token, X-SET-AUTH-ACCOUNT-ID, X-AUTH-ACCOUNT-ID")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func compressor() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression)
}

func panicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := ctxutils.RequestId(c.Request.Context())

				c.Writer.WriteHeader(http.StatusInternalServerError)
				log.Printf("Failed to recover from panic: %v", err)
				debug.PrintStack()

				if os.Getenv("ENVIRONMENT") != "production" {
					fmt.Fprintf(
						c.Writer,
						`{"error":"panic: %s","details":"See logs for more information (%s)."}`,
						err,
						requestID,
					)
				} else {
					fmt.Fprintf(
						c.Writer,
						`{"error":"internal server error (%s)"}`,
						ctxutils.RequestId(c.Request.Context()),
					)
				}
			}
		}()

		c.Next()
	}
}

func securer() gin.HandlerFunc {
	return secure.Secure(secure.Options{
		SSLRedirect:          strings.ToLower(os.Getenv("FORCE_SSL")) == "true",
		SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:           315360000,
		STSIncludeSubdomains: true,
		FrameDeny:            true,
		ContentTypeNosniff:   true,
		BrowserXssFilter:     true,
	})
}

func setRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := custom_types.GenerateUUID()
		ctx := ctxutils.WithRequestId(c.Request.Context(), requestId)
		c.Request = c.Request.WithContext(ctx)
		c.Header(requestIDHeaderKey, requestId)
		c.Next()
	}
}

func setupAppContext(
	authSession sessions.AuthSession,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authSession == nil {
			c.Next()
		}

		ctx := c.Request.Context()

		userAgent := c.Request.Header.Get(userAgentHeaderKey)
		ctx = ctxutils.WithUserAgent(ctx, userAgent)
		// ctx = ctxutils.W

		ipAddress, err := getIP(c.Request)
		if err != nil {
			logger.Warnf("Unable to parse ipAddress %v from remote address: %v", c.Request.RemoteAddr, err)
		} else {
			ctx = ctxutils.WithIpAddress(ctx, ipAddress)
		}

		tokenInfo, err := authSession.TokenInfo(c)
		if err != nil {
			appError := apperrs.New(
				err,
				apperrs.Internal,
			)

			logger.Infof("unable to parse token from request: token not provided err: [%+v]", appError)

		} else {
			ctx = ctxutils.WithTokenInfo(ctx, tokenInfo)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func maintenanceMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("MAINTENANCE_MODE_ENABLED") == "true" {
			appError := apperrs.New(
				errors.New("maintenance mode"),
				apperrs.ServiceUnavailable,
			)

			c.JSON(appError.Status(), appError.JsonResponse())
			c.Abort()
			return
		}

		c.Next()
	}
}

func latencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if time.Since(start) > maxLatency {
			log.Printf(
				"slow response: %v %v took %v for user: %v",
				c.Request.Method,
				c.Request.URL,
				time.Since(start),
				ctxutils.UserID(c.Request.Context()),
			)
		}
	}
}

func logTrailMiddleware(
	driver logtrail.Driver,
	opts ...logtrail.Option,
) gin.HandlerFunc {
	options := logtrail.NewMiddlewareOptions(
		driver,
	)

	for _, opt := range opts {
		opt.Apply(&options)
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		entry := logtrail.Entry{
			Time:       options.Now(),
			RequestID:  ctxutils.RequestId(c.Request.Context()),
			MacAddress: ctxutils.MacAddress(c.Request.Context()),
			HTTP: logtrail.HTTPEntry{
				ClientIP:  ctxutils.IPAddress(c.Request.Context()),
				UserAgent: c.Request.UserAgent(),
				Method:    c.Request.Method,
				Path:      path,
			},
		}

		sensitiveCall := options.SetupSentivePaths(c)

		// Only override the request body if there is actually one and it doesn't contain sensitive information.
		saveBody := c.Request.Body != nil && !sensitiveCall

		var buf bytes.Buffer

		if saveBody {
			// This should be ok, because the server keeps a reference to the original body,
			// so it can close the original request itself.
			c.Request.Body = io.NopCloser(io.TeeReader(c.Request.Body, &buf))
		}

		c.Next() // process request

		entry.UserID = ctxutils.UserID(c.Request.Context())

		// Consider making this configurable if you need to log unauthorized requests,
		// but keep in mind that in case of a public installation it's a potential DoS attack vector.
		if c.Writer.Status() == http.StatusUnauthorized {
			return
		}

		entry.HTTP.StatusCode = c.Writer.Status()
		entry.HTTP.ResponseSize = c.Writer.Size()
		entry.HTTP.ResponseTime = int(options.SinceInMilliSeconds(entry.Time))

		if saveBody {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil && err != io.EOF {
				logger.Errorf("reading http body err: [%+v]", e_errors.WithStack(err))
			}

			// Convert the body back to a string for further processing
			bodyString := string(bodyBytes)

			// Sanitize the body to remove sensitive information
			sanitizedBody := sanitizeBody(bodyString)

			entry.HTTP.RequestBody = sanitizedBody
		}

		if c.IsAborted() {
			for _, e := range c.Errors {
				_e, _ := e.MarshalJSON()

				entry.HTTP.Errors = append(entry.HTTP.Errors, string(_e))
			}
		}

		err := driver.Store(entry)
		if err != nil {
			logger.Errorf("reading http body err: [%+v]", e_errors.WithStack(err))
			return
		}
	}
}

// sanitizeBody removes sensitive information from the request body.
func sanitizeBody(body string) string {
	sanitized := strings.Replace(body, "password", "********", -1)
	sanitized = strings.Replace(sanitized, "password_confirmation", "********", -1)

	return sanitized
}

func getIP(r *http.Request) (string, error) {
	ips := r.Header.Get("X-FORWARDED-FOR")
	realIp := r.Header.Get("X-REAL-IP")
	forwarded := strings.Split(realIp, ",")

	logger.Infof("X-FORWARDED-FOR: %v", ips)
	logger.Infof("X-REAL-IP: %v", realIp)
	logger.Infof("Forwarded: %v", forwarded)

	if realIp != "" && len(forwarded) > 0 {
		return forwarded[0], nil
	}

	return realIp, nil
}
