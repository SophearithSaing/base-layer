import { requireUser } from '../auth/requireUser.ts';
import { Route } from '../routes.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import { validateChatInput } from './aiValidation.ts';
import { together } from './together.ts';

interface ChatPayload {
  message: string;
}

export const aiRoutes: Route[] = [
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/ai/chat' }),
    handler: async (req: Request) => {
      await requireUser(req);
      const payload = await readJson<ChatPayload>(req);
      validateChatInput(payload);
      const response = await together.chat.completions.create({
        model: 'google/gemma-4-31B-it',
        messages: [
          {
            role: 'user',
            content: payload.message,
          },
        ],
      });

      return json({ message: response.choices[0].message?.content });
    },
  },
];
