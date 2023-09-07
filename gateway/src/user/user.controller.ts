import { Body, Controller, Get, Inject, Param, Post, Query } from '@nestjs/common';
import { USER_SERVICE_NAME, UserServiceClient, RegisterUser, LoginUser, Token, Users, UsersQuery, User } from './user.pb';
import { ClientGrpc } from '@nestjs/microservices';
import { Observable } from 'rxjs';
import { LoginUserDto, RegisterUserDto, TokenDto, UserDto, UsersDto, UsersQueryDto } from './user.dto';

@Controller('users')
export class UserController {
    private svc: UserServiceClient;

    @Inject(USER_SERVICE_NAME)
    private readonly client: ClientGrpc;

    public onModuleInit(): void {
        this.svc = this.client.getService<UserServiceClient>(USER_SERVICE_NAME);
    }

    @Post('register')
    private async register(@Body() body: RegisterUserDto): Promise<Observable<TokenDto>> {
        return this.svc.register(body);
    }

    @Post('login')
    private async login(@Body() body: LoginUserDto): Promise<Observable<TokenDto>> {
        return this.svc.login(body);
    }

    @Get()
    private async getAll(@Query() query: UsersQueryDto): Promise<Observable<UsersDto>> {
        return this.svc.getAll(query);
    }

    @Get(':id')
    private async getById(@Param('id') id: string): Promise<Observable<UserDto>> {
        return this.svc.getById({ id });
    }
}
