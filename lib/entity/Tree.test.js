/* eslint-disable no-new */
import Tree from './Tree';
import Node from './Node';

let tree;

beforeAll(() => {
  tree = new Tree(new Node('F'));
  tree.root.left = new Node('B');
  tree.root.left.left = new Node('A');
  tree.root.left.right = new Node('D');
  tree.root.left.right.left = new Node('C');
  tree.root.left.right.right = new Node('E');
  tree.root.right = new Node('G');
  tree.root.right.right = new Node('I');
  tree.root.right.right.left = new Node('H');
});

describe('Tree', () => {
  describe('constructor()', () => {
    test('one node tree', () => {
      const t1 = new Tree(new Node(100));
      expect(t1).toBeTruthy();
      expect(t1.root).toBeDefined();
    });

    test('invalid arg', () => {
      expect(() => {
        new Tree(100);
      }).toThrow(TypeError);
    });
  });

  describe('traversal', () => {
    let arr;

    function callback(node) {
      arr.push(node.value);
    }

    beforeEach(() => {
      arr = [];
    });

    test('three-node tree', () => {
      const t1 = new Tree(new Node(10));
      t1.root.left = new Node(5);
      t1.root.right = new Node(15);
      t1.traverseIn(callback);
      expect(arr).toEqual([5, 10, 15]);
    });

    test('traverseIn tree', () => {
      tree.traverseIn(callback);
      expect(arr).toEqual(['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I']);
    });

    test('traversePre(cb)', () => {
      tree.traversePre(callback);
      expect(arr).toEqual(['F', 'B', 'A', 'D', 'C', 'E', 'G', 'I', 'H']);
    });

    test('traversePost(cb)', () => {
      tree.traversePost(callback);
      expect(arr).toEqual(['A', 'C', 'E', 'D', 'B', 'H', 'I', 'G', 'F']);
    });

    test('bad traverse arg', () => {
      expect(() => {
        tree.traverseIn(10);
      }).toThrow(TypeError);
      expect(() => {
        tree.traversePre(new Node(20));
      }).toThrow(TypeError);
      expect(() => {
        tree.traversePost([1, 2, 3]);
      }).toThrow(TypeError);
    });
  });
});
