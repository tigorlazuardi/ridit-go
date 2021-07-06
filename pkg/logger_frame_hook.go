package pkg

import (
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

type FrameHook struct {
	Disabled bool
}

func (f FrameHook) Fire(entry *logrus.Entry) error {
	if !f.Disabled {
		entry.Data["_location"] = GetCallerLocation(6)
	}
	return nil
}

func (FrameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func GetCallerLocation(skip int) string {
	// skips runtime.Callers, getFrame, and wrapper func itself
	targetFrameIndex := skip + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown", File: "unknown", Line: 0}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	line := strconv.Itoa(frame.Line)
	return frame.File + ":" + line
}
