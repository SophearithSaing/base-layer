import { HttpError } from '../shared/http.ts';

export function parseExtractedReceiptObject(extractedText: string): unknown {
  try {
    return JSON.parse(extractedText);
  } catch {
    throw new HttpError('AI returned invalid receipt JSON', 502);
  }
}
