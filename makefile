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

gen_order_proto:
	cp ./proto/order.proto ./gateway/proto/order.proto
	protoc --proto_path=proto proto/order.proto --go_out=order/api --go-grpc_out=order/api
	protoc --plugin=gateway/node_modules/.bin/protoc-gen-ts_proto -I=./proto --ts_proto_out=gateway/src/order/ proto/order.proto --ts_proto_opt=nestJs=true --ts_proto_opt=fileSuffix=.pb

image:
	docker image tag pad-labs-gateway:latest marcelvlasenco/pad-labs-gateway:latest
	docker image tag pad-labs-product:latest marcelvlasenco/pad-labs-product:latest
	docker image tag pad-labs-order:latest marcelvlasenco/pad-labs-order:latest
	docker image tag pad-labs-user:latest marcelvlasenco/pad-labs-user:latest
	docker image tag pad-labs-etl:latest marcelvlasenco/pad-labs-etl:latest

push:
	docker image push marcelvlasenco/pad-labs-gateway:latest
	docker image push marcelvlasenco/pad-labs-product:latest
	docker image push marcelvlasenco/pad-labs-order:latest
	docker image push marcelvlasenco/pad-labs-user:latest
	docker image push marcelvlasenco/pad-labs-etl:latest