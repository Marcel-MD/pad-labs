import { Timestamp } from './google/protobuf/timestamp.pb';
import { IsEmail, IsUUID, Length, Min } from 'class-validator';

export class ProductDto {
  id: string;
  ownerId: string;
  name: string;
  price: number;
  stock: number;
  createdAt: Timestamp | undefined;
  updatedAt: Timestamp | undefined;
}

export class ProductsDto {
  products: ProductDto[];
}

export interface ProductIdDto {
  id: string;
}

export class CreateProductDto {
  @Length(3)
  name: string;

  @Min(1)
  price: number;

  @Min(0)
  stock: number;
}

export class UpdateProductDto {
  @Length(3)
  name: string;

  @Min(1)
  price: number;

  @Min(0)
  stock: number;
}

export class ProductsQueryDto {
  page: number;
  size: number;
  ownerId: string;
}
