package graceful

import (
	"os"
	"syscall"
	"time"
)

// DefaultSignals is a slice of os.Signal that the application listens to for graceful shutdown.
// The application will start the shutdown process when it receives one of these signals.
// By default, it includes SIGTERM, SIGINT, SIGQUIT, and SIGKILL.
var DefaultSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL}

// DefaultCleanDuration is the duration the application waits for ongoing tasks to finish before it shuts down.
// If this duration is exceeded, the application will forcefully terminate the tasks and shut down.
// By default, it is set to 5 seconds.
var DefaultCleanDuration = 5 * time.Second
