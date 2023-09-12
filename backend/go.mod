module github.com/rajatjindal/test-dashboard/backend

go 1.20

require (
	github.com/fermyon/spin/sdk/go v1.4.2
	github.com/jmoiron/sqlx v1.3.5
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/fermyon/spin/sdk/go => github.com/fermyon/spin/sdk/go v1.3.1-0.20230911190838-f554a45ce1f2
