### md5-fetcher-test-task

This is a test task for a company that requested to stay anonymous.

The task itself is as follows:

Write a CLI tool that fetches sites and outputs their MD5 hashes to the console. The sites must be fetched in parallel, number of parallel requests should be specified by a CLI flag. The result ordering is irrelevant.

Usage:

 - Compile the application (`go build main.go`)
 - Run the produced binary (i.e. `./main -parallel 3 adjust.com google.com reddit.com`)

CLI Flags: 

 - `-parallel` the number of requests allowed to run in parallel to each other. Negative, zero or unset values will default to 10. Values greater than the number of urls will be rounded down to the number of urls.

