import { HttpError } from './http.ts';

export async function getUploadedFile(
  req: Request,
  fieldName = 'image',
): Promise<File> {
  const formData = await req.formData();
  const file = formData.get(fieldName);

  if (!(file instanceof File)) {
    throw new HttpError(`Missing file field: ${fieldName}`, 400);
  }

  return file;
}

export async function fileToDataUrl(file: File): Promise<string> {
  const bytes = new Uint8Array(await file.arrayBuffer());
  const base64 = encodeBase64(bytes);

  return `data:${file.type};base64,${base64}`;
}

function encodeBase64(bytes: Uint8Array): string {
  let binary = '';

  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }

  return btoa(binary);
}
