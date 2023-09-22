/* eslint-disable */
import { GrpcMethod, GrpcStreamMethod } from "@nestjs/microservices";
import { Observable } from "rxjs";
import { Timestamp } from "./google/protobuf/timestamp.pb";

export const protobufPackage = "user";

export interface User {
  id: string;
  name: string;
  email: string;
  createdAt: Timestamp | undefined;
  updatedAt: Timestamp | undefined;
}

export interface Users {
  users: User[];
}

export interface RegisterUser {
  email: string;
  name: string;
  password: string;
}

export interface LoginUser {
  email: string;
  password: string;
}

export interface UserId {
  id: string;
}

export interface Token {
  token: string;
}

export interface UsersQuery {
  page: number;
  size: number;
}

export const USER_PACKAGE_NAME = "user";

export interface UserServiceClient {
  register(request: RegisterUser): Observable<Token>;

  login(request: LoginUser): Observable<Token>;

  validate(request: Token): Observable<User>;

  getAll(request: UsersQuery): Observable<Users>;

  getById(request: UserId): Observable<User>;
}

export interface UserServiceController {
  register(request: RegisterUser): Promise<Token> | Observable<Token> | Token;

  login(request: LoginUser): Promise<Token> | Observable<Token> | Token;

  validate(request: Token): Promise<User> | Observable<User> | User;

  getAll(request: UsersQuery): Promise<Users> | Observable<Users> | Users;

  getById(request: UserId): Promise<User> | Observable<User> | User;
}

export function UserServiceControllerMethods() {
  return function (constructor: Function) {
    const grpcMethods: string[] = ["register", "login", "validate", "getAll", "getById"];
    for (const method of grpcMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcMethod("UserService", method)(constructor.prototype[method], method, descriptor);
    }
    const grpcStreamMethods: string[] = [];
    for (const method of grpcStreamMethods) {
      const descriptor: any = Reflect.getOwnPropertyDescriptor(constructor.prototype, method);
      GrpcStreamMethod("UserService", method)(constructor.prototype[method], method, descriptor);
    }
  };
}

export const USER_SERVICE_NAME = "UserService";
