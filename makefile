swag:
	swag init --parseDependency

up:
	docker-compose up --build

gen_user_proto:
	protoc --proto_path=proto proto/user.proto --go_out=user/api --go-grpc_out=user/api

