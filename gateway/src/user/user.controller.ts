import {
  Body,
  Controller,
  Get,
  Inject,
  Param,
  Post,
  Query,
  Req,
  UseGuards,
  UseInterceptors,
} from '@nestjs/common';
import { USER_SERVICE_NAME, UserServiceClient } from './user.pb';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable, TimeoutError, throwError, timeout } from 'rxjs';
import {
  LoginUserDto,
  RegisterUserDto,
  TokenDto,
  UserDto,
  UsersDto,
  UsersQueryDto,
} from './user.dto';
import { ApiBearerAuth, ApiTags } from '@nestjs/swagger';
import { UserGuard, UserRequest } from './user.guard';
import { CacheInterceptor } from '@nestjs/cache-manager';

var nrOfTimeouts = 0;

const TIMEOUT = {
  each: 2000,
  with: () =>
    throwError(() => {
      nrOfTimeouts++;
      if (nrOfTimeouts > 2) {
        console.error('User services timeout 3 times');
      }
      return new TimeoutError();
    }),
};

@ApiTags('User')
@Controller('users')
@UseInterceptors(CacheInterceptor)
export class UserController {
  private svc: UserServiceClient;

  @Inject(USER_SERVICE_NAME)
  private readonly client: ClientGrpc;

  public onModuleInit(): void {
    this.svc = this.client.getService<UserServiceClient>(USER_SERVICE_NAME);
  }

  @Post('register')
  private async register(
    @Body() body: RegisterUserDto,
  ): Promise<Observable<TokenDto>> {
    return this.svc.register(body).pipe(timeout(TIMEOUT));
  }

  @Post('login')
  private async login(
    @Body() body: LoginUserDto,
  ): Promise<Observable<TokenDto>> {
    return this.svc.login(body).pipe(timeout(TIMEOUT));
  }

  @Get()
  private async getAll(
    @Query() query: UsersQueryDto,
  ): Promise<Observable<UsersDto>> {
    var pbQuery = {
      page: query.page || 0,
      size: query.size || 0,
    };

    return this.svc.getAll(pbQuery).pipe(timeout(TIMEOUT));
  }

  @ApiBearerAuth()
  @Get('current')
  @UseGuards(UserGuard)
  private async createProduct(
    @Req() req: UserRequest,
  ): Promise<Observable<UserDto>> {
    var id = req.user;
    return this.svc.getById({ id }).pipe(timeout(TIMEOUT));
  }

  @Get(':id')
  private async getById(@Param('id') id: string): Promise<Observable<UserDto>> {
    return this.svc.getById({ id }).pipe(timeout(TIMEOUT));
  }
}
