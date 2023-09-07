import { Injectable, Inject } from '@nestjs/common';
import { USER_SERVICE_NAME, User, UserServiceClient } from './user.pb';
import { ClientGrpc } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class UserService {
    private svc: UserServiceClient;

    @Inject(USER_SERVICE_NAME)
    private readonly client: ClientGrpc;

    public onModuleInit() {
        this.svc = this.client.getService<UserServiceClient>(USER_SERVICE_NAME);
    }

    public async validate(token: string): Promise<User> {
        return firstValueFrom(this.svc.validate({ token }));
    }
}
