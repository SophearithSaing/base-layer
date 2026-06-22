import { jwtVerify, SignJWT } from 'jose';
import { env } from '../config/env.ts';

const secret = new TextEncoder().encode(env.JWT_SECRET);

export type AccessTokenPayload = {
  userId: string;
  email: string;
};

export async function signAccessToken(
  payload: AccessTokenPayload,
): Promise<string> {
  return await new SignJWT(payload)
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('15m')
    .sign(secret);
}

export async function verifyAccessToken(
  token: string,
): Promise<AccessTokenPayload> {
  const result = await jwtVerify<AccessTokenPayload>(token, secret);
  return result.payload;
}
