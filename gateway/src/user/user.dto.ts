import { Timestamp } from "./google/protobuf/timestamp.pb";
import { IsEmail, Length } from 'class-validator';

export class UserDto {
    id: string;
    name: string;
    email: string;
    createdAt: Timestamp | undefined;
    updatedAt: Timestamp | undefined;
    deletedAt: Timestamp | undefined;
}

export class UsersDto {
    users: UserDto[];
}

export class RegisterUserDto {
    @IsEmail()
    email: string;

    @Length(3)
    name: string;

    @Length(8)
    password: string;
}

export class LoginUserDto {
    @IsEmail()
    email: string;

    @Length(8)
    password: string;
}

export class TokenDto {
    token: string;
}

export class UsersQueryDto {
    page: number;
    size: number;
}