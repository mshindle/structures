import ListItem from './ListItem';

class LinkedList {
  /**
   * LinkedList constructor
   * The zero element in the array will be first element in the _list.
   *
   * @param {Array} arr If an array is provided, linked _list will contain array elements.
   */
  constructor(arr) {
    this._head = null;
    this._tail = null;
    this._count = 0;
    if (arr && Array.isArray(arr)) {
      for (let i = arr.length - 1; i >= 0; i--) {
        this.prepend(arr[i]);
      }
    }
  }

  /**
   * Create a new item which is unattached to the List.
   *
   * @param val
   * @returns {ListItem}
   */
  static createItem(val) {
    return val instanceof ListItem ? val : new ListItem(val);
  }

  get count() {
    return this._count;
  }

  get head() {
    return this._head;
  }

  get tail() {
    return this._tail;
  }

  /**
   *
   * @param item
   * @param val
   * @returns {LinkedList}
   */
  insert(before, val) {
    if (!(before instanceof ListItem)) {
      throw new TypeError('item should be of type ListItem');
    }
    const newItem = LinkedList.createItem(val);
    const prev = before.prev;
    if (prev) {
      prev.next = newItem;
    }
    before.prev = newItem;
    newItem.next = before;
    this._count++;
    return newItem;
  }

  shift() {
    const head = this._head;
    if (head) {
      this._head = head.next;
    }
    if (this._head == null) {
      this._tail = null;
    }
    this._count--;
    return head;
  }

  prepend(val) {
    const newItem = LinkedList.createItem(val);
    if (this._head) {
      this._head.prev = newItem;
      newItem.next = this._head;
    } else {
      this._tail = newItem;
    }
    this._head = newItem;
    this._count++;
    return this;
  }

  append(val) {
    const newItem = LinkedList.createItem(val);
    if (this._tail) {
      this._tail.next = newItem;
      newItem.prev = this._tail;
    } else {
      this._head = newItem;
    }
    this._tail = newItem;
    this._count++;
    return this;
  }

  /**
   * Remove item from linked _list.
   * @param item
   * @returns {LinkedList}
   */
  removeItem(item) {
    if (!item && !(item instanceof ListItem)) {
      throw new TypeError('item should be of type ListItem');
    }
    const prev = item.prev;
    const next = item.next;
    if (prev) {
      prev.next = next;
    } else {
      // removed item is head.
      this._head = item.next;
      if (this._head) {
        this._head.prev = null;
      }
    }
    if (next) {
      next.prev = prev;
    } else {
      this._tail = item.prev;
      if (this._tail) {
        this._tail.next = null;
      }
    }
    this._count--;
    return this;
  }
}

module.exports = LinkedList;
