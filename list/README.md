## Thought process for building List in this format

* **Stable Pointers:** Unlike a slice, `PushBack` or `PushFront` will never move existing nodes in memory. This is critical for your drone system if you're holding onto references of pending tasks.
* **Predictable Latency:** No "growth spikes." Every insertion is just one small allocation and a pointer swap.
* **Clean API:** We've avoided the `interface{}` trap of the standard library's `container/list`.

### Looking ahead to v0.0.3 Release Notes
When we're ready to tag v0.0.3, the headline will be:
> **"Linear Performance & Stable Memory: Introducing Singly Linked Lists"**
> * Featuring $O(1)$ stack and queue operations.
> * Fully generic and thread-safe.
> * Native iterator support for sequential processing.

### Potential features to add
 * A `Wait/Notify` mechanism (like a `Cond` variable) so goroutines can block until an item is available in the list