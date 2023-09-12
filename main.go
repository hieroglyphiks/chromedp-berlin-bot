package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chromedp/chromedp"
)

func main() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3830.0 Safari/537.36"),
		// any port other than "0" will work since 86.0.4199.0
		chromedp.Flag("remote-debugging-port", "9222"),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.DisableGPU,
	)

	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// create context
	ctx, cancel := chromedp.NewContext(
		ctx,
		// chromedp.WithDebugf(log.Printf),
	)

	defer cancel()
	go waitForSignal(cancel)

	// capture screenshot of an element
	if err := chromedp.Run(ctx,
		chromedp.ActionFunc(bot),
	); err != nil {
		log.Fatal(err)
	}
}

func waitForSignal(cf context.CancelFunc) {
	shutdownSignalChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownSignalChannel, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-shutdownSignalChannel:
		close(shutdownSignalChannel)

		log.Printf("Exiting due to user signal: %s\n", sig.String())
		cf()
	}
}
