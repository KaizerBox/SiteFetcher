package main

import (
	"siteFetcher/components/fetcher"
	"siteFetcher/components/gmailer"
)

type FetcherMailer struct {
	Fetcher fetcher.Fetcher
	Emailer gmailer.Gmailer
}
