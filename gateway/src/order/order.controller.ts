import {
  Controller,
  Get,
  Inject,
  OnModuleInit,
  Param,
  UseGuards,
  Post,
  Body,
  Req,
  Query,
  Put,
} from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable, timeout } from 'rxjs';
import { OrderServiceClient, ORDER_SERVICE_NAME } from './order.pb';
import { UserGuard, UserRequest } from '../user/user.guard';
import {
  CreateOrderDto,
  OrderDto,
  OrderIdDto,
  OrdersDto,
  OrdersQueryDto,
  UpdateOrderDto,
} from './order.dto';
import { Empty } from './google/protobuf/empty.pb';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';

const TIMEOUT = 5000;

@ApiTags('Order')
@Controller('order')
export class OrderController implements OnModuleInit {
  private svc: OrderServiceClient;

  @Inject(ORDER_SERVICE_NAME)
  private readonly client: ClientGrpc;

  public onModuleInit(): void {
    this.svc = this.client.getService<OrderServiceClient>(ORDER_SERVICE_NAME);
  }

  @ApiBearerAuth()
  @Post(':productId')
  @UseGuards(UserGuard)
  private async createOrder(
    @Req() req: UserRequest,
    @Param('productId') productId: string,
    @Body() body: CreateOrderDto,
  ): Promise<Observable<OrderIdDto>> {
    var userId = req.user;
    return this.svc
      .create({ userId, productId, ...body })
      .pipe(timeout(TIMEOUT));
  }

  @ApiBearerAuth()
  @Put(':id')
  @UseGuards(UserGuard)
  private async updateOrder(
    @Req() req: UserRequest,
    @Param('id') id: string,
    @Body() body: UpdateOrderDto,
  ): Promise<Observable<Empty>> {
    var productOwnerId = req.user;
    return this.svc
      .update({ id, productOwnerId, ...body })
      .pipe(timeout(TIMEOUT));
  }

  @Get()
  private async getAll(
    @Query() query: OrdersQueryDto,
  ): Promise<Observable<OrdersDto>> {
    var pbQuery = {
      page: query.page || 0,
      size: query.size || 0,
      productOwnerId: query.productOwnerId || '',
      userId: query.userId || '',
    };

    return this.svc.getAll(pbQuery).pipe(timeout(TIMEOUT));
  }

  @Get(':id')
  private async getById(
    @Param('id') id: string,
  ): Promise<Observable<OrderDto>> {
    return this.svc.getById({ id }).pipe(timeout(TIMEOUT));
  }
}
