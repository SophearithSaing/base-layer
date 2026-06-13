import { ObjectId } from 'mongodb';
import { signAccessToken } from './jwt.ts';
import { users } from '../db/collections.ts';
import {
  createRefreshToken,
  findValidRefreshToken,
  revokeRefreshToken,
  saveRefreshToken,
} from './refreshTokens.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import { clearCookie, getCookie, serializeCookie } from '../shared/cookies.ts';
import { requireUser } from './requireUser.ts';
import { hashPassword, verifyPassword } from './password.ts';
import { HttpError } from '../shared/http.ts';
import { Route } from '../routes.ts';
type AuthRequestPayload = {
  email: string;
  password: string;
};

const ACCESS_MAX_AGE = 60 * 15;
const REFRESH_MAX_AGE = 60 * 60 * 24 * 30;

export const authRoutes: Route[] = [
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/auth/register' }),
    handler: register,
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/auth/login' }),
    handler: login,
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/auth/logout' }),
    handler: logout,
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/auth/refresh' }),
    handler: refresh,
  },
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/me' }),
    handler: me,
  },
];

async function register(req: Request): Promise<Response> {
  const { email, password } = await readJson<AuthRequestPayload>(req);
  const existing = await users.findOne({ email });

  if (existing) {
    throw new HttpError('Email is already registered', 409);
  }

  const now = new Date();
  const userId = new ObjectId();

  await users.insertOne({
    _id: userId,
    email,
    passwordHash: await hashPassword(password),
    createdAt: now,
    updatedAt: now,
  });

  return issueAuthResponse(userId, email);
}

async function login(req: Request): Promise<Response> {
  const { email, password } = await readJson<AuthRequestPayload>(req);
  const user = await users.findOne({ email });

  if (!user) {
    throw new HttpError('Invalid credentials', 401);
  }

  const valid = await verifyPassword(password, user.passwordHash);
  if (!valid) {
    throw new HttpError('Invalid credentials', 401);
  }

  return issueAuthResponse(user._id, email);
}

async function refresh(req: Request): Promise<Response> {
  const oldToken = getCookie(req, 'refresh_token');
  if (!oldToken) {
    throw new HttpError('Unauthorized', 401);
  }

  const existing = await findValidRefreshToken(oldToken);
  if (!existing) {
    throw new HttpError('Unauthorized', 401);
  }

  await revokeRefreshToken(existing._id);

  const user = await users.findOne({ _id: existing.userId });
  if (!user) {
    throw new HttpError('Unauthorized', 401);
  }

  return issueAuthResponse(user._id, user.email);
}

async function logout(req: Request): Promise<Response> {
  const token = getCookie(req, 'refresh_token');

  if (token) {
    const refreshToken = await findValidRefreshToken(token);
    if (refreshToken) {
      await revokeRefreshToken(refreshToken._id);
    }
  }

  return json(
    { ok: true },
    {
      headers: {
        'set-cookie': [
          clearCookie('access_token'),
          clearCookie('refresh_token'),
        ].join(', '),
      },
    },
  );
}

async function me(req: Request): Promise<Response> {
  const user = await requireUser(req);

  return json({ user });
}

async function issueAuthResponse(
  userId: ObjectId,
  email: string,
): Promise<Response> {
  const accessToken = await signAccessToken({
    userId: userId.toString(),
    email,
  });

  const refreshToken = createRefreshToken();
  await saveRefreshToken(userId, refreshToken);

  return json({
    user: { id: userId.toString(), email },
    headers: {
      'set-cookie': [
        serializeCookie('access_token', accessToken, {
          httpOnly: true,
          sameSite: 'Lax',
          path: '/',
          maxAge: ACCESS_MAX_AGE,
        }),
        serializeCookie('refresh_token', refreshToken, {
          httpOnly: true,
          sameSite: 'Lax',
          path: '/auth/refresh',
          maxAge: REFRESH_MAX_AGE,
        }),
      ].join(', '),
    },
  });
}
