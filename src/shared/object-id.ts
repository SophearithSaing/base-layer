import { ObjectId } from 'mongodb';
import { HttpError } from './http.ts';

export function toObjectId(id: string | ObjectId) {
  if (!ObjectId.isValid(id)) {
    throw new HttpError('Invalid id', 400);
  }

  return new ObjectId(id);
}
