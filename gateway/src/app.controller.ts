import { Controller, Get } from '@nestjs/common';
import { AppService, StatusResponse } from './app.service';
import { ApiTags } from '@nestjs/swagger';

@ApiTags('Hello')
@Controller()
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getStatus(): StatusResponse {
    return this.appService.getStatus();
  }
}
