import { Injectable, CanActivate, ExecutionContext, HttpStatus, UnauthorizedException, Inject } from '@nestjs/common';
import { Request } from 'express';
import { User } from './user.pb';
import { UserService } from './user.service';

export interface UserRequest extends Request {
    user: string;
}

@Injectable()
export class UserGuard implements CanActivate {
    @Inject(UserService)
    public readonly service: UserService;

    public async canActivate(ctx: ExecutionContext): Promise<boolean> | never {
        const req: UserRequest = ctx.switchToHttp().getRequest();
        const authorization: string = req.headers['authorization'];

        if (!authorization) {
            throw new UnauthorizedException();
        }

        const bearer: string[] = authorization.split(' ');

        if (!bearer || bearer.length < 2) {
            throw new UnauthorizedException();
        }

        const token: string = bearer[1];

        const { id }: User = await this.service.validate(token);

        if (!id) {
            throw new UnauthorizedException();
        }

        req.user = id;

        return true;
    }
}