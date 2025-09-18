package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// -------------------- LOG LEVELS --------------------
type Level int

const (
	DEBUG Level = iota
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	FATAL
	PANIC
)

var levelNames = map[string]Level{
	"DEBUG":    DEBUG,
	"INFO":     INFO,
	"NOTICE":   NOTICE,
	"WARNING":  WARNING,
	"ERROR":    ERROR,
	"CRITICAL": CRITICAL,
	"FATAL":    FATAL,
	"PANIC":    PANIC,
}

func (l Level) String() string {
	for k, v := range levelNames {
		if v == l {
			return k
		}
	}
	return "UNKNOWN"
}

// -------------------- CONTEXT SUPPORT --------------------
type ctxKey string

const RequestIDKey ctxKey = "requestId"

func WithRequestID(ctx context.Context, requestID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// -------------------- LOG ENTRY --------------------
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Component string `json:"component"`
	PodID     string `json:"podId"`
	Env       string `json:"env"`
	RequestID string `json:"requestId,omitempty"`
	File      string `json:"file"`
	Function  string `json:"function"`
	Line      int    `json:"line"`
	Message   string `json:"message"`
}

// -------------------- LOGGER --------------------
type Logger struct {
	component string
	minLevel  Level
	console   *log.Logger
	file      *log.Logger
}

var (
	Database *Logger
	RestAPI  *Logger
	Orders   *Logger
	Payment  *Logger
	Products *Logger
	Shipping *Logger
	Users    *Logger
	Admin    *Logger

	loggers   = make(map[string]*Logger)
	configMap = make(map[string]Level)
	once      sync.Once

	podID = os.Getenv("POD_ID")
	env   = os.Getenv("ENV")
)

func init() {
	// Automatic initialization with default path
	if err := Init("config/logging-level.json"); err != nil {
		log.Printf("Logging auto-init failed: %v\n", err)
	}
}

// -------------------- INIT --------------------
func Init(configPath string) error {
	var err error
	once.Do(func() {
		// Load config
		configFile, e := os.Open(configPath)
		if e != nil {
			err = fmt.Errorf("failed to open logging config: %w", e)
			return
		}
		defer configFile.Close()

		var raw map[string]string
		if e := json.NewDecoder(configFile).Decode(&raw); e != nil {
			err = fmt.Errorf("failed to parse logging config: %w", e)
			return
		}

		for comp, lvl := range raw {
			lvl = strings.ToUpper(lvl)
			if parsed, ok := levelNames[lvl]; ok {
				configMap[comp] = parsed
			} else {
				configMap[comp] = INFO
			}
		}

		// Initialize static loggers
		Database = setupLogger("database")
		RestAPI = setupLogger("restapi")
		Orders = setupLogger("orders")
		Payment = setupLogger("payment")
		Products = setupLogger("products")
		Shipping = setupLogger("shipping")
		Users = setupLogger("users")
		Admin = setupLogger("admin")
	})
	return err
}

// -------------------- COMPONENT LOGGER --------------------
func Component(name string) *Logger {
	switch strings.ToLower(name) {
	case "database":
		return Database
	case "restapi":
		return RestAPI
	case "orders":
		return Orders
	case "payment":
		return Payment
	case "products":
		return Products
	case "shipping":
		return Shipping
	case "users":
		return Users
	case "admin":
		return Admin
	default:
		return setupLogger(name)
	}
}

// -------------------- SETUP LOGGER --------------------
func setupLogger(component string) *Logger {
	if lg, ok := loggers[component]; ok {
		return lg
	}

	_ = os.MkdirAll("logs", 0755)

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", component+".log"),
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     28,
		Compress:   true,
	}
	fileLogger := log.New(fileWriter, "", 0)
	consoleLogger := log.New(os.Stdout, "", 0)

	minLevel := INFO
	if lvl, ok := configMap[component]; ok {
		minLevel = lvl
	}

	lg := &Logger{
		component: component,
		minLevel:  minLevel,
		console:   consoleLogger,
		file:      fileLogger,
	}
	loggers[component] = lg
	return lg
}

// -------------------- COLORS --------------------
func levelColor(level Level) string {
	switch level {
	case DEBUG:
		return "\033[36m" // Cyan
	case INFO:
		return "\033[32m" // Green
	case NOTICE:
		return "\033[34m" // Blue
	case WARNING:
		return "\033[33m" // Yellow
	case ERROR:
		return "\033[31m" // Red
	case CRITICAL:
		return "\033[35m" // Magenta
	case FATAL:
		return "\033[41m" // Red BG
	case PANIC:
		return "\033[41;37m"
	default:
		return "\033[0m"
	}
}

// -------------------- CORE LOG --------------------
func (l *Logger) log(ctx context.Context, level Level, msg string) {
	if level < l.minLevel {
		return
	}
	pc, file, line, ok := runtime.Caller(2)
	funcName := "-"
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			fullName := fn.Name() // e.g. "github.com/yourorg/Adornme/databases.connectOpenSearch"
			parts := strings.Split(fullName, ".")
			funcName = parts[len(parts)-1] // take only the last part
		}
		file = filepath.Base(file)
	}

	rid := "-"
	if ctx != nil {
		if r, ok := ctx.Value(RequestIDKey).(string); ok && r != "" {
			rid = r
		}
	}

	pod := podID
	if pod == "" {
		pod = "-"
	}
	environment := env
	if environment == "" {
		environment = "-"
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Component: defaultIfEmpty(l.component),
		PodID:     pod,
		Env:       environment,
		RequestID: rid,
		File:      file,
		Function:  funcName,
		Line:      line,
		Message:   defaultIfEmpty(msg),
	}

	// JSON log (structured)
	jsonData, _ := json.Marshal(entry)
	l.file.Println(string(jsonData))

	// Console log with |line|
	color := levelColor(level)
	consoleMsg := fmt.Sprintf(
		"%s|%s|%s|%s|%s|%s|%s|%s|%d|%s",
		entry.Timestamp, entry.Level, entry.Component,
		entry.PodID, entry.Env, entry.RequestID,
		entry.File, entry.Function, entry.Line, entry.Message,
	)
	l.console.Println(color + consoleMsg + "\033[0m")

	if level == FATAL {
		os.Exit(1)
	}
	if level == PANIC {
		panic(msg)
	}
}

func defaultIfEmpty(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// -------------------- PUBLIC METHODS --------------------
func (l *Logger) Debug(ctx context.Context, msg string)    { l.log(ctx, DEBUG, msg) }
func (l *Logger) Info(ctx context.Context, msg string)     { l.log(ctx, INFO, msg) }
func (l *Logger) Notice(ctx context.Context, msg string)   { l.log(ctx, NOTICE, msg) }
func (l *Logger) Warning(ctx context.Context, msg string)  { l.log(ctx, WARNING, msg) }
func (l *Logger) Error(ctx context.Context, msg string)    { l.log(ctx, ERROR, msg) }
func (l *Logger) Critical(ctx context.Context, msg string) { l.log(ctx, CRITICAL, msg) }
func (l *Logger) Fatal(ctx context.Context, msg string)    { l.log(ctx, FATAL, msg) }
func (l *Logger) Panic(ctx context.Context, msg string)    { l.log(ctx, PANIC, msg) }

// -------------------- FORMATTED METHODS --------------------
func (l *Logger) Debugf(ctx context.Context, f string, v ...any) {
	l.log(ctx, DEBUG, fmt.Sprintf(f, v...))
}
func (l *Logger) Infof(ctx context.Context, f string, v ...any) {
	l.log(ctx, INFO, fmt.Sprintf(f, v...))
}
func (l *Logger) Noticef(ctx context.Context, f string, v ...any) {
	l.log(ctx, NOTICE, fmt.Sprintf(f, v...))
}
func (l *Logger) Warningf(ctx context.Context, f string, v ...any) {
	l.log(ctx, WARNING, fmt.Sprintf(f, v...))
}
func (l *Logger) Errorf(ctx context.Context, f string, v ...any) {
	l.log(ctx, ERROR, fmt.Sprintf(f, v...))
}
func (l *Logger) Criticalf(ctx context.Context, f string, v ...any) {
	l.log(ctx, CRITICAL, fmt.Sprintf(f, v...))
}
func (l *Logger) Fatalf(ctx context.Context, f string, v ...any) {
	l.log(ctx, FATAL, fmt.Sprintf(f, v...))
}
func (l *Logger) Panicf(ctx context.Context, f string, v ...any) {
	l.log(ctx, PANIC, fmt.Sprintf(f, v...))
}
