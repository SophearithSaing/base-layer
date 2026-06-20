import { HttpError } from '../shared/http.ts';

const allowedTypes = ['image/png', 'image/jpeg', 'image/webp'];
const maxSize = 5 * 1024 * 1024;

export function validateImageFile(file: File) {
  if (!allowedTypes.includes(file.type)) {
    throw new HttpError('Image must be PNG, JPEG, or WebP', 400);
  }

  if (file.size > maxSize) {
    throw new HttpError('Image must be under 5MB', 400);
  }
}
