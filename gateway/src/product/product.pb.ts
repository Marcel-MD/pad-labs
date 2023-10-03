/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Empty } from "./google/protobuf/empty.pb";
import { Timestamp } from "./google/protobuf/timestamp.pb";

export const protobufPackage = "product";

export interface Product {
  id: string;
  ownerId: string;
  name: string;
  price: number;
  stock: number;
  createdAt: Timestamp | undefined;
  updatedAt: Timestamp | undefined;
}

export interface Products {
  products: Product[];
}

export interface CreateProduct {
  ownerId: string;
  name: string;
  price: number;
  stock: number;
}

export interface UpdateProduct {
  id: string;
  ownerId: string;
  name: string;
  price: number;
  stock: number;
}

export interface ProductId {
  id: string;
}

export interface ProductsQuery {
  page: number;
  size: number;
  ownerId: string;
}

export const PRODUCT_PACKAGE_NAME = "product";

export interface ProductServiceClient {
  getAll(request: ProductsQuery): Observable<Products>;

  getById(request: ProductId): Observable<Product>;

  create(request: CreateProduct): Observable<ProductId>;

  update(request: UpdateProduct): Observable<Empty>;
}

export interface ProductServiceController {
  getAll(request: ProductsQuery): Promise<Products> | Observable<Products> | Products;

  getById(request: ProductId): Promise<Product> | Observable<Product> | Product;

  create(request: CreateProduct): Promise<ProductId> | Observable<ProductId> | ProductId;

  update(request: UpdateProduct): void;
}

export function ProductServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["getAll", "getById", "create", "update"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("ProductService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("ProductService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const PRODUCT_SERVICE_NAME = "ProductService";
