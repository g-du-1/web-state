# page-state-saver

## LLM Base Query

You are an experienced senior software engineer. Help me with designing a system. I have a starting point like:

- Functional Requirements
- NFRs
- Technologies
- Diagrams

What else do I need?

It will be a system that saves the state of browser tabs in a database. Like scroll position, number of scroll more buttons clicked etc. It should solve the issue when a mobile browser refreshes the tab on reddit and losing where it was. I need to be able to load back or see where I was on the tab. Make it as simple as possible but also as decoupled as possible. It should be a server that is called by a mobile browser extension (chrome). 

## Functional Requirements
- Save web page state to database through a mobile browser extension (or initially a userscript)
- When the page stops scrolling, save (needs to be general for all sites):
  - Visible text in the viewport
  - Scrollposition
## NFRs
- KISS
### Logging
## Technologies
- Spring Boot
- Postgres
- Docker
- Chrome Extension - Vanilla JS - as simple as possible
## Diagrams


## Example Code

```js
function getVisibleTextInViewport() {
  const walker = document.createTreeWalker(document.body, NodeFilter.SHOW_TEXT);
  let text = '', node;

  while (node = walker.nextNode()) {
    const value = node.nodeValue.trim();

    if (!value) continue;

    const range = document.createRange();
    range.selectNodeContents(node);
    const rects = range.getClientRects();

    for (const { top, bottom, left, right } of rects) {
      if (top >= 0 && bottom <= innerHeight && left >= 0 && right <= innerWidth) {
        text += value + ' ';
        break;
      }
    }
  }

  return text.trim();
}

console.log(getVisibleTextInViewport());
```
