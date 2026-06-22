import { ObjectId, WithId } from 'mongodb';
import { hashPassword, verifyPassword } from './password.ts';
import { RefreshTokenDoc, refreshTokens } from '../db/collections.ts';

const REFRESH_TOKEN_DAYS = 30;

export function createRefreshToken(): string {
  return crypto.randomUUID() + '.' + crypto.randomUUID();
}

export function getRefreshExpiry(): Date {
  const expiresAt = new Date();
  expiresAt.setDate(expiresAt.getDate() + REFRESH_TOKEN_DAYS);
  return expiresAt;
}

export async function saveRefreshToken(
  userId: ObjectId,
  token: string,
): Promise<void> {
  const tokenHash = await hashPassword(token);

  await refreshTokens.insertOne({
    _id: new ObjectId(),
    userId,
    tokenLookupHash: await createTokenLookupHash(token),
    tokenHash,
    expiresAt: getRefreshExpiry(),
    createdAt: new Date(),
  });
}

export async function findValidRefreshToken(
  token: string,
): Promise<WithId<RefreshTokenDoc> | null> {
  const tokenLookupHash = await createTokenLookupHash(token);
  const doc = await refreshTokens.findOne({
    tokenLookupHash,
    revokedAt: { $exists: false },
    expiresAt: { $gt: new Date() },
  });

  if (!doc) {
    return null;
  }

  const valid = await verifyPassword(token, doc.tokenHash);

  return valid ? doc : null;
}

export async function revokeRefreshToken(id: ObjectId): Promise<void> {
  await refreshTokens.updateOne(
    { _id: id },
    { $set: { revokedAt: new Date() } },
  );
}

async function createTokenLookupHash(token: string): Promise<string> {
  const data = new TextEncoder().encode(token);
  const digest = await crypto.subtle.digest('SHA-256', data);

  return [...new Uint8Array(digest)]
    .map((byte) => byte.toString(16).padStart(2, '0'))
    .join('');
}
