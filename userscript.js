// ==UserScript==
// @name         Scroll State Saver
// @namespace    http://tampermonkey.net/
// @version      2025-08-06
// @description  Saves page state.
// @author       You
// @match        https://www.reddit.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=reddit.com
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
  "use strict";

  const baseUrl = "http://192.168.0.11:8080/api/v1";
  const saveUrl = `${baseUrl}/pagestate/save`;
  const getLatestUrl = `${baseUrl}/pagestate?url=`;

  let latestState = null;

  const createButtons = () => {
    const container = document.createElement("div");
    container.id = "scroll-state-saver-container";
    container.style.cssText = `
      position: fixed;
      bottom: 0;
      right: 0;
      z-index: 10000;
    `;

    const button1 = document.createElement("button");

    button1.id = "scroll-state-saver-btn";

    button1.style.cssText = `
      background: rgba(0, 0, 0, 0.5);
      color: white;
      cursor: pointer;
      border-radius: 0;
      font-size: 11px;
      min-width: 50px;
    `;

    button1.textContent = "0";

    container.onclick = () => {
      if (latestState.url) {
        alert(
          latestState.url +
            "\n\n" +
            latestState.scrollPos +
            "\n\n" +
            latestState.visibleText
        );
      }
    };

    container.appendChild(button1);

    const button2 = document.createElement("button");

    button2.id = "scroll-state-saver-scrollpos-btn";

    button2.style.cssText = `
      background: rgba(0, 0, 0, 0.5);
      color: white;
      border-radius: 0;
      cursor: pointer;
      font-size: 11px;
      min-width: 50px;
    `;

    button2.textContent = latestState?.scrollPos ? latestState.scrollPos : "0";

    container.appendChild(button2);

    window.addEventListener("pagestateloaded", () => {
      button2.textContent = latestState?.scrollPos
        ? latestState.scrollPos
        : "0";
    });

    document.body.appendChild(container);
  };

  const getLatestPageState = () => {
    GM_xmlhttpRequest({
      method: "GET",
      url: getLatestUrl + encodeURIComponent(window.location.href),
      onload: (response) => {
        latestState = JSON.parse(response.responseText);

        window.dispatchEvent(
          new CustomEvent("pagestateloaded", {
            detail: {
              name: "latestState",
              value: latestState,
            },
          })
        );
      },
      onerror: (error) => {},
    });
  };

  window.addEventListener("load", () => {
    createButtons();
    getLatestPageState();
  });

  let lastUrl = location.href;

  setInterval(() => {
    if (location.href !== lastUrl) {
      lastUrl = location.href;
      getLatestPageState();
    }
  }, 500);

  const getVisibleText = () => {
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
      null,
      false
    );

    let visibleText = [];
    let node;

    while ((node = walker.nextNode())) {
      const rect = node.parentElement.getBoundingClientRect();

      if (
        rect.top >= topMargin &&
        rect.bottom <= bottomMargin &&
        rect.left < viewport.width &&
        rect.right > 0
      ) {
        const text = node.textContent.trim();

        if (text && text.length > 25) {
          visibleText.push(text + "\n\n");
        }
      }
    }

    return visibleText.join(" ");
  };

  const updateButtonText = () => {
    const stateButton = document.getElementById("scroll-state-saver-btn");

    if (stateButton) {
      stateButton.textContent = `${Math.trunc(window.scrollY)}`;
    }
  };

  window.addEventListener("scroll", updateButtonText);

  const savePageState = () => {
    const payload = {
      url: window.location.href,
      scrollPos: Math.trunc(window.scrollY),
      visibleText: getVisibleText(),
    };

    GM_xmlhttpRequest({
      method: "POST",
      url: saveUrl,
      data: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      onload: (response) => {},
      onerror: (error) => {},
    });
  };

  let saveTimeout;

  window.addEventListener("scrollend", () => {
    clearTimeout(saveTimeout);
    saveTimeout = setTimeout(savePageState, 2500);
  });
})();
