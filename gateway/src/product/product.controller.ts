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
  UseInterceptors,
} from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable, TimeoutError, throwError, timeout } from 'rxjs';
import { ProductServiceClient, PRODUCT_SERVICE_NAME } from './product.pb';
import { UserGuard, UserRequest } from '../user/user.guard';
import {
  CreateProductDto,
  ProductDto,
  ProductIdDto,
  ProductsDto,
  ProductsQueryDto,
  UpdateProductDto,
} from './product.dto';
import { Empty } from './google/protobuf/empty.pb';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { CacheInterceptor } from '@nestjs/cache-manager';

var nrOfTimeouts = 0;

const TIMEOUT = {
  each: 2000,
  with: () =>
    throwError(() => {
      nrOfTimeouts++;
      if (nrOfTimeouts > 2) {
        console.error('Product services timeout 3 times');
      }
      return new TimeoutError();
    }),
};

@ApiTags('Product')
@Controller('product')
@UseInterceptors(CacheInterceptor)
export class ProductController implements OnModuleInit {
  private svc: ProductServiceClient;

  @Inject(PRODUCT_SERVICE_NAME)
  private readonly client: ClientGrpc;

  public onModuleInit(): void {
    this.svc =
      this.client.getService<ProductServiceClient>(PRODUCT_SERVICE_NAME);
  }

  @ApiBearerAuth()
  @Post()
  @UseGuards(UserGuard)
  private async createProduct(
    @Req() req: UserRequest,
    @Body() body: CreateProductDto,
  ): Promise<Observable<ProductIdDto>> {
    var ownerId = req.user;
    return this.svc.create({ ownerId, ...body }).pipe(timeout(TIMEOUT));
  }

  @ApiBearerAuth()
  @Put(':id')
  @UseGuards(UserGuard)
  private async updateProduct(
    @Req() req: UserRequest,
    @Param('id') id: string,
    @Body() body: UpdateProductDto,
  ): Promise<Observable<Empty>> {
    var ownerId = req.user;
    return this.svc.update({ id, ownerId, ...body }).pipe(timeout(TIMEOUT));
  }

  @Get()
  private async getAll(
    @Query() query: ProductsQueryDto,
  ): Promise<Observable<ProductsDto>> {
    var pbQuery = {
      page: query.page || 0,
      size: query.size || 0,
      ownerId: query.ownerId || '',
    };

    return this.svc.getAll(pbQuery).pipe(timeout(TIMEOUT));
  }

  @Get(':id')
  private async getById(
    @Param('id') id: string,
  ): Promise<Observable<ProductDto>> {
    return this.svc.getById({ id }).pipe(timeout(TIMEOUT));
  }
}
