import { Injectable } from '@nestjs/common';

export class StatusResponse {
  status: string;
}

@Injectable()
export class AppService {
  getStatus(): StatusResponse {
    return { status: 'ok' };
  }
}
