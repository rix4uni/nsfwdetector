package main

import (
        "bufio"
        "io/ioutil"
        "net/http"
        "os"
        "strings"
        "sync"
        "time"

        "github.com/sirupsen/logrus"
        "github.com/spf13/pflag"
        "github.com/rix4uni/nsfwdetector/banner"
)

func checkNSFW(content string, keywords []string) []string {
        var matchedKeywords []string
        for _, keyword := range keywords {
                if strings.Contains(content, keyword) {
                        matchedKeywords = append(matchedKeywords, keyword)
                }
        }
        return matchedKeywords
}

func checkURL(url string, keywords []string, timeout int, wg *sync.WaitGroup) {
        defer wg.Done()

        // Create an HTTP client with a timeout
        client := &http.Client{
                Timeout: time.Duration(timeout) * time.Second,
        }

        // Try both https:// and http:// if the protocol is missing
        if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
                // Try https:// first
                logrus.Debugf("Trying https:// for URL: %s", url)
                tryURLs := []string{"https://" + url, "http://" + url}

                for _, tryURL := range tryURLs {
                        resp, err := client.Get(tryURL)
                        if err != nil {
                                logrus.Warnf("Error fetching URL %s: %v", tryURL, err)
                                continue // Try the next protocol
                        }
                        defer resp.Body.Close()

                        // Read the response body
                        body, err := ioutil.ReadAll(resp.Body)
                        if err != nil {
                                logrus.Warnf("Error reading response body for %s: %v", tryURL, err)
                                continue
                        }

                        // Check for any matched keywords
                        matchedKeywords := checkNSFW(string(body), keywords)
                        // Format the matched keywords as a comma-separated string
                        if len(matchedKeywords) > 0 {
                                logrus.Infof("%s [%s]", tryURL, strings.Join(matchedKeywords, ", "))
                        } else {
                                logrus.Infof("%s []", tryURL)
                        }
                        return // Exit the loop after successfully processing one URL
                }
                logrus.Warnf("Both protocols failed for %s, skipping.", url)
                return
        }

        // If URL has a valid protocol, continue fetching
        resp, err := client.Get(url)
        if err != nil {
                logrus.Warnf("Error fetching URL %s: %v", url, err)
                return
        }
        defer resp.Body.Close()

        // Read the response body
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                logrus.Warnf("Error reading response body for %s: %v", url, err)
                return
        }

        // Check for any matched keywords
        matchedKeywords := checkNSFW(string(body), keywords)
        // Format the matched keywords as a comma-separated string
        if len(matchedKeywords) > 0 {
                logrus.Infof("%s [%s]", url, strings.Join(matchedKeywords, ", "))
        } else {
                logrus.Infof("%s []", url)
        }
}

func main() {
	// Initialize the flags
        concurrency := pflag.IntP("concurrency", "c", 50, "Number of concurrent workers")
        timeout := pflag.IntP("timeout", "t", 30, "Timeout for each HTTP request in seconds")
        verbose := pflag.Bool("verbose", false, "Enable verbose logging")
        wordlist := pflag.StringP("wordlist", "w", "keywords.txt", "Path to the file containing keywords to check")
        silent := pflag.Bool("silent", false, "silent mode.")
	versionFlag := pflag.Bool("version", false, "Print the version of the tool and exit.")

        // Set the logging level based on the verbose flag
        logrus.SetOutput(os.Stdout)
        if *verbose {
                logrus.SetLevel(logrus.DebugLevel)
        } else {
                logrus.SetLevel(logrus.InfoLevel)
        }

        // Parse the command-line arguments
        pflag.Parse()

        if *versionFlag {
			banner.PrintBanner()
			banner.PrintVersion()
			return
		}

		if !*silent {
			banner.PrintBanner()
		}

        // Read the keywords from the wordlist file
        file, err := os.Open(*wordlist)
        if err != nil {
                logrus.Fatal("Error opening keyword file:", err)
        }
        defer file.Close()

        var keywords []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                keywords = append(keywords, scanner.Text())
        }
        if err := scanner.Err(); err != nil {
                logrus.Fatal("Error reading keyword file:", err)
        }

        // Read URLs from standard input
        var urls []string
        scanner = bufio.NewScanner(os.Stdin)
        for scanner.Scan() {
                url := scanner.Text()
                if url != "" {
                        urls = append(urls, url)
                }
        }
        if err := scanner.Err(); err != nil {
                logrus.Fatal("Error reading input:", err)
        }

        // Create a wait group to manage concurrent workers
        var wg sync.WaitGroup
        sem := make(chan struct{}, *concurrency) // Semaphore to limit the number of concurrent workers

        // Process the URLs concurrently
        for _, url := range urls {
                wg.Add(1)

                // Acquire a slot in the semaphore
                sem <- struct{}{}

                go func(url string) {
                        defer func() {
                                // Release the slot in the semaphore
                                <-sem
                        }()
                        checkURL(url, keywords, *timeout, &wg)
                }(url)
        }

        // Wait for all workers to finish
        wg.Wait()
}
