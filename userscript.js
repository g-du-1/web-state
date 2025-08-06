// ==UserScript==
// @name         Scroll State Saver
// @namespace    http://tampermonkey.net/
// @version      2025-08-06
// @description  try to take over the world!
// @author       You
// @match        https://www.reddit.com/*
// @icon         https://www.google.com/s2/favicons?sz=64&domain=reddit.com
// @grant        GM_xmlhttpRequest
// ==/UserScript==

(function () {
  "use strict";

  const isFullyInViewport = (element) => {
    const rect = element.getBoundingClientRect();

    return (
      rect.top >= 0 &&
      rect.left >= 0 &&
      rect.bottom <= window.innerHeight &&
      rect.right <= window.innerWidth
    );
  };

  const getElementsAtViewportCenter = () => {
    const centerX = window.innerWidth / 2;
    const centerY = window.innerHeight / 2;
    return document.elementsFromPoint(centerX, centerY);
  };

  const getVisibleTextAtCenter = () => {
    const elements = getElementsAtViewportCenter();

    const visibleText = elements
      .filter((el) => isFullyInViewport(el) && el.textContent?.trim())
      .map((el) => el.textContent.trim())
      .join(" ")
      .replace(/\s+/g, " ")
      .trim();

    return visibleText;
  };

  const savePageState = async () => {
    const payload = {
      url: window.location.href,
      scrollPos: Math.trunc(window.scrollY),
    };

    GM_xmlhttpRequest({
      method: "POST",
      url: "http://localhost:8080/pagestate",
      data: JSON.stringify(payload),
      headers: {
        "Content-Type": "application/json",
      },
      onload: function (response) {
        console.log("Success:", response.responseText);
      },
      onerror: function (error) {
        console.error("Error:", error);
      },
    });
  };

  document.addEventListener("scrollend", async () => {
    console.log("Visible: ", getVisibleTextAtCenter());
    savePageState;
  });
})();
