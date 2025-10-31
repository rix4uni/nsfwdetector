## nsfwdetector

A tool to detect website for NSFW content based on a keyword list like `porn`, `adult`, `xxx`, `hentai`.

## Why?
- I created this tool to get NSFW website quickly with no time after that i can give this output to any `Image Recognition` tool for more accurate results
- If i directly use all of website to `Image Recognition` tool that will take too much time.
- These NSFW website you can use in pihole or any other purpose.

## Installation
```
go install github.com/rix4uni/nsfwdetector@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/nsfwdetector/releases/download/v0.0.2/nsfwdetector-linux-amd64-0.0.2.tgz
tar -xvzf nsfwdetector-linux-amd64-0.0.2.tgz
rm -rf nsfwdetector-linux-amd64-0.0.2.tgz
mv nsfwdetector ~/go/bin/nsfwdetector
```
Or download [binary release](https://github.com/rix4uni/nsfwdetector/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/nsfwdetector.git
cd nsfwdetector; go install
```

## Usage
```yaml
Usage of nsfwdetector:
  -c, --concurrency int   Number of concurrent workers (default 50)
  -o, --output string     Path to the output file to save results
      --silent            silent mode.
  -t, --timeout int       Timeout for each HTTP request in seconds (default 30)
      --verbose           Enable verbose logging
      --version           Print the version of the tool and exit.
  -w, --wordlist string   Path to the file containing keywords to check (default "keywords.txt")
```

## Usage Examples
Single Target:
```yaml
▶ echo "domain.com" | nsfwdetector --silent --wordlist keywords.txt
```

Multiple Targets:
```yaml
▶ cat targets.txt
domain.com
anotherdomain.com

▶ cat targets.txt | nsfwdetector --silent --wordlist keywords.txt
```
