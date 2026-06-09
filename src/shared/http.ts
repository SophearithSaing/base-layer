export enum HttpMethod {
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
  PATCH = 'PATCH',
  DELETE = 'DELETE',
  OPTIONS = 'OPTIONS',
  HEAD = 'HEAD',
}

export type Handler = (
  req: Request,
  params: Record<string, string>,
) => Response | Promise<Response>;

export function json(data: unknown, init: ResponseInit = {}) {
  return Response.json(data, {
    ...init,
    headers: {
      'content-type': 'application/json',
      ...(init.headers ?? {}),
    },
  });
}

export function error(message: string, status = 400) {
  return json({ error: message }, { status });
}

export async function readJson<T>(req: Request): Promise<T> {
  try {
    return (await req.json()) as T;
  } catch {
    throw new HttpError('Invalid JSON body', 400);
  }
}

export class HttpError extends Error {
  constructor(
    message: string,
    public status = 500,
  ) {
    super(message);
  }
}
