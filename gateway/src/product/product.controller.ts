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
  Delete,
} from '@nestjs/common';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';
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

@Controller('product')
export class ProductController implements OnModuleInit {
  private svc: ProductServiceClient;

  @Inject(PRODUCT_SERVICE_NAME)
  private readonly client: ClientGrpc;

  public onModuleInit(): void {
    this.svc =
      this.client.getService<ProductServiceClient>(PRODUCT_SERVICE_NAME);
  }

  @Post()
  @UseGuards(UserGuard)
  private async createProduct(
    @Req() req: UserRequest,
    @Body() body: CreateProductDto,
  ): Promise<Observable<ProductIdDto>> {
    var ownerId = req.user;
    return this.svc.create({ ownerId, ...body });
  }

  @Put(':id')
  @UseGuards(UserGuard)
  private async updateProduct(
    @Req() req: UserRequest,
    @Param('id') id: string,
    @Body() body: UpdateProductDto,
  ): Promise<Observable<Empty>> {
    var ownerId = req.user;
    return this.svc.update({ id, ownerId, ...body });
  }

  @Delete(':id')
  @UseGuards(UserGuard)
  private async deleteProduct(
    @Req() req: UserRequest,
    @Param('id') id: string,
  ): Promise<Observable<Empty>> {
    var ownerId = req.user;
    return this.svc.delete({ id, ownerId });
  }

  @Get()
  private async getAll(
    @Query() query: ProductsQueryDto,
  ): Promise<Observable<ProductsDto>> {
    return this.svc.getAll(query);
  }

  @Get(':id')
  private async getById(
    @Param('id') id: string,
  ): Promise<Observable<ProductDto>> {
    return this.svc.getById({ id });
  }
}
