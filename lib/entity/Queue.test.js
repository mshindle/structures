/* eslint-disable no-new */
import Queue from './Queue';

const fruits = [
  'apple',
  'banana',
  'cherry',
  'date',
  'elderberry',
  'fig',
  'guava',
];

describe('Queue', () => {
  test('constructor()', () => {
    const q = new Queue();
    expect(q).toBeTruthy();
    expect(q.count).toBe(0);
  });

  test('peek()', () => {
    const q = new Queue();
    expect(q.peek()).toBeNull();
    q.push(fruits[0]);
    expect(q.count).toBe(1);  // should only be 1 item
    expect(q.peek()).toBe(fruits[0]); // value should be first item in array
    expect(q.count).toBe(1); // should still be 1 item
  });

  test('push()', () => {
    const q = new Queue();
    const l = fruits.length;
    for (let i = 0; i < l; i++) {
      q.push(fruits[i]);
    }
    expect(q.count).toBe(l);
    for (let i = 0; i < l; i++) {
      expect(q.pop()).toBe(fruits[i]);
    }
    expect(q.pop()).toBeNull();
  });
});
