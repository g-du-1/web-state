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
- 
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

  let node, visibleText = '';

  while (node = walker.nextNode()) {
    const text = node.nodeValue.trim();

    if (!text) continue;

    const range = document.createRange();
    range.selectNodeContents(node);

    for (const rect of range.getClientRects()) {
      const { top, bottom, left, right } = rect;
      const { innerWidth: w, innerHeight: h } = window;

      if (top >= 0 && bottom <= h && left >= 0 && right <= w) {
        visibleText += text + ' ';
        break;
      }
    }
  }

  return visibleText.trim();
}

console.log(getVisibleTextInViewport());
```
