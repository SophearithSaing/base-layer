import { requireUser } from '../auth/requireUser.ts';
import { Route } from '../routes.ts';
import { fileToDataUrl, getUploadedFile } from '../shared/files.ts';
import { HttpMethod, json, readJson } from '../shared/http.ts';
import { validateChatInput } from './ai-validation.ts';
import {
  AIModels,
  extractReceiptFromImage,
  extractTextFromImage,
} from './image-extraction.service.ts';
import { validateImageFile } from './image-validation.ts';
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
        model: AIModels.GLM_5_2,
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
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/ai/extract-text' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const file = await getUploadedFile(req, 'image');
      validateImageFile(file);
      const imageDataUrl = await fileToDataUrl(file);
      const result = await extractTextFromImage(user.userId, {
        imageDataUrl,
        fileName: file.name,
        mimeType: file.type,
      });

      return json({ result }, { status: 201 });
    },
  },
  {
    method: HttpMethod.POST,
    pattern: new URLPattern({ pathname: '/ai/extract-receipt' }),
    handler: async (req: Request) => {
      const user = await requireUser(req);
      const file = await getUploadedFile(req, 'image');
      validateImageFile(file);
      const imageDataUrl = await fileToDataUrl(file);
      const result = await extractReceiptFromImage(user.userId, {
        imageDataUrl,
        fileName: file.name,
        mimeType: file.type,
      });

      return json({ result }, { status: 201 });
    },
  },
];
