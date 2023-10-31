
migration_dir := $(shell echo /pkg/migrations)

migration_path_all_dir := $(shell pwd)$(migration_dir)

dev:
	air -c .air.dev.toml

build:
	go build -o out/ .

up.migrate:
	migrate -database 'postgres://banana:123456@0.0.0.0:4444/ecommerce?sslmode=disable' -source file://$(migration_path_all_dir) -verbose up
down.migrate:
	migrate -database 'postgres://banana:123456@0.0.0.0:4444/ecommerce?sslmode=disable' -source file://$(migration_path_all_dir) -verbose down

# how to use 
# make gen.migrate name=init
gen.migrate:
	migrate create -ext sql -dir pkg/migrations -seq $(name)