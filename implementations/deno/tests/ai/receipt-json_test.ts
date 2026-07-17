import { assertEquals, assertThrows } from '@std/assert';
import { HttpError } from '../../src/shared/http.ts';
import { parseExtractedReceiptObject } from '../../src/ai/receipt-json.ts';

Deno.test('parseExtractedReceiptObject parses valid receipt JSON', () => {
  const result = parseExtractedReceiptObject(
    '[{"name":"Chicken Sandwich","amount":5.99}]',
  );

  assertEquals(result, [{ name: 'Chicken Sandwich', amount: 5.99 }]);
});

Deno.test('parseExtractedReceiptObject throws controlled error for invalid JSON', () => {
  const error = assertThrows(
    () => parseExtractedReceiptObject('not json'),
    HttpError,
  );

  assertEquals(error.status, 502);
  assertEquals(error.message, 'AI returned invalid receipt JSON');
});
