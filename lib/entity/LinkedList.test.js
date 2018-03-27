import LinkedList from './LinkedList';
import ListItem from './ListItem';

const fruits = [
  'apple',
  'banana',
  'cherry',
  'date',
  'elderberry',
  'fig',
  'guava',
];

describe('LinkedList', () => {
  test('constructor()', () => {
    const ll = new LinkedList();
    expect(ll).toBeTruthy();
    expect(ll.count).toBe(0);
  });

  test('push()', () => {
    const ll = new LinkedList();
    ll.append(100);
    expect(ll.count).toBe(1);
    ll.append(200);
    expect(ll.count).toBe(2);
  });

  test('shift()', () => {
    const ll = new LinkedList(fruits);
    expect(ll).toBeTruthy();
    expect(ll.count).toBe(fruits.length);
    expect(ll.shift().value).toBe(fruits[0]);
    expect(ll.shift().value).toBe(fruits[1]);
    expect(ll.shift().value).toBe(fruits[2]);
  });

  test('head()', () => {
    const ll = new LinkedList(fruits);
    const item = ll.head;
    expect(ll).toBeTruthy();
    expect(item instanceof ListItem).toBeTruthy();
    expect(item.value).toBe('apple');
    expect(item.next.value).toBe('banana');
  });

  test('removeItem()', () => {
    const ll = new LinkedList(fruits);
    const item = ll.head.next.next;
    expect(item instanceof ListItem).toBeTruthy();
    expect(item.value).toBe('cherry');
    ll.removeItem(item);
    expect(ll.head.next.next.value).toBe('date');
  });

  test('insert()', () => {
    const ll = new LinkedList(fruits);
    const item = ll.head.next.next;
    expect(item instanceof ListItem).toBeTruthy();
    expect(item.value).toBe('cherry');
    ll.insert(item, 'blueberry');
    expect(ll.count).toBe(fruits.length + 1);
    expect(ll.head.next.next.value).toBe('blueberry');
    expect(ll.head.next.next.next.value).toBe('cherry');
  });
});
