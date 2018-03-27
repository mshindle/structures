/**
 * ListItem wraps a value with pointers to another item either
 * next or previous in the list
 */
class ListItem {
  /**
   * ListItem constructor
   * @param val
   */
  constructor(val) {
    this._val = val;
    this.next = null;
    this.prev = null;
  }

  get value() {
    return this._val;
  }
}

export default ListItem;
