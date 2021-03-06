package logrus

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

// An entry is the final or intermediate Logrus logging entry. It containts all
// the fields passed with WithField{,s}. It's finally logged when Debug, Info,
// Warn, Error, Fatal or Panic is called on it. These objects can be reused and
// passed around as much as you wish to avoid field duplication.
type Entry struct {
	Logger *Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Debug, Info, Warn, Error, Fatal or Panic
	Level Level

	// Message passed to Debug, Info, Warn, Error, Fatal or Panic
	Message string
}

var baseTimestamp time.Time

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, give a little extra room
		Data: make(Fields, 5),
	}
}

// Returns a reader for the entry, which is a proxy to the formatter.
func (entry *Entry) Reader() (*bytes.Buffer, error) {
	serialized, err := entry.Logger.Formatter.Format(entry)
	return bytes.NewBuffer(serialized), err
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	reader, err := entry.Reader()
	if err != nil {
		return "", err
	}

	return reader.String(), err
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return entry.WithFields(Fields{key: value})
}

// Add a map of fields to the Entry.
func (entry *Entry) WithFields(fields Fields) *Entry {
	data := Fields{}
	for k, v := range entry.Data {
		data[k] = v
	}
	for k, v := range fields {
		data[k] = v
	}
	return &Entry{Logger: entry.Logger, Data: data}
}

func (entry *Entry) log(level Level, msg string) string {
	entry.Time = time.Now()
	entry.Level = level
	entry.Message = msg

	if err := entry.Logger.Hooks.Fire(level, entry); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fire hook", err)
	}

	reader, err := entry.Reader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to obtain reader, %v", err)
	}

	entry.Logger.mu.Lock()
	defer entry.Logger.mu.Unlock()

	_, err = io.Copy(entry.Logger.Out, reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v", err)
	}

	return reader.String()
}

func (entry *Entry) Debug(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Print(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Info(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Warn(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Error(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprint(args...))
	}
}

func (entry *Entry) Fatal(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprint(args...))
	}
	os.Exit(1)
}

func (entry *Entry) Panic(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		msg := entry.log(PanicLevel, fmt.Sprint(args...))
		panic(msg)
	}
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Debugf(format string, args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, fmt.Sprintf(format, args...))
	}
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, fmt.Sprintf(format, args...))
	}
}

// Entry Println family functions

func (entry *Entry) Debugln(args ...interface{}) {
	if entry.Logger.Level >= DebugLevel {
		entry.log(DebugLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Infoln(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Println(args ...interface{}) {
	if entry.Logger.Level >= InfoLevel {
		entry.log(InfoLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warnln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Warningln(args ...interface{}) {
	if entry.Logger.Level >= WarnLevel {
		entry.log(WarnLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Errorln(args ...interface{}) {
	if entry.Logger.Level >= ErrorLevel {
		entry.log(ErrorLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Fatalln(args ...interface{}) {
	if entry.Logger.Level >= FatalLevel {
		entry.log(FatalLevel, entry.sprintlnn(args...))
	}
}

func (entry *Entry) Panicln(args ...interface{}) {
	if entry.Logger.Level >= PanicLevel {
		entry.log(PanicLevel, entry.sprintlnn(args...))
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// fmt.Sprintln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func (entry *Entry) sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
