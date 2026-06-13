const isProd = Deno.env.get('APP_ENV') === 'prod';

type CookieOptions = {
  httpOnly?: boolean;
  secure?: boolean;
  sameSite?: 'Lax' | 'Strict' | 'None';
  path?: string;
  maxAge?: number;
};

export function serializeCookie(
  name: string,
  value: string,
  options: CookieOptions,
): string {
  const parts = [`${name}=${value}`];

  if (options.httpOnly) {
    parts.push('HttpOnly');
  }

  if (options.secure ?? isProd) {
    parts.push('Secure');
  }

  if (options.sameSite) {
    parts.push(options.sameSite);
  }

  if (options.path) {
    parts.push(`Path=${options.path}`);
  }

  if (options.maxAge !== undefined) {
    parts.push(`Max-Age=${options.maxAge}`);
  }

  return parts.join(';');
}

export function clearCookie(name: string, path = '/'): string {
  return serializeCookie(name, '', {
    path,
    maxAge: 0,
    httpOnly: true,
    sameSite: 'Lax',
  });
}

export function getCookie(req: Request, name: string): string | undefined {
  const cookie = req.headers.get('cookie');

  if (!cookie) {
    return undefined;
  }

  return cookie
    .split(';')
    .map((part) => part.trim())
    .find((part) => part.startsWith(`${name}=`))
    ?.split('=')[1];
}
