import { Timestamp } from './google/protobuf/timestamp.pb';
import { Length, Min } from 'class-validator';

export class OrderDto {
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

export class OrdersDto {
  orders: OrderDto[];
}

export class OrderIdDto {
  id: string;
}

export class CreateOrderDto {
  @Length(3)
  shippingAddress: string;

  @Min(1)
  quantity: number;
}

export class UpdateOrderDto {
  @Length(3)
  status: string;

  @Min(1)
  cost: number;
}

export class OrdersQueryDto {
  page?: number;
  size?: number;
  userId?: string;
  productOwnerId?: string;
}
