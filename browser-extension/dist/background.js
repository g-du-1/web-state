// src/background.ts
chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
  if (changeInfo.status === "complete" && tab.url?.includes("reddit.com")) {
    console.log("Tab updated", tab.url);
    const { redditPageStateUrl } = await chrome.storage.local.get(
      "redditPageStateUrl"
    );
    try {
      const getLatestUrl = `${redditPageStateUrl}/pagestate?url=${encodeURIComponent(
        tab.url
      )}`;
      const response = await fetch(getLatestUrl);
      const data = await response.json();
      console.log("Page state fetched", data);
      await chrome.tabs.sendMessage(tabId, {
        type: "PAGE_STATE_LOADED",
        data
      });
    } catch (error) {
      console.log("Failed to fetch page state:", error);
      await chrome.tabs.sendMessage(tabId, {
        type: "PAGE_STATE_LOADED",
        data: null
      });
    }
  }
});
chrome.runtime.onMessage.addListener(async (message, sender) => {
  if (message.type === "scrollStopped") {
    console.log("scrollStopped message received", message.data);
    const { redditPageStateUrl } = await chrome.storage.local.get(
      "redditPageStateUrl"
    );
    try {
      const response = await fetch(`${redditPageStateUrl}/pagestate/save`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(message.data)
      });
      const data = await response.json();
      console.log("Page state saved", data);
    } catch (error) {
      console.log("Failed to save page state:", error);
    }
  }
  if (message.type === "triggerHealthCheck") {
    const targetTabId = sender.tab?.id;
    if (!targetTabId) return;
    const [activeTab] = await chrome.tabs.query({
      active: true,
      lastFocusedWindow: true
    });
    if (!activeTab || activeTab.id !== targetTabId) return;
    const { redditPageStateUrl } = await chrome.storage.local.get(
      "redditPageStateUrl"
    );
    try {
      const response = await fetch(`${redditPageStateUrl}/health`);
      if (response.ok) {
        console.log("Health check successful");
        await chrome.tabs.sendMessage(targetTabId, {
          type: "HEALTH_CHECK_SUCCESSFUL",
          data: null
        });
      } else {
        console.log("Health check failed");
        await chrome.tabs.sendMessage(targetTabId, {
          type: "HEALTH_CHECK_UNSUCCESSFUL",
          data: null
        });
      }
    } catch (err) {
      console.log("Health check error:", err);
      await chrome.tabs.sendMessage(targetTabId, {
        type: "HEALTH_CHECK_UNSUCCESSFUL",
        data: null
      });
    }
  }
});
