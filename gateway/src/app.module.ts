import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { UserModule } from './user/user.module';
import { ConfigModule } from '@nestjs/config';
import { ProductModule } from './product/product.module';

@Module({
  imports: [UserModule, ConfigModule.forRoot({ isGlobal: true }), ProductModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule { }
