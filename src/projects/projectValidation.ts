import { HttpError } from '../shared/http.ts';

export function validateCreateProject(data: {
  name?: unknown;
  description?: unknown;
}) {
  if (typeof data.name !== 'string' || data.name.trim().length === 0) {
    throw new HttpError('Project name is required', 400);
  }
}

export function validateUpdateProject(data: {
  name?: unknown;
  description?: unknown;
}) {
  if (data.name !== undefined) {
    if (typeof data.name !== 'string' || data.name.trim().length === 0) {
      throw new HttpError('Project name cannot be empty', 400);
    }
  }

  if (data.description !== undefined) {
    if (typeof data.description !== 'string') {
      throw new HttpError('Project description must be a string', 400);
    }
  }
}
