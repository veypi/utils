package log

// 封装自 zero log

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"runtime"
	"strconv"
)

type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

func (l Level) String() string {
	return zerolog.Level(l).String()
}

func SetLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

func ParseLevel(s string) (Level, error) {
	l, e := zerolog.ParseLevel(s)
	return Level(l), e
}

var (
	enableCaller bool = true
)

func DisableCaller() {
	enableCaller = false
	SetLogger(originLoger)
}

var fileHook = lumberjack.Logger{
	Filename:   "",
	MaxSize:    128, // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups: 10,  // 日志文件最多保存多少个备份
	MaxAge:     60,  // 文件最多保存多少天
	LocalTime:  true,
	Compress:   true, // 是否压缩
}

// Logger just for dev env, low performance but human-friendly
var originLoger *zerolog.Logger
var logger zerolog.Logger
var WithDeepCaller zerolog.Logger
var WithNoCaller zerolog.Logger

func init() {
	// 输出trace信息，err以上级别调用Err或者Errs使触发
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return string(PanicTrace())
	}
	SetLogger(ConsoleLogger())
}

func SetLogger(l *zerolog.Logger) {
	originLoger = l
	if enableCaller {
		logger = l.With().Timestamp().CallerWithSkipFrameCount(2).Logger()
		WithDeepCaller = l.With().Timestamp().CallerWithSkipFrameCount(3).Logger()
	} else {
		logger = l.With().Timestamp().Logger()
		WithDeepCaller = l.With().Timestamp().Logger()
	}
	WithNoCaller = l.With().Timestamp().Logger()
}

// FileLogger for product, height performance
func FileLogger(fileName string) *zerolog.Logger {
	// 会创建644权限文件夹，导致其他用户无法读取
	fileHook.Filename = fileName
	l := zerolog.New(&fileHook).With().Caller().Timestamp().Logger()
	return &l
}

func ConsoleLogger() *zerolog.Logger {
	l := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.FormatFieldValue = func(i interface{}) string {
			switch i := i.(type) {
			case string:
				s, e := strconv.Unquote(i)
				if e != nil {
					return i
				}
				return s
			default:
				return fmt.Sprintf("%s", i)
			}
		}
	}))
	return &l
}

func HandlerErrs(errs ...error) {
	for _, e := range errs {
		if e != nil {
			WithDeepCaller.Warn().Msg(e.Error())
		}
	}
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *zerolog.Event {
	return logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *zerolog.Event {
	return logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *zerolog.Event {
	return logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *zerolog.Event {
	return logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *zerolog.Event {
	return logger.Error().Stack()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *zerolog.Event {
	return logger.Fatal().Stack()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *zerolog.Event {
	return logger.Panic()
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *zerolog.Event {
	return logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

func PanicTrace() []byte {
	buf := make([]byte, 10240)
	n := runtime.Stack(buf, false)
	if n > 0 {
		count := 0
		for i, v := range buf[:n] {
			if v == '\n' {
				count++
				if count == 11 {
					return buf[i:n]
				}
			}
		}
	}
	return buf[:n]
}
