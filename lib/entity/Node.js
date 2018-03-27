/**
 * BinaryTree node.
 */
class Node {
  constructor(value) {
    this._val = value;
    this.left = null;
    this.right = null;
  }

  get value() {
    return this._val;
  }

  /**
   * Returns true if node has no descendants in the tree.
   * @returns {boolean}
   */
  isLeaf() {
    return !(this.left || this.right);
  }
}

export default Node;
