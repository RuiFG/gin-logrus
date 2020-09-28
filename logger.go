package gin_logrus

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	// defaultTransformer is the default log transform function Logger middleware uses.
	defaultTransformer = func(logger *logrus.Logger, params FieldsParams) {
		fields := logrus.Fields{}
		//
		option := params.Option
		if option.Host {
			fields["Host"] = params.Request.Host
		}
		if option.Header {
			fields["Header"] = params.Request.Header
		}
		if option.UserAgent {
			fields["UserAgent"] = params.Request.UserAgent()
		}
		if option.Referer {
			fields["Referer"] = params.Request.Referer()
		}
		fields["Latency"] = params.Latency
		fields["ClientIP"] = params.ClientIP
		fields["Path"] = params.Path
		fields["TimeStamp"] = params.TimeStamp
		entry := logger.WithFields(fields)
		code := params.StatusCode
		switch {
		case code >= http.StatusOK && code < http.StatusBadRequest:
			entry.Debugf("[gin-logrus]%d %s %3s", params.StatusCode, params.Method, http.StatusText(params.StatusCode))
		case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
			entry.Warnf("[gin-logrus]%d %s %3s", params.StatusCode, params.Method, http.StatusText(params.StatusCode))
		default:
			entry.Errorf("[gin-logrus]%d %s %3s", params.StatusCode, params.Method, http.StatusText(params.StatusCode))
		}
	}
)

type LogTransformer func(logger *logrus.Logger, params FieldsParams)

// FieldsParams is the logrus Fields params·
type FieldsParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}

	Option OptionalFieldsParams
}

type OptionalFieldsParams struct {
	//Optional fields
	Host      bool
	Referer   bool
	UserAgent bool
	Header    bool
}

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	Logger    *logrus.Logger
	Formatter LogTransformer
	SkipPaths []string
	Option    OptionalFieldsParams
}

func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Logger: logrus.New(),
		Option: OptionalFieldsParams{},
	})
}
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	logger := config.Logger
	transformer := config.Formatter
	if transformer == nil {
		transformer = defaultTransformer
	}
	var skip map[string]struct{}
	if length := len(config.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range config.SkipPaths {
			skip[path] = struct{}{}
		}
	}
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		// Process request
		ctx.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := FieldsParams{
				Request: ctx.Request,
				Keys:    ctx.Keys,
				Option:  config.Option,
			}
			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = ctx.ClientIP()
			param.Method = ctx.Request.Method
			param.StatusCode = ctx.Writer.Status()
			param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = ctx.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			transformer(logger, param)
		}
	}
}
