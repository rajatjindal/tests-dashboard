curl -vXPOST http://localhost:3000/api/run/db6fcf6a-b155-4cb3-8506-671f98bd9924 \
	-H "Content-Type: multipart/mixed" \
	-F "results=@/Users/rajatjindal/Downloads/testresults4.json" \
	-F "metadata={\"runId\":\"db6fcf6a-b155-4cb3-8506-671f98bd9924\",\"repo\":\"dagger/dagger\",\"branch\":\"gotestsum\",\"commitSha\":\"dfb38abfe2f60bdfc1915d88540ce03e27f440e9\",\"format\":\"gojson\",\"link\":\"https://github.com/dagger/dagger/actions/runs/12171942795\",\"tags\":{\"cli_version\":\"dagger v0.14.0 \",\"dagger_version\":\"v0.15.0-010101000000-dev-e5e4c8b7ddd0\",\"engine_version\":\"v0.14.0\",\"job\":\"test-provision4\",\"runner\":\"GitHub Actions 49\",\"workflow\":\"Engine \\u0026 CLI\"},\"createdAt\":\"\"};type=application/json"
	