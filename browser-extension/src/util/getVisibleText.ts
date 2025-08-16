export const getVisibleText = () => {
  const viewport = {
    top: window.scrollY,
    left: window.scrollX,
    width: window.innerWidth,
    height: window.innerHeight,
  };

  const topMargin = viewport.height * 0.1;
  const bottomMargin = viewport.height * 0.9;

  const walker = document.createTreeWalker(
    document.body,
    NodeFilter.SHOW_TEXT,
    null
  );

  let visibleText = [];

  let node;

  while ((node = walker.nextNode())) {
    if (node.parentElement) {
      const rect = node.parentElement.getBoundingClientRect();

      if (
        rect.top >= topMargin &&
        rect.bottom <= bottomMargin &&
        rect.left < viewport.width &&
        rect.right > 0
      ) {
        const text = node.textContent?.trim();

        if (text && text.length > 25) {
          visibleText.push(text + "\n\n");
        }
      }
    }
  }

  return visibleText.join(" ");
};
