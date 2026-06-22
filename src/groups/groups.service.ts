import { ObjectId, WithId } from 'mongodb';
import { GroupDoc, groups } from '../db/collections.ts';
import { HttpError } from '../shared/http.ts';
import { toObjectId } from '../shared/object-id.ts';

export type GroupPayload = {
  name: string;
  people: string[];
};

export type UpdateGroupPayload = Partial<GroupPayload>;

export async function createGroup(
  creatorId: string,
  payload: GroupPayload,
): Promise<GroupDoc> {
  const existing = await groups.findOne({ creatorId: toObjectId(creatorId) });
  if (existing) {
    throw new HttpError('User already has a group', 409);
  }

  const now = new Date();
  const group: GroupDoc = {
    _id: new ObjectId(),
    creatorId: toObjectId(creatorId),
    name: payload.name,
    people: payload.people,
    createdAt: now,
    updatedAt: now,
  };

  await groups.insertOne(group);

  return group;
}

export async function getGroupByCreatorId(
  creatorId: string,
): Promise<WithId<GroupDoc>> {
  const group = await groups.findOne({ creatorId: toObjectId(creatorId) });

  if (!group) {
    throw new HttpError('Group not found', 404);
  }

  return group;
}

export async function updateGroup(
  creatorId: string,
  payload: UpdateGroupPayload,
): Promise<WithId<GroupDoc>> {
  const result = await groups.findOneAndUpdate(
    {
      creatorId: toObjectId(creatorId),
    },
    {
      $set: { ...payload, updatedAt: new Date() },
    },
    {
      returnDocument: 'after',
    },
  );

  if (!result) {
    throw new HttpError('Group not found', 404);
  }

  return result;
}
