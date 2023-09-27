import {
  ExceptionFilter,
  Catch,
  ArgumentsHost,
  HttpException,
  HttpStatus,
} from '@nestjs/common';
import { HttpAdapterHost } from '@nestjs/core';

@Catch()
export class AllExceptionFilter implements ExceptionFilter {
  constructor(private readonly httpAdapterHost: HttpAdapterHost) {}

  catch(exception: unknown, host: ArgumentsHost): void {
    // In certain situations `httpAdapter` might not be available in the
    // constructor method, thus we should resolve it here.
    const { httpAdapter } = this.httpAdapterHost;

    console.log(exception);

    const ctx = host.switchToHttp();

    if (exception instanceof HttpException) {
      httpAdapter.reply(
        ctx.getResponse(),
        exception.getResponse(),
        exception.getStatus(),
      );
      return;
    }

    var httpStatus = HttpStatus.FAILED_DEPENDENCY;

    const responseBody = {
      statusCode: httpStatus,
      error: 'Failed Dependency',
      message: 'Unknown error',
    };

    if (
      typeof exception === 'object' &&
      exception !== null &&
      exception.hasOwnProperty('details')
    ) {
      responseBody.message = exception['details'];
    }

    httpAdapter.reply(ctx.getResponse(), responseBody, httpStatus);
  }
}
