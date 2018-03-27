import Node from './Node';

class Tree {
  /**
   * Constructs a binary tree using node as the root
   * @param {Node} node
   */
  constructor(node) {
    if (!(node instanceof Node)) {
      throw new TypeError('Tree can only be created with a Node object');
    }
    this._root = node;
  }

  /**
   * Returns the beginning of the tree.
   * @returns {Node|*}
   */
  get root() {
    return this._root;
  }

  traversePre(cb) {
    if (typeof cb !== 'function') {
      throw new TypeError('callback function expected as arg.');
    }
    this._preOrder(this._root, cb);
  }

  _preOrder(node, cb) {
    if (!node) {
      return;
    }
    cb(node);
    this._preOrder(node.left, cb);
    this._preOrder(node.right, cb);
  }

  traverseIn(cb) {
    if (typeof cb !== 'function') {
      throw new TypeError('callback function expected as arg.');
    }
    this._inOrder(this._root, cb);
  }

  _inOrder(node, cb) {
    if (!node) {
      return;
    }
    this._inOrder(node.left, cb);
    cb(node);
    this._inOrder(node.right, cb);
  }

  traversePost(cb) {
    if (typeof cb !== 'function') {
      throw new TypeError('callback function expected as arg.');
    }
    this._postOrder(this._root, cb);
  }

  _postOrder(node, cb) {
    if (!node) {
      return;
    }
    this._postOrder(node.left, cb);
    this._postOrder(node.right, cb);
    cb(node);
  }
}

export default Tree;
