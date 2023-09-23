swag:
	swag init --parseDependency

up:
	docker-compose up --build

gen_user_proto:
	cp ./proto/user.proto ./gateway/proto/user.proto
	protoc --proto_path=proto proto/user.proto --go_out=user/api --go-grpc_out=user/api
	protoc --plugin=gateway/node_modules/.bin/protoc-gen-ts_proto -I=./proto --ts_proto_out=gateway/src/user/ proto/user.proto --ts_proto_opt=nestJs=true --ts_proto_opt=fileSuffix=.pb

gen_product_proto:
	cp ./proto/product.proto ./gateway/proto/product.proto
	protoc --proto_path=proto proto/product.proto --go_out=product/api --go-grpc_out=product/api
	protoc --plugin=gateway/node_modules/.bin/protoc-gen-ts_proto -I=./proto --ts_proto_out=gateway/src/product/ proto/product.proto --ts_proto_opt=nestJs=true --ts_proto_opt=fileSuffix=.pb
