import { ObjectId } from 'mongodb';
import { toObjectId } from '../shared/objectId.ts';
import { AiExtractionDoc, aiExtractions } from '../db/collections.ts';
import { together } from './together.ts';

export enum AIModels {
  Gemma_4 = 'google/gemma-4-31B-it',
  KimiK_2_6 = 'moonshotai/Kimi-K2.6',
  GLM_5_2 = 'zai-org/GLM-5.2',
}

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
    model: AIModels.KimiK_2_6,
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

export async function extractReceiptFromImage(
  imageDataUrl: string,
): Promise<Array<unknown>> {
  const response = await together.chat.completions.create({
    model: AIModels.KimiK_2_6,
    messages: [
      {
        role: 'system',
        content: `
          You are a receipt extraction assistant. Extract only purchased line items from the receipt image.
          Return only a valid JSON array. Each object must have exactly:

          * \`name\`: item name as shown
          * \`amount\`: final line item price as a number

          Rules:
          * Exclude subtotal, tax, discounts, tips, fees, totals, payment info, dates, store info, and receipt metadata.
          * If quantity is shown, return one object per purchased unit.
          * If the receipt shows \`2 Chicken Sandwich 11.98\`, return two objects, each with \`"amount": 5.99\`.
          * If only the total line amount is shown for multiple units, divide it by the quantity.
          * Omit items with unreadable names or prices.
          * Do not guess, explain, or add Markdown.
          * If no items are found, return [].

          Example:
          [
            {
              "name": "Chicken Sandwich",
              "amount": 5.99
            }
          ]
        `,
      },
      {
        role: 'user',
        content: [
          {
            type: 'image_url',
            image_url: {
              url: imageDataUrl,
            },
          },
        ],
      },
    ],
  });

  const extractedText = response.choices[0].message?.content ?? '{}';
  const sanitizedText = extractedText
    .trim()
    .replace(/^```json\s*/i, '')
    .replace(/^```\s*/i, '')
    .replace(/```$/i, '')
    .trim();

  return JSON.parse(sanitizedText);
}
