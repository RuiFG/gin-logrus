# Gin-Logrus
######
gin-logrus is a gin middleware implemented using logrus
# Installation
```shell script
go get github.com/RuiFG/gin-logrus
```
# Simple Example
```go
func main() {
	engine := gin.New()
	engine.Use(gin_logrus.Logger())
	engine.GET("", func(context *gin.Context) {
		context.Status(200)
	})
	_ = engine.Run("127.0.0.1:8080")
}
```
# Customize
You can also define the implementation of LogTransformer yourself
```go
func main() {
	logger := logrus.New()
	engine := gin.New()
	engine.Use(gin_logrus.LoggerWithConfig(gin_logrus.LoggerConfig{
		Logger: logger, 
		Formatter: func(logger *logrus.Logger, params gin_logrus.FieldsParams) {
		//noting to do
	}}))
	engine.GET("", func(context *gin.Context) {
		context.Status(200)
	})
	_ = engine.Run("127.0.0.1:8080")
}
```

# License
This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.