import quickSort from './quickSort';

function sort(method) {
  test('empty', () => {
    expect(method([])).toEqual([]);
  });

  test('one or two elements', () => {
    expect(method([1])).toEqual([1]);
    expect(method([2, 1])).toEqual([1, 2]);
    expect(method([1, 2])).toEqual([1, 2]);
  });
}

describe('quicksort', () => {
  sort(quickSort);
});
