package apiserver

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

import (
	"ascendex.io/act-aws-lambda-s3/library/ecode"

	"github.com/gorilla/mux"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	cfg    *Config
	log    *zap.Logger
	Router *mux.Router
}

func (a *App) Initialize(cfg *Config) error {

	a.cfg = cfg
	a.log = a.initLogger()
	a.initializeRoutes()

	return nil
}

func (a *App) initLogger() *zap.Logger {
	// 设置一些基本日志格式 具体含义还比较好理解，直接看zap源码也不难懂
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})

	// 实现两个判断日志等级的interface (其实 zapcore.*Level 自身就是 interface)
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	// 获取 info、warn日志文件的io.Writer 抽象 getWriter() 在下方实现
	infoWriter := getWriter(a.cfg.LogInfoPath)
	warnWriter := getWriter(a.cfg.LogErrPath)

	// 最后创建具体的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(warnWriter), warnLevel),
	)

	zlogger := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	return zlogger
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func (a *App) initializeRoutes() {
	r := mux.NewRouter()
	// ping api服务在线
	r.HandleFunc("/ping", a.handlerPing).Methods("GET")

	// 处理 短视频(图片，视频), 影片(图片，视频), 头像的上传回调
	r.HandleFunc("/v1/sns/s3/callback", a.handlerSNSS3Callback).Methods("POST")

	a.Router = r
}

func (a *App) Run() {
	var err error
	var hs *http.Server

	hs = &http.Server{
		Addr: a.cfg.ListenAddr,
		Handler: cors.New(cors.Options{
			AllowedMethods:   []string{http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
			AllowedHeaders:   []string{"*"},
			AllowedOrigins:   []string{"*"}, //a.cfg.CorsDomains, []string{"*"}
			AllowCredentials: true,
		}).Handler(a.Router),
	}

	err = hs.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *App) handlerPing(w http.ResponseWriter, r *http.Request) {
	var body []byte
	body, _ = json.Marshal(&Response{Code: ecode.OK.Code(), Error: "", Data: "pong"})
	_, _ = w.Write(body)
	return
}

func (a *App) WriteErrorResponse(w http.ResponseWriter, err error) {
	var code ecode.Codes
	var body []byte

	code = errors.Cause(err).(ecode.Codes)
	if code.Code() < 0 {
		w.WriteHeader(-code.Code())
	} else {
		body, _ = json.Marshal(&Response{Code: code.Code(), Error: code.Error(), Data: ""})
		_, _ = w.Write(body)
	}
}
