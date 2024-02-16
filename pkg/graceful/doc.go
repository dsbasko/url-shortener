// Package graceful provides a simple way to manage goroutines synchronization
// for a graceful shutdown. It uses sync.WaitGroup under the hood.
//
// The package provides a global instance of a Graceful object, which encapsulates
// a sync.WaitGroup. This instance is used to wait for all added goroutines to finish.
//
// The Add function should be called before starting a goroutine, which increments
// the WaitGroup counter by one.
//
// The Done function should be called when a goroutine finishes, which decrements
// the WaitGroup counter by one.
//
// The Wait function should be called where you want to wait for all goroutines to finish.
// It blocks until the WaitGroup counter is zero.
//
// Example:
//
//	// Some function that starts a goroutine
//	func goFn1(ctx context.Context) {
//		defer graceful.Done()
//		<-ctx.Done()
//		// do some cleanup
//	}
//
//	// Another function with timeout for cleanup function
//	// graceful.Done not called here, because it's called in CleanFn
//	func goFn2(ctx context.Context) {
//		<-ctx.Done()
//		graceful.CleanFn(graceful.DefaultCleanDuration, func() {
//			// do some cleanup
//		})
//	}
//
//	func main() {
//		// Create a context that listens for the SIGINT signal
//		ctx, cancel := graceful.Context(context.Background(), syscall.SIGTERM, syscall.SIGINT)
//		defer cancel()
//
//		// Add 5 goroutines to the Graceful object
//		for i := 0; i < 5; i++ {
//			graceful.Add()
//			go goFn1(ctx)
//			go goFn2(ctx)
//		}
//
//		// Wait for all goroutines to finish
//		graceful.Wait()
//		fmt.Println("Program finished.")
//	}
package graceful
