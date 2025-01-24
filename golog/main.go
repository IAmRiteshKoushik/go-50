package main

import (
	// "errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) LogValue() slog.Value {
	// return slog.IntValue(u.ID)
	return slog.StringValue(fmt.Sprintf(
		"User ID: %d - for user %s", u.ID, u.Username,
	))
}

// We do not have all the features that we would need in the default package
// We would like structured-logging package with log-levels, kv-pairs, better
// flag support
func defaultLogger() {
	// 1.Writing to a file-descriptor
	// 2. Add a prefix
	// 3. Add the message body
	logger := log.New(
		os.Stderr,
		"orders: ",
		log.Ldate|log.Ltime,
	)
	logger.Println("Order submitted")
}

func structuredLogging() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}
	logger := slog.New(NewHandler(opts))
	slog.SetDefault(logger)

	// Passing custom arguments and key-value pairs to the slog.Info method
	slog.Info(
		"Order Submitted",
		slog.String("user", "john"),
		slog.Int("id", rand.Intn(50)),
	)
	slog.Error("Error occured")
}

func main() {
	LOG_FILE := os.Getenv("LOG_FILE")
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	w := io.MultiWriter(os.Stderr, file)
	handlerOpts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}

	logger := slog.New(slog.NewJSONHandler(w, handlerOpts))
	slog.SetDefault(logger)

	// userGroup := slog.Group(
	// 	"user",
	// 	slog.Int("id", rand.Intn(50)),
	// 	slog.String("username", "john"),
	// )

	reqGroup := slog.Group(
		"request",
		slog.String("method", "GET"),
		slog.String("Content-Type", "application/json"),
	)
	requestLogger := logger.With(reqGroup)
	requestLogger.Info("Order submitted", reqGroup)
	requestLogger.Error("Error occured")
	user := &User{rand.Intn(50), "johndoe", "secret"}
	requestLogger.Info("Order submitted", "user", user)

	// opts := &slog.HandlerOptions{
	// 	Level:     slog.LevelDebug,
	// 	AddSource: true,
	// }
	// logger := slog.New(NewHandler(opts))
	//
	// logger.Debug("This is a debug message")
	// logger.Info("This is a info message")
	// logger.Warn("This is a warn message")
	// logger.Error("Oops, something went wrong!", "error", errors.New("Bad Luck my friend"))
}
