package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	var (
		server          string
		dur             time.Duration
		subject         string
		includeInbox    bool
		includeInternal bool
		dedupe          bool
	)

	flag.StringVar(&server, "server", "localhost:4222", "NATS server URLs.")
	flag.DurationVar(&dur, "duration", 0, "Duration to profile. Set to zero, to wait until interrupt.")
	flag.StringVar(&subject, "subject", ">", "Subject pattern to profile.")
	flag.BoolVar(&includeInbox, "include-inbox", false, "Include _INBOX.> subjects.")
	flag.BoolVar(&includeInternal, "include-internal", false, "Include internal subjects starting with $.")
	flag.BoolVar(&dedupe, "dedupe", true, "Only print each unique subject once.")

	flag.Parse()

	nc, err := nats.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()
	defer nc.Close()

	// Define context to run for a set period of time.
	ctx := context.Background()
	if dur != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dur)
		defer cancel()
	}

	// Map to keep track of subjects for deduping (if enabled).
	subs := make(map[string]struct{})

	// Subscribe to subject.
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		s := msg.Subject

		// Function to determine if the subject is relevant to track.
		if strings.HasPrefix(s, "_INBOX.") && !includeInbox {
			return
		}

		if strings.HasPrefix(s, "$") && !includeInternal {
			return
		}

		if dedupe {
			if _, ok := subs[s]; ok {
				return
			}
			subs[s] = struct{}{}
		}

		fmt.Println(s)
	})
	if err != nil {
		log.Fatal(err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	// Wait for timeout or interrupt.
	select {
	case <-ctx.Done():
	case <-sigch:
	}
}
