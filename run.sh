curl -vXPOST http://localhost:3000/api/run/unique-run-id \
	-H "Content-Type: multipart/mixed" \
	-F "results=@/Users/rajatjindal/Downloads/testresults5.json" \
	-F "metadata={\"runId\":\"unique-run-id\",\"commitSha\": \"shaa\", \"jobName\": \"jobName\", \"repo\":\"dagger/ci-tests\",\"branch\":\"main\",\"format\":\"gojson\", \"link\":\"http://link/to/github/actions/or/elsewhere\", \"tags\":{\"first\":\"tag1\",\"second\": \"tag2\"}};type=application/json"
	