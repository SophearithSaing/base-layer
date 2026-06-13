import { getCookie } from '../shared/cookies.ts';
import { HttpError } from '../shared/http.ts';
import { verifyAccessToken } from './jwt.ts';

export async function requireUser(req: Request) {
  const token = getCookie(req, 'access_token');

  if (!token) {
    throw new HttpError('Unauthorized', 401);
  }

  try {
    return await verifyAccessToken(token);
  } catch {
    throw new HttpError('Invalid or expired token', 401);
  }
}
