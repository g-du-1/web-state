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

  const baseUrl = "http://192.168.0.11:8080";
  const saveUrl = `${baseUrl}/pagestate`;
  const getLatestUrl = `${baseUrl}/pagestate/latest`;

  let latestState = null;
  let stateButton = null;

  const createButton = () => {
    const button = document.createElement("button");

    button.style.cssText = `
      position: fixed;
      bottom: 20px;
      right: 20px;
      z-index: 10000;
      background: rgba(0, 0, 0, 0.5);
      color: white;
      border-radius: 4px;
      cursor: pointer;
      font-size: 11px;
      min-width: 100px;
    `;

    button.onclick = () => {
      if (latestState) {
        alert(latestState.visibleText);
      } else {
        alert("No saved state found");
      }
    };

    document.body.appendChild(button);
    stateButton = button;
    updateButtonText();
  };

  const updateButtonText = () => {
    if (!stateButton) return;

    const currentScroll = Math.trunc(window.scrollY);

    const maxScroll = Math.trunc(
      document.documentElement.scrollHeight - window.innerHeight
    );

    if (latestState) {
      stateButton.textContent = `${currentScroll} / ${latestState.scrollPos}`;
    } else {
      stateButton.textContent = `${currentScroll} / ${maxScroll}`;
    }
  };

  const savePageState = () => {
    const payload = {
      url: window.location.href,
      scrollPos: Math.trunc(window.scrollY),
      visibleText: "TODO: Not implemented in js",
    };

    GM_xmlhttpRequest({
      method: "POST",
      url: saveUrl,
      data: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      onload: (response) => {
        console.log("Saved Page State:", JSON.parse(response.responseText));
      },
      onerror: (error) => {
        console.error("Page State Save:", error);
      },
    });
  };

  const getLatestPageState = () => {
    GM_xmlhttpRequest({
      method: "GET",
      url: getLatestUrl + "?url=" + encodeURIComponent(window.location.href),
      onload: (response) => {
        latestState = JSON.parse(response.responseText);
        console.log("Latest Page State:", latestState);
        updateButtonText();
      },
      onerror: (error) => {
        console.error("Latest Page State:", error);
      },
    });
  };

  window.addEventListener("load", () => {
    createButton();
    getLatestPageState();
  });

  window.addEventListener("scroll", updateButtonText);
  window.addEventListener("scrollend", savePageState);
})();
