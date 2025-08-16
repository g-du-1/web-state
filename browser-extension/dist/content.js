// src/util/getVisibleText.ts
var getVisibleText = () => {
  const viewport = {
    top: window.scrollY,
    left: window.scrollX,
    width: window.innerWidth,
    height: window.innerHeight
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
  while (node = walker.nextNode()) {
    if (node.parentElement) {
      const rect = node.parentElement.getBoundingClientRect();
      if (rect.top >= topMargin && rect.bottom <= bottomMargin && rect.left < viewport.width && rect.right > 0) {
        const text = node.textContent?.trim();
        if (text && text.length > 25) {
          visibleText.push(text + "\n\n");
        }
      }
    }
  }
  return visibleText.join(" ");
};

// src/util/isUrlDisallowed.ts
var isUrlDisallowed = async (url) => {
  const { whitelistSites } = await chrome.storage.local.get("whitelistSites");
  if (whitelistSites === "") return true;
  let result = true;
  const whitelistedArr = whitelistSites?.split(",");
  whitelistedArr.forEach((site) => {
    if (url.includes(site)) {
      result = false;
    }
  });
  return result;
};

// src/content.ts
(async () => {
  const isDisallowed = await isUrlDisallowed(window.location.href);
  if (isDisallowed) {
    return;
  }
  console.log("GD Page State Saver Loading...");
  let pageState = null;
  let pageStateLoaded = false;
  const createButtons = () => {
    const container = document.createElement("div");
    container.id = "scroll-state-saver-container";
    const button1 = document.createElement("button");
    button1.id = "scroll-state-saver-btn";
    button1.textContent = "0";
    container.onclick = () => {
      if (pageState && pageState.url) {
        alert(pageState.visibleText);
      }
    };
    container.appendChild(button1);
    const button2 = document.createElement("button");
    button2.id = "scroll-state-saver-scrollpos-btn";
    container.appendChild(button2);
    document.body.appendChild(container);
  };
  document.addEventListener("scrollend", async () => {
    if (pageStateLoaded) {
      await chrome.runtime.sendMessage({
        type: "scrollStopped",
        data: {
          url: window.location.href,
          scrollPos: Math.trunc(window.scrollY),
          visibleText: getVisibleText()
        }
      });
    }
  });
  chrome.runtime.onMessage.addListener((message) => {
    if (message.type === "PAGE_STATE_LOADED") {
      console.log("Page state loaded", message.data);
      pageState = message.data;
      pageStateLoaded = true;
      const btn = document.getElementById("scroll-state-saver-scrollpos-btn");
      if (btn) {
        btn.textContent = pageState?.scrollPos.toString() || "0";
      }
      void sendHealthCheck();
    }
    if (message.type === "HEALTH_CHECK_SUCCESSFUL") {
      console.log("Health check successful");
      const container = document.getElementById("scroll-state-saver-container");
      if (container) {
        container.style.border = "2px solid green";
      }
    }
    if (message.type === "HEALTH_CHECK_UNSUCCESSFUL") {
      console.log("Health check unsuccessful");
      const container = document.getElementById("scroll-state-saver-container");
      if (container) {
        container.style.border = "2px solid red";
      }
    }
  });
  createButtons();
  window.addEventListener("scroll", () => {
    const btn = document.getElementById("scroll-state-saver-btn");
    if (btn) {
      btn.textContent = Math.trunc(window.scrollY).toString();
    }
  });
  let healthIntervalId = null;
  const sendHealthCheck = async () => {
    try {
      await chrome.runtime.sendMessage({ type: "triggerHealthCheck" });
    } catch (_err) {
    }
  };
  const startHealthChecks = () => {
    if (healthIntervalId !== null) return;
    healthIntervalId = window.setInterval(() => {
      if (document.visibilityState === "visible") {
        void sendHealthCheck();
      }
    }, 3e4);
  };
  const stopHealthChecks = () => {
    if (healthIntervalId === null) return;
    clearInterval(healthIntervalId);
    healthIntervalId = null;
  };
  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "visible") {
      startHealthChecks();
      void sendHealthCheck();
    } else {
      stopHealthChecks();
    }
  });
  if (document.visibilityState === "visible") {
    startHealthChecks();
  }
})();
