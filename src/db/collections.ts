import { ObjectId } from 'mongodb';
import { db } from './mongo.ts';

export type UserDoc = {
  _id: ObjectId;
  email: string;
  passwordHash: string;
  createdAt: Date;
  updatedAt: Date;
};

export type RefreshTokenDoc = {
  _id: ObjectId;
  userId: ObjectId;
  tokenLookupHash: string;
  tokenHash: string;
  expiresAt: Date;
  revokedAt?: Date;
  createdAt: Date;
};

export type ProjectDoc = {
  _id: ObjectId;
  userId: ObjectId;
  name: string;
  description?: string;
  createdAt: Date;
  updatedAt: Date;
};

export type AiMessageDoc = {
  _id: ObjectId;
  userId: ObjectId;
  role: 'user' | 'assistant' | 'system';
  content: string;
  createdAt: Date;
};

export type AiExtractionDoc = {
  _id: ObjectId;
  userId: ObjectId;
  fileName?: string;
  mimeType: string;
  extractedText: string;
  createdAt: Date;
};

export const users = db.collection<UserDoc>('users');
export const refreshTokens = db.collection<RefreshTokenDoc>('refresh_tokens');
export const projects = db.collection<ProjectDoc>('projects');
export const aiMessages = db.collection<AiMessageDoc>('ai_messages');
export const aiExtractions = db.collection<AiExtractionDoc>('ai_extractions');

export async function ensureIndexes() {
  await users.createIndex({ email: 1 }, { unique: true });

  await refreshTokens.createIndex({ tokenLookupHash: 1 }, { unique: true });
  await refreshTokens.createIndex({ expiresAt: 1 }, { expireAfterSeconds: 0 });
  await refreshTokens.createIndex({ userId: 1 });

  await projects.createIndex({ userId: 1 });

  await aiExtractions.createIndex({ userId: 1, createdAt: -1 });
}
