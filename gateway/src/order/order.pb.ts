/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Empty } from "./google/protobuf/empty.pb";
import { Timestamp } from "./google/protobuf/timestamp.pb";

export const protobufPackage = "order";

export interface Order {
  id: string;
  productId: string;
  userId: string;
  quantity: number;
  cost: number;
  status: string;
  shippingAddress: string;
  createdAt: Timestamp | undefined;
  updatedAt: Timestamp | undefined;
}

export interface Orders {
  orders: Order[];
}

export interface CreateOrder {
  productId: string;
  userId: string;
  quantity: number;
  shippingAddress: string;
}

export interface UpdateOrder {
  id: string;
  productOwnerId: string;
  status: string;
  cost: number;
}

export interface OrderId {
  id: string;
}

export interface OrdersQuery {
  page: number;
  size: number;
  userId: string;
  productOwnerId: string;
}

export const ORDER_PACKAGE_NAME = "order";

export interface OrderServiceClient {
  getAll(request: OrdersQuery): Observable<Orders>;

  getById(request: OrderId): Observable<Order>;

  create(request: CreateOrder): Observable<OrderId>;

  update(request: UpdateOrder): Observable<Empty>;
}

export interface OrderServiceController {
  getAll(request: OrdersQuery): Promise<Orders> | Observable<Orders> | Orders;

  getById(request: OrderId): Promise<Order> | Observable<Order> | Order;

  create(request: CreateOrder): Promise<OrderId> | Observable<OrderId> | OrderId;

  update(request: UpdateOrder): void;
}

export function OrderServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["getAll", "getById", "create", "update"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("OrderService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("OrderService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const ORDER_SERVICE_NAME = "OrderService";
