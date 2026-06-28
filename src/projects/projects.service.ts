import { InsertOneResult, ObjectId, WithId } from 'mongodb';
import { ProjectDoc, projects } from '../db/collections.ts';
import { HttpError } from '../shared/http.ts';
import { toObjectId } from '../shared/object-id.ts';

export type ProjectPayload = Omit<
  ProjectDoc,
  '_id' | 'createdAt' | 'updatedAt'
>;

export async function listProjects(): Promise<WithId<ProjectDoc>[]> {
  return await projects.find().toArray();
}

export async function createProject(
  payload: ProjectPayload,
): Promise<InsertOneResult<ProjectDoc>> {
  const now = new Date();

  return await projects.insertOne({
    ...payload,
    _id: new ObjectId(),
    createdAt: now,
    updatedAt: now,
  });
}

export async function getProject(
  projectId: string,
): Promise<WithId<ProjectDoc>> {
  const project = await projects.findOne({
    _id: toObjectId(projectId),
  });

  if (!project) {
    throw new HttpError('Project not found', 404);
  }

  return project;
}

export async function updateProject(
  projectId: string,
  payload: ProjectPayload,
): Promise<WithId<ProjectDoc>> {
  const result = await projects.findOneAndUpdate(
    {
      _id: toObjectId(projectId),
    },
    {
      $set: { ...payload, updatedAt: new Date() },
    },
    {
      returnDocument: 'after',
    },
  );

  if (!result) {
    throw new HttpError('Project not found', 404);
  }

  return result;
}

export async function deleteProject(projectId: string): Promise<boolean> {
  const result = await projects.deleteOne({
    _id: toObjectId(projectId),
  });

  if (!result.deletedCount) {
    throw new HttpError('Project not found', 404);
  }

  return result.deletedCount === 1;
}
