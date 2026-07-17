import { HttpError } from '../shared/http.ts';

export type GroupPayload = {
  name?: unknown;
  people?: unknown;
};

export function validateCreateGroup(data: GroupPayload): {
  name: string;
  people: string[];
} {
  if (typeof data.name !== 'string' || data.name.trim().length === 0) {
    throw new HttpError('Group name is required', 400);
  }

  if (!Array.isArray(data.people)) {
    throw new HttpError('Group people must be an array', 400);
  }

  if (!data.people.every((person) => typeof person === 'string')) {
    throw new HttpError('Group people must contain only strings', 400);
  }

  return {
    name: data.name.trim(),
    people: data.people.map((person) => person.trim()).filter(Boolean),
  };
}

export function validateUpdateGroup(data: GroupPayload): {
  name?: string;
  people?: string[];
} {
  const payload: { name?: string; people?: string[] } = {};

  if (data.name !== undefined) {
    if (typeof data.name !== 'string' || data.name.trim().length === 0) {
      throw new HttpError('Group name cannot be empty', 400);
    }

    payload.name = data.name.trim();
  }

  if (data.people !== undefined) {
    if (!Array.isArray(data.people)) {
      throw new HttpError('Group people must be an array', 400);
    }

    if (!data.people.every((person) => typeof person === 'string')) {
      throw new HttpError('Group people must contain only strings', 400);
    }

    payload.people = data.people.map((person) => person.trim()).filter(Boolean);
  }

  return payload;
}
