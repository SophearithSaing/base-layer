import { requireUser } from '../auth/require-user.ts';
import { Route } from '../routes.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import {
  validateCreateGroup,
  validateUpdateGroup,
} from './group-validation.ts';
import {
  createGroup,
  getGroupByCreatorId,
  updateGroup,
} from './groups.service.ts';

export const groupRoutes: Route[] = [
  {
    method: HttpMethod.GET,
    pattern: new URLPattern({ pathname: '/groups' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const group = await getGroupByCreatorId(user.userId);

      return json({ group });
    },
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/groups/create' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const payload = validateCreateGroup(await readJson(req));
      const group = await createGroup(user.userId, payload);

      return json({ group }, { status: 201 });
    },
  },
  {
    method: HttpMethod.PATCH,
    pattern: new URLPattern({ pathname: '/groups' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const payload = validateUpdateGroup(await readJson(req));
      const group = await updateGroup(user.userId, payload);

      return json({ group });
    },
  },
];
