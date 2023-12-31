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
import { Observable, firstValueFrom, timeout } from 'rxjs';
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

const TIMEOUT = 5000;

@ApiTags('Product')
@Controller('product')
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
  private async getAll(@Query() query: ProductsQueryDto): Promise<ProductsDto> {
    var pbQuery = {
      page: query.page || 0,
      size: query.size || 0,
      ownerId: query.ownerId || '',
    };

    return await firstValueFrom(
      this.svc.getAll(pbQuery).pipe(timeout(TIMEOUT)),
    );
  }

  @Get(':id')
  private async getById(@Param('id') id: string): Promise<ProductDto> {
    return await firstValueFrom(
      this.svc.getById({ id }).pipe(timeout(TIMEOUT)),
    );
  }
}
