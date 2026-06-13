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
    tokenHash,
    expiresAt: getRefreshExpiry(),
    createdAt: new Date(),
  });
}

export async function findValidRefreshToken(
  token: string,
): Promise<WithId<RefreshTokenDoc> | null> {
  const candidates = await refreshTokens
    .find({
      revokedAt: { $exists: false },
      expiresAt: { $gt: new Date() },
    })
    .toArray();

  for (const doc of candidates) {
    if (await verifyPassword(token, doc.tokenHash)) {
      return doc;
    }
  }

  return null;
}

export async function revokeRefreshToken(id: ObjectId): Promise<void> {
  await refreshTokens.updateOne(
    { _id: id },
    { $set: { revokedAt: new Date() } },
  );
}
