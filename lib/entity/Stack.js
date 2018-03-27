import LinkedList from './LinkedList';

class Stack {
  constructor() {
    this._list = new LinkedList(null);
  }

  get count() {
    return this._list.count;
  }

  peek() {
    if (this._list.head == null) {
      return null;
    }
    return this._list.head.value;
  }

  push(val) {
    this._list.prepend(val);
  }

  pop() {
    const item = this._list.shift();
    if (item === null) return null;
    return item.value;
  }
}

export default Stack;
