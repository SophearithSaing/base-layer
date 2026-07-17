import { HttpError } from '../shared/http.ts';

type ChatPayload = {
  message?: unknown;
};

export function validateChatInput(data: ChatPayload) {
  if (typeof data.message !== 'string' || data.message.trim().length === 0) {
    throw new HttpError('Message is required', 400);
  }
}
