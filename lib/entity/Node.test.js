import Node from './Node';

describe('Node', () => {
  test('constructor()', () => {
    const n = new Node(100);
    expect(n).toBeTruthy();
    expect(n.value).toBe(100);
    expect(n.left).toBe(null);
    expect(n.right).toBe(null);
  });

  test('value', () => {
    const n1 = new Node(100);
    expect(n1.value).toBe(100);
    const n2 = new Node('foo bar');
    expect(n2.value).toBe('foo bar');
  });

  test('left', () => {
    const n1 = new Node(100);
    const n2 = new Node(200);
    n1.left = n2;
    expect(n1.left.value).toBe(n2.value);
  });

  test('right', () => {
    const n1 = new Node(100);
    const n2 = new Node(200);
    n2.right = n1;
    expect(n2.right.value).toBe(n1.value);
  });

  test('isLeaf()', () => {
    const n1 = new Node(100);
    expect(n1.isLeaf()).toBeTruthy();
    const n2 = new Node(200);
    n1.left = n2;
    expect(n1.isLeaf()).toBeFalsy();
  });
});
