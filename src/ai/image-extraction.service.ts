import { ObjectId } from 'mongodb';
import { toObjectId } from '../shared/objectId.ts';
import { AiExtractionDoc, aiExtractions } from '../db/collections.ts';
import { together } from './together.ts';

export interface ExtractTextFromImagePayload {
  imageDataUrl: string;
  fileName: string;
  mimeType: string;
}

export async function extractTextFromImage(
  userId: string,
  data: ExtractTextFromImagePayload,
): Promise<AiExtractionDoc> {
  const response = await together.chat.completions.create({
    model: 'google/gemma-4-31B-it',
    messages: [
      {
        role: 'system',
        content:
          'You are an OCR text extraction assistant. Your task is to extract all visible text from the provided image as accurately as possible.',
      },
      {
        role: 'user',
        content: [
          {
            type: 'image_url',
            image_url: {
              url: data.imageDataUrl,
            },
          },
        ],
      },
    ],
  });

  const extractedText = response.choices[0].message?.content ?? '';

  const doc = {
    _id: new ObjectId(),
    userId: toObjectId(userId),
    fileName: data.fileName,
    mimeType: data.mimeType,
    extractedText,
    createdAt: new Date(),
  };

  await aiExtractions.insertOne(doc);

  return doc;
}
